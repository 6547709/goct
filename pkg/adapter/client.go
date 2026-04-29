package adapter

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/openlyinc/pointy"
	apiclient "github.com/smartxworks/cloudtower-go-sdk/v2/client"
	userPkg "github.com/smartxworks/cloudtower-go-sdk/v2/client/user"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"

	"github.com/6547709/goct/pkg/debug"
)

// 业务级错误哨兵；可被上层 errors.Is 识别用于映射 exit code。
var (
	ErrAuth       = errors.New("authentication failed")
	ErrNotFound   = errors.New("not found")
	ErrTaskFailed = errors.New("task failed")
	ErrUnsupported = errors.New("operation not supported by SDK")
)

// Client 是 adapter 的聚合接口；各资源 sub-interface 通过内嵌扩展。
type Client interface {
	About(ctx context.Context) (TowerInfo, error)
	VMOps
	SnapshotOps
	HostOps
	ClusterOps
	DatastoreOps
	NetworkOps
	VLANOps
	TaskOps
	AlertOps
	UserOps
	// GetTaskProgress 实现 task.Ops 接口，让 watcher 能轮询任务状态。
	GetTaskProgress(ctx context.Context, id string) (percent int, status string, err error)
}

// Options 是 NewClient 入参。
//
//	URL          形如 https://tower.example.com 或带 /v2/api 前缀
//	Source       登录源（local/ldap/sso/authn），缺省 local
//	Insecure     true 跳过 TLS 校验（自签名内网用）
//	Token        非空时跳过登录直接复用（session 命中场景）
type Options struct {
	URL      string
	Username string
	Password string
	Source   string
	Insecure bool
	Token    string
}

// defaultClient 是 Client 的默认实现。
// 持有 *apiclient.Cloudtower（SDK 主客户端）+ *httptransport.Runtime
// （用于动态注入鉴权头）。
type defaultClient struct {
	api       *apiclient.Cloudtower
	transport *httptransport.Runtime
}

// NewClient 构造 adapter.Client。
//
//   - 提供 Options.Token 时直接注入鉴权头跳过登录
//   - 未提供时调用 SDK Login 并把返回的 token 作为 SessionToken 回吐给上层缓存
//
// 返回的 SessionToken.Value 即可用于后续 Options.Token 复用。
func NewClient(_ context.Context, opts Options) (Client, SessionToken, error) {
	host, basePath, schemes, err := splitURL(opts.URL)
	if err != nil {
		return nil, SessionToken{}, err
	}
	debug.Debugf("adapter: connecting to %s (basePath=%s, insecure=%v)", host, basePath, opts.Insecure)
	tr := newTransport(host, basePath, schemes, opts.Insecure)

	if opts.Token != "" {
		debug.Debug("adapter: reusing cached token (skip login)")
		tr.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", opts.Token)
		api := apiclient.New(tr, strfmt.Default)
		return &defaultClient{api: api, transport: tr}, SessionToken{Value: opts.Token}, nil
	}

	api := apiclient.New(tr, strfmt.Default)
	debug.Debugf("adapter: performing login as %q (source=%s)", opts.Username, opts.Source)
	tok, err := login(api, tr, opts)
	if err != nil {
		return nil, SessionToken{}, err
	}
	return &defaultClient{api: api, transport: tr}, tok, nil
}

// login 调 SDK User.Login，成功后把 token 注入 transport 默认鉴权。
//
// CloudTower 未公开 token TTL；此处假设 8 小时（典型 web 会话），
// session 缓存到期前即使未真到期也会因 401 被自动刷新。
func login(api *apiclient.Cloudtower, tr *httptransport.Runtime, opts Options) (SessionToken, error) {
	src := models.UserSourceLOCAL
	if opts.Source != "" {
		switch strings.ToLower(opts.Source) {
		case "local":
			src = models.UserSourceLOCAL
		case "ldap":
			src = models.UserSource("LDAP")
		case "sso":
			src = models.UserSource("SSO")
		case "authn":
			src = models.UserSource("AUTHN")
		default:
			src = models.UserSource(strings.ToUpper(opts.Source))
		}
	}

	p := userPkg.NewLoginParams()
	p.RequestBody = &models.LoginInput{
		Username: pointy.String(opts.Username),
		Password: pointy.String(opts.Password),
		Source:   models.NewUserSource(src),
	}
	res, err := api.User.Login(p)
	if err != nil {
		return SessionToken{}, fmt.Errorf("%w: %v", ErrAuth, err)
	}
	if res.Payload == nil || res.Payload.Data == nil || res.Payload.Data.Token == nil {
		return SessionToken{}, fmt.Errorf("%w: empty token in response", ErrAuth)
	}
	tok := *res.Payload.Data.Token
	tr.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", tok)
	return SessionToken{Value: tok, ExpireAt: time.Now().Add(8 * time.Hour)}, nil
}

// newTransport 构造 OpenAPI runtime；insecure=true 时跳过 TLS 校验。
func newTransport(host, basePath string, schemes []string, insecure bool) *httptransport.Runtime {
	if insecure {
		hc := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 — 内网自签名场景显式启用
			},
		}
		return httptransport.NewWithClient(host, basePath, schemes, hc)
	}
	return httptransport.New(host, basePath, schemes)
}

// splitURL 解析用户传入的 endpoint，返回 OpenAPI runtime 需要的三元组。
//
// 接受形式：
//
//	https://tower.example.com
//	https://tower.example.com/
//	https://tower.example.com/v2/api（自动归一化）
func splitURL(raw string) (host, basePath string, schemes []string, err error) {
	u, perr := url.Parse(raw)
	if perr != nil || u.Host == "" || u.Scheme == "" {
		return "", "", nil, fmt.Errorf("invalid URL %q", raw)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", nil, fmt.Errorf("unsupported scheme %q (want http|https)", u.Scheme)
	}
	return u.Host, "/v2/api", []string{u.Scheme}, nil
}

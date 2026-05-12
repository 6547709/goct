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
	TemplateOps
	MetricsOps
	LabelOps
	VMFolderOps
	VMPlacementGroupOps
	SnapshotPlanOps
	ElfStoragePolicyOps
	GlobalSettingsOps
	UsbDeviceOps
	ApplicationOps
	DeployOps
	LicenseOps
	ClusterSettingsOps
	NtpOps
	AlertRuleOps
	ContentLibraryImageOps
	CloudTowerApplicationOps
	EventOps
	// GetTaskProgress 实现 task.Ops 接口，让 watcher 能轮询任务状态。
	GetTaskProgress(ctx context.Context, id string) (percent int, status string, err error)
}

// Options 是 NewClient 入参。
//
//	URL          形如 https://tower.example.com 或带 /v2/api 前缀
//	Source       登录源（local/ldap/sso/authn），缺省 local
//	Insecure     true 跳过 TLS 校验（自签名内网用）
//	Token        非空时跳过登录直接复用（session 命中场景）
//	Transport    可选自定义 http.RoundTripper，优先级高于 insecure 设置
type Options struct {
	URL       string
	Username  string
	Password  string
	Source    string
	Insecure  bool
	Token     string
	Transport http.RoundTripper
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
	tr := newTransport(host, basePath, schemes, opts.Insecure, opts.Transport)

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
// CloudTower 未在 login 响应里返回 token TTL；为了让缓存的 token 在真实失效后
// 尽快被 invalidate（不至于让用户感觉每次都"先失败一次"），这里采用
// **更短的 TTL（30 分钟）作为缓存窗口**。即使服务端实际 TTL 更长，缓存提前过期
// 的代价只是多一次登录；反过来，如果服务端实际 TTL 比 8h 短，原实现会让客户端
// 每次都先 401 一次再 fallback，体验更差。
//
// 401 自动 fallback 已在 pkg/client.New 实现，因此这里的 TTL 只是缓存生命期的
// "悲观估计"，不是 token 的真实有效期。
const cachedTokenTTL = 30 * time.Minute

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
	return SessionToken{Value: tok, ExpireAt: time.Now().Add(cachedTokenTTL)}, nil
}

// newTransport 构造 OpenAPI runtime。
// insecure=true 时构造跳过 TLS 校验的 base transport。
// customTransport 不为空时作为 wrapper 套在 base 之上（典型用法：trace round-tripper），
// 它的 Base 字段会被回填为下层 base，从而 trace 与 insecure 可以**叠加**而非互斥。
//
// customTransport 实现的接口约定：
//   - 必须是 *debug.TraceRoundTripper 或类似带 Base 字段的 wrapper；
//   - 或者用户已经显式塞了 base（baseSetter 接口）。
//
// Bug fix（v0.2.1）：原实现 `if customTransport != nil { hc.Transport = customTransport }`
// 会把 insecure 设置整段覆盖，导致 `--trace --insecure` 同开时仍做 TLS 校验。
func newTransport(host, basePath string, schemes []string, insecure bool, customTransport http.RoundTripper) *httptransport.Runtime {
	var base http.RoundTripper
	if insecure {
		base = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 — 内网自签名场景显式启用
		}
	}

	hc := &http.Client{}
	switch {
	case customTransport != nil && base != nil:
		// 把 base 注入到 wrapper 里（让 trace 走 insecure 的 transport）。
		if bs, ok := customTransport.(baseSetter); ok {
			bs.SetBase(base)
		}
		hc.Transport = customTransport
	case customTransport != nil:
		hc.Transport = customTransport
	case base != nil:
		hc.Transport = base
	}
	return httptransport.NewWithClient(host, basePath, schemes, hc)
}

// baseSetter 让 wrapper 类型（如 debug.TraceRoundTripper）声明"我可以接 base transport"。
// 见 pkg/debug/trace.go: TraceRoundTripper.SetBase。
type baseSetter interface {
	SetBase(http.RoundTripper)
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

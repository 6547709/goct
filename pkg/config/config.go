// Package config 加载并合并 goct 运行时配置。
//
// 三级优先级（高 → 低）：
//  1. CLI flag（通过 Override 传入）
//  2. 环境变量 GOCT_*
//  3. 配置文件 ~/.goct.yaml（或显式 ConfigFile）
//
// Resolved 结构是其余包的唯一真相来源；不要直接读 viper。
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Resolved 是合并后最终生效的配置。
type Resolved struct {
	URL      string
	Username string
	Password string
	Cluster  string
	Source   string // local|ldap|sso|authn，缺省 local
	Insecure bool
}

// Override 由调用方（通常是 cmd 层）传入，对应 CLI flag。
// 字段为空字符串视为「未设置」；bool 类型用 InsecureSet 显式标注是否提供。
type Override struct {
	URL        string
	Username   string
	Password   string
	Cluster    string
	Source     string
	ConfigFile string // 可选：覆盖默认 ~/.goct.yaml

	Insecure    bool
	InsecureSet bool // true 表示用户显式传过 --insecure
}

// Resolve 按三级优先级合并配置并返回 Resolved。
//
// 当 ConfigFile 显式给出但读取失败时返回错误；
// 否则即使 ~/.goct.yaml 缺失也按"无 file 输入"继续。
func Resolve(o Override) (Resolved, error) {
	v := viper.New()
	v.SetEnvPrefix("GOCT")
	v.AutomaticEnv()
	for _, k := range []string{"url", "username", "password", "insecure", "cluster", "source"} {
		_ = v.BindEnv(k)
	}
	v.SetDefault("source", "local")

	if o.ConfigFile != "" {
		v.SetConfigFile(o.ConfigFile)
		if err := v.ReadInConfig(); err != nil {
			return Resolved{}, fmt.Errorf("read config %s: %w", o.ConfigFile, err)
		}
	} else {
		v.SetConfigName(".goct")
		v.SetConfigType("yaml")
		v.AddConfigPath("$HOME")
		v.AddConfigPath(".")
		if err := v.ReadInConfig(); err != nil {
			// 文件缺失允许，其它 IO 错误抛出
			var nf viper.ConfigFileNotFoundError
			if !errors.As(err, &nf) {
				return Resolved{}, fmt.Errorf("read default config: %w", err)
			}
		}
	}

	r := Resolved{
		URL:      pick(o.URL, v.GetString("url")),
		Username: pick(o.Username, v.GetString("username")),
		Password: pick(o.Password, v.GetString("password")),
		Cluster:  pick(o.Cluster, v.GetString("cluster")),
		Source:   pick(o.Source, v.GetString("source")),
	}
	if o.InsecureSet {
		r.Insecure = o.Insecure
	} else {
		r.Insecure = v.GetBool("insecure")
	}
	return r, nil
}

// pick 返回第一个非空字符串；用于实现"高优先级覆盖"。
func pick(xs ...string) string {
	for _, x := range xs {
		if x != "" {
			return x
		}
	}
	return ""
}

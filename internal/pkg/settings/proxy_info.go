package settings

import "fmt"

type ProxyInfo struct {
	Enable   bool   `json:"enable"`
	HttpUrl  string `json:"http_url"`
	HttpPort int    `json:"http_port"`
}

func NewProxyInfo() *ProxyInfo {
	return &ProxyInfo{
		Enable:   false,
		HttpUrl:  "",
		HttpPort: 0,
	}
}

func (p *ProxyInfo) GetHttpProxyUrl() string {
	return fmt.Sprintf("http://%s:%d", p.HttpUrl, p.HttpPort)
}

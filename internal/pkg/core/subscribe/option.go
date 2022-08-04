package subscribe

import (
	"time"
)

type ProxyProtocol int

const (
	NONE ProxyProtocol = iota + 1
	SOCKS
	HTTP
)

type UpdateOption struct {
	Key       string
	ProxyMode ProxyProtocol
	Addr      string
	Port      int
	Timeout   time.Duration
}

func NewUpdateOption(proxyMode ProxyProtocol, addr string, port int, timeout time.Duration) *UpdateOption {
	return &UpdateOption{ProxyMode: proxyMode, Addr: addr, Port: port, Timeout: timeout}
}

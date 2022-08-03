package types

type ProxyProtocol int

const (
	NONE ProxyProtocol = iota + 1
	SOCKS
	HTTP
)

package settings

type OneProxySettings struct {
	PID                  int    // Xray 程序进程的PID
	HttpPort             int    // HTTP 代理的端口
	SocksPort            int    // SOCKS 代理的端口
	AllowLanConn         bool   // 允许局域网连接
	Sniffing             bool   // 流量地址监听
	RelayUDP             bool   // 转发 UDP
	DNSPort              int    // DNS 端口
	DNSForeign           string // 国外的 DNS
	DNSDomestic          string // 国内的 DNS
	DNSDomesticBackup    string // 国内的 DNS 备用
	BypassLANAndMainLand bool   // 绕过局域网和大陆
	RoutingStrategy      string // 路由策略
	Mux                  bool   // 多路复用
}

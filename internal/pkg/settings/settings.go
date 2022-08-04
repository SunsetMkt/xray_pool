package settings

type AppSettings struct {
	OneNodeTestTimeOut      int    // 单个节点测试超时时间，单位：秒
	BatchNodeTestMaxTimeOut int    // 批量节点测试的最长超时时间，单位：秒
	TestUrl                 string // 测试代理访问速度的url
	MainProxySettings       ProxySettings
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		OneNodeTestTimeOut:      3,
		BatchNodeTestMaxTimeOut: 100,
		TestUrl:                 "https://google.com",
		MainProxySettings: *NewProxySettings(
			1080,
			0,
			false,
			true,
			true,
			13500,
			"1.1.1.1",
			"119.29.29.29",
			"114.114.114.114",
			true,
			"IPIfNonMatch",
			false,
		),
	}
}

type ProxySettings struct {
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

func NewProxySettings(httpPort int, socksPort int, allowLanConn bool, sniffing bool, relayUDP bool, DNSPort int, DNSForeign string, DNSDomestic string, DNSDomesticBackup string, bypassLANAndMainLand bool, routingStrategy string, mux bool) *ProxySettings {
	return &ProxySettings{HttpPort: httpPort, SocksPort: socksPort, AllowLanConn: allowLanConn, Sniffing: sniffing, RelayUDP: relayUDP, DNSPort: DNSPort, DNSForeign: DNSForeign, DNSDomestic: DNSDomestic, DNSDomesticBackup: DNSDomesticBackup, BypassLANAndMainLand: bypassLANAndMainLand, RoutingStrategy: routingStrategy, Mux: mux}
}

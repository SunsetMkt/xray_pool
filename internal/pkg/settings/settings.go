package settings

type AppSettings struct {
	AppStartPort            int           `json:"app_start_port"`               // 本程序启动的端口，用于 WebUI 登录
	XrayPortRange           string        `json:"xray_port_range"`              // xray 程序的端口范围，从这个范围中找到空闲的端口来使用 36000-  只需要填写起始的端口号，range，会根据 Node 的数量取补全
	XrayInstanceCount       int           `json:"xray_instance_count"`          // Xray 程序的实例数量，简单说就是开启多少个代理
	XrayOpenSocksAndHttp    bool          `json:"xray_open_socks_and_http"`     // 是否开启 socks 和 http 端口，默认只开启 socks 端口
	OneNodeTestTimeOut      int           `json:"one_node_test_time_out"`       // 单个节点测试超时时间，单位：秒
	BatchNodeTestMaxTimeOut int           `json:"batch_node_test_max_time_out"` // 批量节点测试的最长超时时间，单位：秒
	TestUrl                 string        `json:"test_url"`                     // 测试代理访问速度的url
	TestUrlThread           int           `json:"test_url_thread"`              // 测试代理访问速度的url的线程数量
	MainProxySettings       ProxySettings `json:"main_proxy_settings"`
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		AppStartPort:            19035,
		XrayPortRange:           "36000",
		XrayInstanceCount:       3,
		XrayOpenSocksAndHttp:    false,
		OneNodeTestTimeOut:      6,
		BatchNodeTestMaxTimeOut: 100,
		TestUrl:                 "https://google.com",
		TestUrlThread:           10,
		MainProxySettings: *NewProxySettings(
			0,
			1080,
			false,
			true,
			true,
			0,
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

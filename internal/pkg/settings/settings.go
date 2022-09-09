package settings

import "github.com/allanpk716/xray_pool/internal/pkg/common"

type AppSettings struct {
	UserName                string   `json:"user_name"`
	Password                string   `json:"password"`
	AppStartPort            int      `json:"app_start_port"`               // 本程序启动的端口，用于 WebUI 登录
	ManualLbPort            int      `json:"manual_lb_port"`               // 手动指定的负载均衡端口，如果为 0 则自动分配
	XrayPortRange           string   `json:"xray_port_range"`              // xray 程序的端口范围，从这个范围中找到空闲的端口来使用 36000-  只需要填写起始的端口号，range，会根据 Node 的数量取补全
	XrayInstanceCount       int      `json:"xray_instance_count"`          // Xray 程序的实例数量，简单说就是开启多少个代理
	XrayOpenSocksAndHttp    bool     `json:"xray_open_socks_and_http"`     // 是否开启 socks 和 http 端口，默认只开启 socks 端口
	OneNodeTestTimeOut      int      `json:"one_node_test_time_out"`       // 单个节点测试超时时间，单位：秒
	BatchNodeTestMaxTimeOut int      `json:"batch_node_test_max_time_out"` // 批量节点测试的最长超时时间，单位：秒
	HealthCheckUrl          string   `json:"health_check_url"`             // glider 健康检查的 Url，可以与 TestUrl 相同，但是这样可能会浪费连接这个网站的次数。比如，一分钟只允许3次，Health Check 也是会浪费的，根据你的情况来考虑
	TestUrl                 string   `json:"test_url"`                     // 测试代理访问速度的url
	TestUrlThread           int      `json:"test_url_thread"`              // 测试代理访问速度的url的线程数量
	TestUrlHardWay          bool     `json:"test_url_hard_way"`            // 使用 go-rod 启动浏览器来进行测试
	TestUrlFailedWords      []string `json:"test_url_failed_words"`        // 测试这个网站是否有效的关键词，注意是失效的关键词，不是正常的。注意这里需要填入的是小写的，在内部会进行大小写的转换
	TestUrlFailedRegex      string   `json:"test_url_failed_regex"`        // 测试这个网站是否有效的正则表达式，如果匹配到这个正则表达式，则认为这个网站是无效的
	TestUrlStatusCode       int      `json:"test_url_status_code"`         // 期望的网页 StatusCode，一般来说是 200 ，默认是0也就是不进行检查
	/*
		Glider 负载均衡策略：
			rr: round robin
			ha: high availability
			lha: latency based high availability
			dh: destination hashing
	*/
	GliderStrategy    string        `json:"glider_strategy"`
	MainProxySettings ProxySettings `json:"main_proxy_settings"`
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		UserName:                "",
		Password:                "",
		AppStartPort:            common.DefAppStartPort,
		ManualLbPort:            0,
		XrayPortRange:           "36000",
		XrayInstanceCount:       3,
		XrayOpenSocksAndHttp:    false,
		OneNodeTestTimeOut:      6,
		BatchNodeTestMaxTimeOut: 100,
		TestUrl:                 "https://google.com",
		TestUrlThread:           10,
		TestUrlHardWay:          false,
		TestUrlFailedWords:      []string{},
		TestUrlFailedRegex:      "",
		TestUrlStatusCode:       0,
		GliderStrategy:          "rr",
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
	/*
		路由策略 "AsIs" | "IPIfNonMatch" | "IPOnDemand"
			"AsIs"：只使用域名进行路由选择。默认值。
			"IPIfNonMatch"：当域名没有匹配任何规则时，将域名解析成 IP（A 记录或 AAAA 记录）再次进行匹配；
			当一个域名有多个 A 记录时，会尝试匹配所有的 A 记录，直到其中一个与某个规则匹配为止；
			解析后的 IP 仅在路由选择时起作用，转发的数据包中依然使用原始域名；
			"IPOnDemand"：当匹配时碰到任何基于 IP 的规则，将域名立即解析为 IP 进行匹配；
	*/
	RoutingStrategy string
	Mux             bool // 多路复用
}

func NewProxySettings(httpPort int, socksPort int, allowLanConn bool, sniffing bool, relayUDP bool, DNSPort int, DNSForeign string, DNSDomestic string, DNSDomesticBackup string, bypassLANAndMainLand bool, routingStrategy string, mux bool) *ProxySettings {
	return &ProxySettings{HttpPort: httpPort, SocksPort: socksPort, AllowLanConn: allowLanConn, Sniffing: sniffing, RelayUDP: relayUDP, DNSPort: DNSPort, DNSForeign: DNSForeign, DNSDomestic: DNSDomestic, DNSDomesticBackup: DNSDomesticBackup, BypassLANAndMainLand: bypassLANAndMainLand, RoutingStrategy: routingStrategy, Mux: mux}
}

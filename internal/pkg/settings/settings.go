package settings

type AppSettings struct {
	OneNodeTestTimeOut      int    // 单个节点测试超时时间
	BatchNodeTestMaxTimeOut int    // 批量节点测试的最长超时时间
	TestUrl                 string // 测试代理访问速度的url
}

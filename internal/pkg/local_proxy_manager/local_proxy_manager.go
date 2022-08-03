package local_proxy_manager

// LocalProxyManager 本地所有代理实例的管理者
type LocalProxyManager struct {
	localProxyList []*OneProxy // 本地开启多个代理的实例，每个对应着一个 Xray 程序
}

func NewLocalProxyManager() *LocalProxyManager {
	return &LocalProxyManager{localProxyList: make([]*OneProxy, 0)}
}

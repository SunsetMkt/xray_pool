package manager

import (
	coreSettings "github.com/allanpk716/xray_pool/internal/pkg/core/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
)

// Manager 本地所有代理实例的管理者
type Manager struct {
	AppSettings    *settings.AppSettings            // 主程序的配置
	localProxyList []*coreSettings.OneProxySettings // 本地开启多个代理的实例，每个对应着一个 Xray 程序
}

func NewManager() *Manager {
	return &Manager{localProxyList: make([]*coreSettings.OneProxySettings, 0)}
}

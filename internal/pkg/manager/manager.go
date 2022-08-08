package manager

import (
	"encoding/json"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/xray_helper"
	"net/http"
	"os"
	"sync"
)

// Manager 本地所有代理实例的管理者
type Manager struct {
	AppSettings           *settings.AppSettings     `json:"app_settings"` // 主程序的配置
	Subscribes            []*subscribe.Subscribe    `json:"subscribes"`   // 订阅地址
	NodeList              []*node.Node              `json:"nodes"`        // 存放所有的节点
	Filter                []*node.Filter            `json:"filter"`       // 存放所有的过滤器
	xrayHelperList        []*xray_helper.XrayHelper // 本地开启多个代理的实例，每个对应着一个 Xray 程序
	xrayPoolRunning       bool                      // Xray 程序是否正在运行
	xrayPoolRunningLock   sync.Mutex                // Xray 程序是否正在运行的锁
	reverseServer         *http.Server              // 反向代理服务器实例
	reverseServerHttpPort int                       // 反向代理服务器的端口
	reverseServerLocker   sync.Mutex                // 反向代理服务器的锁
	reverseServerRunning  bool                      // 反向代理服务器是否正在运行
	transportManager      *TransportManager         // 代理实例的管理者
	route                 *routing.Routing          // 路由
	wg                    sync.WaitGroup
}

func NewManager() *Manager {

	manager := &Manager{
		AppSettings:    settings.NewAppSettings(),
		Subscribes:     make([]*subscribe.Subscribe, 0),
		NodeList:       make([]*node.Node, 0),
		Filter:         make([]*node.Filter, 0),
		xrayHelperList: make([]*xray_helper.XrayHelper, 0),
		route:          routing.NewRouting(),
	}
	if _, err := os.Stat(core.AppSettings); os.IsNotExist(err) {
		manager.Save()
	} else {
		file, _ := os.Open(core.AppSettings)
		defer func() {
			_ = file.Close()
		}()
		err = json.NewDecoder(file).Decode(manager)
		if err != nil {
			logger.Error(err)
		}
		manager.NodeForEach(func(i int, n *node.Node) {
			n.ParseData()
		})
	}
	return manager
}

// Save 保存数据
func (m *Manager) Save() {
	err := pkg.WriteJSON(m, core.AppSettings)
	if err != nil {
		logger.Error(err)
	}
}

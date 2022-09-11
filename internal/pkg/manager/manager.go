package manager

import (
	"encoding/json"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"github.com/allanpk716/xray_pool/internal/pkg/glider_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/xray_aio"
	"os"
	"sync"
)

// Manager 本地所有代理实例的管理者
type Manager struct {
	AppSettings           *settings.AppSettings       `json:"app_settings"` // 主程序的配置
	Subscribes            []*subscribe.Subscribe      `json:"subscribes"`   // 订阅地址
	NodeList              []*node.Node                `json:"nodes"`        // 存放所有的节点
	Filter                []*node.Filter              `json:"filter"`       // 存放所有的过滤器
	xrayPoolRunning       bool                        // Xray 程序是否正在运行
	xrayPoolRunningLock   sync.Mutex                  // Xray 程序是否正在运行的锁
	xrayAIO               *xray_aio.XrayAIO           // Xray 程序的管理者
	gliderHelper          *glider_helper.GliderHelper // 正向代理的实例
	forwardServerLocker   sync.Mutex                  // 正向代理服务器的锁
	forwardServerHttpPort int                         // 正向代理服务器的端口
	forwardServerRunning  bool                        // 正向代理服务器是否正在运行
	routing               *routing.Routing            // 路由
	wg                    sync.WaitGroup
}

func NewManager() *Manager {

	manager := &Manager{
		AppSettings: settings.NewAppSettings(),
		Subscribes:  make([]*subscribe.Subscribe, 0),
		NodeList:    make([]*node.Node, 0),
		Filter:      make([]*node.Filter, 0),
		routing:     routing.NewRouting(),
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
			logger.Panicf("Decode Config xray_pool_config.json, %v", err)
		}
		manager.NodeForEach(func(i int, n *node.Node) {
			n.ParseData()
		})
	}
	nowXrayHelper := xray_aio.NewXrayOne(0, nil, manager.AppSettings, manager.AppSettings.MainProxySettings, manager.routing, nil)
	defer nowXrayHelper.Stop()
	//if nowXrayHelper.Check() == false {
	//	logger.Panic("NewXrayHelper Check == false")
	//}

	manager.gliderHelper = glider_helper.NewGliderHelper()
	//if manager.gliderHelper.Check() == false {
	//	logger.Panic("glider Check == false")
	//}

	return manager
}

func (m *Manager) CheckGliderStatus() bool {
	return m.gliderHelper.Check()
}

func (m *Manager) CheckXrayStatus() bool {
	return m.xrayAIO.Check()
}

// Save 保存数据
func (m *Manager) Save() {
	err := pkg.WriteJSON(m, core.AppSettings)
	if err != nil {
		logger.Error(err)
	}
}

// Start 启动
func (m *Manager) Start(targetSiteUrl string) bool {

	m.xrayPoolRunningLock.Lock()
	defer m.xrayPoolRunningLock.Unlock()

	if targetSiteUrl != "" {
		m.AppSettings.TestUrl = targetSiteUrl
		m.Save()
	}
	// 检查可用的端口和可用的Node
	bok, aliveNodeIndexList, alivePorts := m.GetsValidNodesAndAlivePorts()
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: GetsValidNodesAndAlivePorts failed")
		return false
	}
	// 开启本地的代理
	bok = m.StartXray(aliveNodeIndexList, alivePorts)
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: StartXray failed")
		return false
	}
	// 开启 glider 前置代理
	bok = m.ForwardProxyStart()
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: ForwardProxyStart failed")
		return false
	}

	m.xrayPoolRunning = true
	return false
}

// Stop 停止
func (m *Manager) Stop() {

	m.xrayPoolRunningLock.Lock()
	defer m.xrayPoolRunningLock.Unlock()

	m.ForwardProxyStop()

	m.StopXray()

	m.xrayPoolRunning = false

	logger.Infof("Stop: xrayPoolRunning = %v", m.xrayPoolRunning)
}

func (m *Manager) XrayPoolRunning() bool {
	m.xrayPoolRunningLock.Lock()
	defer m.xrayPoolRunningLock.Unlock()
	return m.xrayPoolRunning
}

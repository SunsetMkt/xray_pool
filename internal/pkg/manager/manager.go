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
	"github.com/allanpk716/xray_pool/internal/pkg/xray_helper"
	"os"
	"sync"
)

// Manager 本地所有代理实例的管理者
type Manager struct {
	AppSettings         *settings.AppSettings     `json:"app_settings"` // 主程序的配置
	Subscribes          []*subscribe.Subscribe    `json:"subscribes"`   // 订阅地址
	NodeList            []*node.Node              `json:"nodes"`        // 存放所有的节点
	Filter              []*node.Filter            `json:"filter"`       // 存放所有的过滤器
	xrayHelperList      []*xray_helper.XrayHelper // 本地开启多个代理的实例，每个对应着一个 Xray 程序
	xrayPoolRunning     bool                      // Xray 程序是否正在运行
	xrayPoolRunningLock sync.Mutex                // Xray 程序是否正在运行的锁

	gliderHelper          *glider_helper.GliderHelper // 正向代理的实例
	forwardServerLocker   sync.Mutex                  // 正向代理服务器的锁
	forwardServerHttpPort int                         // 正向代理服务器的端口
	forwardServerRunning  bool                        // 正向代理服务器是否正在运行

	route *routing.Routing // 路由
	wg    sync.WaitGroup
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
		manager.save()
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

	manager.gliderHelper = glider_helper.NewGliderHelper()
	if manager.gliderHelper.Check() == false {
		logger.Panic("glider Check == false")
	}

	return manager
}

// save 保存数据
func (m *Manager) save() {
	err := pkg.WriteJSON(m, core.AppSettings)
	if err != nil {
		logger.Error(err)
	}
}

// Start 启动
func (m *Manager) Start() {

	m.xrayPoolRunningLock.Lock()
	defer m.xrayPoolRunningLock.Unlock()

	// 检查可用的端口和可用的Node
	bok, aliveNodeIndexList, alivePorts := m.GetsValidNodesAndAlivePorts()
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: GetsValidNodesAndAlivePorts failed")
		return
	}
	// 开启本地的代理
	bok = m.StartXray(aliveNodeIndexList, alivePorts)
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: StartXray failed")
		return
	}
	// 开启 glider 前置代理
	bok = m.ForwardProxyStart()
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: ForwardProxyStart failed")
		return
	}

	m.xrayPoolRunning = true
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

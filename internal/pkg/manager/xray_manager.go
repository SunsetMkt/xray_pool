package manager

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/xray_aio"
	"github.com/go-rod/rod"
	"github.com/panjf2000/ants/v2"
	"github.com/tklauser/ps"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// GetsValidNodesAndAlivePorts 获取有效的节点和端口信息
func (m *Manager) GetsValidNodesAndAlivePorts(browser *rod.Browser) (bool, []int, []int) {

	defer pkg.TimeCost()("GetsValidNodesAndAlivePorts")

	aliveNodeIndexList := make([]int, 0)
	defer func() {
		logger.Infoln("------------------------------")
		logger.Infof("Alive Node Count: %v", len(aliveNodeIndexList))
		for _, nodeIndex := range aliveNodeIndexList {
			if nodeIndex <= 0 {
				continue
			}
			logger.Infof("alive node: %v -- %v", nodeIndex, m.GetNode(nodeIndex).GetName())
		}
		logger.Infoln("------------------------------")
	}()

	defer func() {
		if browser != nil {
			_ = browser.Close()
		}
	}()

	// 然后需要扫描一个连续的端口段，便于后续的分配
	// 这里需要根据 Node 的数量来推算一个连续的端口段
	needTestPortCount := m.NodeLen()
	if m.AppSettings.XrayOpenSocksAndHttp == true {
		needTestPortCount *= 5
	}
	startRange, err := strconv.Atoi(m.AppSettings.XrayPortRange)
	if err != nil {
		logger.Errorf("xray port range Atoi error: %v", err)
		return false, nil, nil
	}
	portRange := fmt.Sprintf("%d-%d", startRange, startRange+needTestPortCount)
	alivePorts := pkg.ScanAlivePortList(portRange)
	if alivePorts == nil || len(alivePorts) == 0 {
		logger.Errorf("没有找到可用的端口段: %s", portRange)
		return false, nil, nil
	}
	// 默认只需要考虑 socks 的端口，如果需要同时开启 http 端口，则需要2倍
	needMinPortsCount := m.NodeLen()
	if m.AppSettings.XrayOpenSocksAndHttp == true {
		needMinPortsCount = needMinPortsCount * 2
	}
	if len(alivePorts) < needMinPortsCount {
		logger.Errorf("没有找到足够的端口段: %s", portRange)
		return false, nil, nil
	}
	// 是否有足够的空闲、有效的节点，进行了一次粗略的 TCP 排序
	m.NodesTCPing()
	// --------------------------------------------
	// 开启所有的节点，然后再一个个进行检测
	firstTimeNodeIndexList := make([]int, 0)
	m.NodeForEach(func(nIndex int, node *node.Node) {
		firstTimeNodeIndexList = append(firstTimeNodeIndexList, nIndex)
	})
	bok := m.StartXray(firstTimeNodeIndexList, alivePorts)
	defer func() {
		m.StopXray()
	}()
	if bok == false {
		logger.Errorf("StartProxyPoolHandler: StartXray failed")
		return false, nil, nil
	}
	// 临时开启的端口有那些
	firTimeOpenedProxyPorts := m.GetOpenedProxyPorts()
	// --------------------------------------------
	// 接收测试结果
	checkResultChan := make(chan CheckResult, 1)
	defer func() {
		close(checkResultChan)
	}()
	exitRevResultChan := make(chan bool, 1)
	defer close(exitRevResultChan)

	go func() {
		for {
			select {
			case revCheckResult, bok := <-checkResultChan:
				if bok == false {
					continue
				}
				aliveNodeIndexList = append(aliveNodeIndexList, revCheckResult.NodeIndex)
			case <-exitRevResultChan:
				return
			}
		}
	}()
	// --------------------------------------------
	var wg sync.WaitGroup
	// 然后需要并发取完成 Xray 的启动并且通过代理访问目标网站取进行延迟的评价
	p, err := ants.NewPoolWithFunc(m.AppSettings.TestUrlThread, func(inData interface{}) {
		deliveryInfo := inData.(DeliveryInfo)
		defer func() {
			deliveryInfo.Wg.Done()
		}()
		// 测试这节点
		var speedResult int
		if m.AppSettings.TestUrlHardWay == true && deliveryInfo.Browser != nil {
			speedResult, _ = xray_aio.TestNodeByRod(m.AppSettings, deliveryInfo.Browser, deliveryInfo.OpenResult.HttpPort)
		} else {
			speedResult, _ = xray_aio.TestNode(m.AppSettings.TestUrl, deliveryInfo.OpenResult.SocksPort, m.AppSettings.OneNodeTestTimeOut)
		}

		if speedResult > 0 {
			checkResultChan <- CheckResult{
				NodeIndex: deliveryInfo.NodeIndex,
				Delay:     speedResult,
			}
		} else {
			logger.Infof("节点 %d %s 测试失败", deliveryInfo.NodeIndex, deliveryInfo.OpenResult.Name)
		}
	})
	if err != nil {
		logger.Errorf("创建 xray 工作池失败: %v", err)
		return false, nil, nil
	}
	defer p.Release()

	for i, openResult := range firTimeOpenedProxyPorts {

		wg.Add(1)
		err = p.Invoke(DeliveryInfo{
			Browser:    browser,
			NodeIndex:  i + 1,
			OpenResult: openResult,
			Wg:         &wg,
		})
		if err != nil {
			logger.Errorf("xray 工作池提交任务失败: %v", err)
			return false, nil, nil
		}
	}

	wg.Wait()

	return true, aliveNodeIndexList, alivePorts
}

// StartXray 批量启动 Xray 开启代理
func (m *Manager) StartXray(aliveNodeIndexList, alivePorts []int) bool {

	defer pkg.TimeCost()("StartXray")
	// 获取有效的节点列表
	nodes := make([]*node.Node, 0)
	for _, aliveNodeIndex := range aliveNodeIndexList {
		nodes = append(nodes, m.GetNode(aliveNodeIndex))
	}
	// 需要切割出 socks 和 http 的端口
	aliveSocksPorts := make([]int, 0)
	aliveHttpPorts := make([]int, 0)
	// 如果开启了 http 才需要分 http 端口出来
	maxLen := len(alivePorts)
	if len(alivePorts)%2 != 0 {
		maxLen -= 1
	}
	if m.AppSettings.XrayOpenSocksAndHttp == true {
		for i := 0; i < maxLen; i++ {
			if i%2 == 0 {
				aliveSocksPorts = append(aliveSocksPorts, alivePorts[i])
			} else {
				aliveHttpPorts = append(aliveHttpPorts, alivePorts[i])
			}
		}
	} else {
		aliveSocksPorts = alivePorts
	}

	m.xrayAIO = xray_aio.NewXrayAIO(nodes, m.AppSettings, m.routing, aliveSocksPorts, aliveHttpPorts)
	if m.xrayAIO.Check() == false {
		logger.Errorf("xray Check Error")
		return false
	}
	return m.xrayAIO.StartMix()
}

func (m *Manager) StopXray() bool {

	if m.xrayAIO != nil {
		m.xrayAIO.Stop()
	}
	m.KillAllXray()
	return true
}

func (m *Manager) GetOpenedProxyPorts() []xray_aio.OpenResult {

	if m.xrayAIO != nil {
		return m.xrayAIO.GetOpenedProxyPorts()
	} else {
		return []xray_aio.OpenResult{}
	}
}

func (m *Manager) KillAllXray() {

	logger.Debugln("结束所有的 xray，开始...")
	defer logger.Debugln("结束所有的 xray，完成")
	processes, err := ps.Processes()
	if err != nil {
		logger.Errorf("get processes error: %v", err)
		return
	}
	xrayName := pkg.GetXrayExeName()

	for _, p := range processes {
		if strings.ToLower(filepath.Base(p.ExecutablePath())) == xrayName {
			x, err := os.FindProcess(p.PID())
			if err != nil {
				logger.Errorf("find process error: %v", err)
				continue
			}
			err = x.Kill()
			if err != nil {
				logger.Errorf("kill process error: %v", err)
				continue
			}
		}
	}
}

type DeliveryInfo struct {
	Browser    *rod.Browser
	NodeIndex  int
	OpenResult xray_aio.OpenResult
	Wg         *sync.WaitGroup
}

type CheckResult struct {
	NodeIndex int // 当前的 Node Index
	Delay     int // ms
}

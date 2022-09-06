package manager

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/rod_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
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
func (m *Manager) GetsValidNodesAndAlivePorts() (bool, []int, []int) {

	defer pkg.TimeCost()("GetsValidNodesAndAlivePorts")

	aliveNodeIndexList := make([]int, 0)

	defer func() {
		logger.Infoln("------------------------------")
		logger.Infof("Alive Node Count: %v", len(aliveNodeIndexList))
		for _, nodeIndex := range aliveNodeIndexList {
			logger.Infof("alive node: %v -- %v", nodeIndex, m.GetNode(nodeIndex).GetName())
		}
		logger.Infoln("------------------------------")
	}()

	browser := rod_helper.NewBrowser()
	defer func() {
		_ = browser.Close()
	}()

	// 首先需要找到当前系统中残留的 xray 程序，结束它们
	m.KillAllXray()
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

	checkResultChan := make(chan CheckResult, 1)
	defer close(checkResultChan)
	exitRevResultChan := make(chan bool, 1)
	defer close(exitRevResultChan)
	go func() {
		for {
			select {
			case revCheckResult := <-checkResultChan:
				aliveNodeIndexList = append(aliveNodeIndexList, revCheckResult.NodeIndex)
			case <-exitRevResultChan:
				return
			}
		}
	}()

	var wg sync.WaitGroup
	// 然后需要并发取完成 Xray 的启动并且通过代理访问目标网站取进行延迟的评价
	p, err := ants.NewPoolWithFunc(m.AppSettings.TestUrlThread, func(inData interface{}) {
		deliveryInfo := inData.(DeliveryInfo)

		var nowXrayOne *xray_aio.XrayAIO
		defer func() {
			if nowXrayOne != nil {
				nowXrayOne.Stop()
			}
			deliveryInfo.Wg.Done()
		}()

		nowXrayOne = xray_aio.NewXrayOne(deliveryInfo.StartIndex,
			m.GetNode(deliveryInfo.NowNodeIndex),
			deliveryInfo.AppSettings,
			deliveryInfo.NowProxySettings,
			m.routing,
			deliveryInfo.Browser)
		if nowXrayOne.Check() == false {
			logger.Errorf("xray Check Error")
			return
		}

		bok, delay := nowXrayOne.StartOne(
			m.AppSettings.TestUrl,
			m.AppSettings.OneNodeTestTimeOut,
			false,
		)
		if bok == true {
			// 需要记录当前的 Node Index 信息
			checkResultChan <- CheckResult{
				NodeIndex: deliveryInfo.NowNodeIndex,
				Delay:     delay,
			}
		}
	})
	if err != nil {
		logger.Errorf("创建 xray 工作池失败: %v", err)
		return false, nil, nil
	}
	defer p.Release()

	alivePortIndex := 0
	m.NodeForEach(func(nIndex int, node *node.Node) {

		// 设置 socks 端口
		nowProxySettings := m.AppSettings.MainProxySettings
		socksPort := alivePorts[alivePortIndex]
		alivePortIndex++
		nowProxySettings.SocksPort = socksPort
		// 设置 http 端口
		if m.AppSettings.XrayOpenSocksAndHttp == true {
			httpPort := alivePorts[alivePortIndex]
			alivePortIndex++
			nowProxySettings.HttpPort = httpPort
		}

		wg.Add(1)
		err = p.Invoke(DeliveryInfo{
			Browser:          browser,
			StartIndex:       nIndex,
			AppSettings:      m.AppSettings,
			NowProxySettings: nowProxySettings,
			NowNodeIndex:     nIndex,
			Wg:               &wg,
		})
		if err != nil {
			logger.Errorf("xray 工作池提交任务失败: %v", err)
			return
		}
	})

	wg.Wait()
	exitRevResultChan <- true

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
	Browser          *rod.Browser
	StartIndex       int
	AppSettings      *settings.AppSettings
	NowProxySettings settings.ProxySettings
	NowNodeIndex     int
	Wg               *sync.WaitGroup
}

type CheckResult struct {
	NodeIndex int // 当前的 Node Index
	Delay     int // ms
}

package manager

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/xray_helper"
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

	// 首先需要找到当前系统中残留的 xray 程序，结束它们
	m.KillAllXray()
	// 然后需要扫描一个连续的端口段，便于后续的分配
	// 这里需要根据 Node 的数量来推算一个连续的端口段
	needTestPortCount := m.NodeLen()
	if m.AppSettings.XrayOpenSocksAndHttp == true {
		needTestPortCount *= 2
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
	needMinPortsCount := m.AppSettings.XrayInstanceCount
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

		var nowXrayHelper *xray_helper.XrayHelper
		defer func() {
			if nowXrayHelper != nil {
				nowXrayHelper.Stop()
			}
			deliveryInfo.Wg.Done()
		}()

		nowXrayHelper = xray_helper.NewXrayHelper(deliveryInfo.StartIndex, deliveryInfo.NowProxySettings, m.routing)
		if nowXrayHelper.Check() == false {
			logger.Errorf("xray Check Error")
			return
		}

		bok, delay := nowXrayHelper.Start(
			m.GetNode(deliveryInfo.NowNodeIndex),
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
			StartIndex:       nIndex,
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

	m.xrayHelperList = make([]*xray_helper.XrayHelper, 0)
	// 开始启动 xray
	selectNodeIndex := 0
	alivePortIndex := 0
	startXrayCount := 0

	var startWg sync.WaitGroup
	for {
		if startXrayCount >= m.AppSettings.XrayInstanceCount || selectNodeIndex > len(aliveNodeIndexList)-1 {
			break
		}
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

		startOne := func(startXrayCount, selectNodeIndex int, nowProxySettings settings.ProxySettings) {
			defer startWg.Done()
			nowXrayHelper := xray_helper.NewXrayHelper(startXrayCount, nowProxySettings, m.routing)
			if nowXrayHelper.Check() == false {
				logger.Errorf("xray Check Error")
				nowXrayHelper.Stop()
				return
			}
			bok, _ := nowXrayHelper.Start(
				m.GetNode(aliveNodeIndexList[selectNodeIndex]),
				m.AppSettings.TestUrl,
				m.AppSettings.OneNodeTestTimeOut,
				true,
			)
			if bok == true {
				m.xrayHelperList = append(m.xrayHelperList, nowXrayHelper)
			}
		}

		startWg.Add(1)
		go startOne(startXrayCount, selectNodeIndex, nowProxySettings)

		startXrayCount++
		selectNodeIndex++
	}

	startWg.Wait()

	return true
}

func (m *Manager) StopXray() bool {

	for _, xrayHelper := range m.xrayHelperList {
		xrayHelper.Stop()
	}

	m.KillAllXray()

	m.xrayHelperList = make([]*xray_helper.XrayHelper, 0)

	return true
}

// GetOpenedProxyPorts 获取 Xray 开启的 socks 端口和 http 端口，是否有 http 端口需要看 AppSettings.XrayOpenSocksAndHttp 设置
func (m *Manager) GetOpenedProxyPorts() []OpenResult {

	openResultList := make([]OpenResult, 0)
	for _, xrayHelper := range m.xrayHelperList {

		now := OpenResult{}
		now.SocksPort = xrayHelper.ProxySettings.SocksPort
		if m.AppSettings.XrayOpenSocksAndHttp == true {
			now.HttpPort = xrayHelper.ProxySettings.HttpPort
		}
		if xrayHelper.Node != nil {
			now.Name = xrayHelper.Node.GetName()
			now.ProtoModel = xrayHelper.Node.GetProtocolMode().String()
			openResultList = append(openResultList, now)
		}
	}

	return openResultList
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
	StartIndex       int
	NowProxySettings settings.ProxySettings
	NowNodeIndex     int
	Wg               *sync.WaitGroup
}

type CheckResult struct {
	NodeIndex int // 当前的 Node Index
	Delay     int // ms
}

type OpenResult struct {
	Name       string `json:"name"`
	ProtoModel string `json:"proto_model"`
	SocksPort  int    `json:"socks_port"`
	HttpPort   int    `json:"http_port"`
}

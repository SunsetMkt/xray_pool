package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/xray_helper"
	"github.com/panjf2000/ants/v2"
	"github.com/tklauser/ps"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GetsValidNodesAndAlivePorts 获取有效的节点和端口信息
func (m *Manager) GetsValidNodesAndAlivePorts() (bool, []int, []int) {

	defer pkg.TimeCost()("GetsValidNodesAndAlivePorts")
	// 首先需要找到当前系统中残留的 xray 程序，结束它们
	m.KillAllXray()
	// 然后需要扫描一个连续的端口段，便于后续的分配
	alivePorts := pkg.ScanAlivePortList(m.AppSettings.XrayPortRange)
	if alivePorts == nil || len(alivePorts) == 0 {
		logger.Errorf("没有找到可用的端口段: %s", m.AppSettings.XrayPortRange)
		return false, nil, nil
	}
	// 默认只需要考虑 socks 的端口，如果需要同时开启 http 端口，则需要2倍
	needMinPortsCount := m.AppSettings.XrayInstanceCount
	if m.AppSettings.XrayOpenSocksAndHttp == true {
		needMinPortsCount = needMinPortsCount * 2
	}
	if len(alivePorts) < needMinPortsCount {
		logger.Errorf("没有找到足够的端口段: %s", m.AppSettings.XrayPortRange)
		return false, nil, nil
	}
	// 是否有足够的空闲、有效的节点，进行了一次粗略的 TCP 排序
	m.NodesTCPing()

	aliveNodeIndexList := make([]int, 0)
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

		nowXrayHelper = xray_helper.NewXrayHelper(deliveryInfo.StartIndex, deliveryInfo.NowProxySettings, m.route)
		if nowXrayHelper.Check() == false {
			logger.Errorf("xray Check Error")
			return
		}

		bok, delay := nowXrayHelper.Start(m.GetNode(deliveryInfo.NowNodeIndex), m.AppSettings.TestUrl, m.AppSettings.OneNodeTestTimeOut)
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

func (m *Manager) StartXray(aliveNodeIndexList, alivePorts []int) bool {

	defer pkg.TimeCost()("StartXray")

	m.xrayHelperList = make([]*xray_helper.XrayHelper, 0)
	// 开始启动 xray
	selectNodeIndex := 0
	alivePortIndex := 0
	startXrayCount := 0
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

		nowXrayHelper := xray_helper.NewXrayHelper(startXrayCount, nowProxySettings, m.route)
		if nowXrayHelper.Check() == false {
			logger.Errorf("xray Check Error")
			return false
		}

		bok, _ := nowXrayHelper.Start(m.GetNode(aliveNodeIndexList[selectNodeIndex]), m.AppSettings.TestUrl, m.AppSettings.OneNodeTestTimeOut)
		if bok == true {
			m.xrayHelperList = append(m.xrayHelperList, nowXrayHelper)
			startXrayCount++
		} else {
			// 如果失败了，那么端口的 Index 需要回退
			alivePortIndex--
			if m.AppSettings.XrayOpenSocksAndHttp == true {
				alivePortIndex--
			}
		}
		selectNodeIndex++
	}

	return true
}

func (m *Manager) StopXray() bool {

	for _, xrayHelper := range m.xrayHelperList {
		xrayHelper.Stop()
	}

	m.KillAllXray()

	err := os.RemoveAll(pkg.GetTmpFolderFPath())
	if err != nil {
		logger.Errorf("remove tmp folder error: %v", err)
		return false
	}

	return true
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

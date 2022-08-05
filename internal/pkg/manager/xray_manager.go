package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/xray_helper"
	"github.com/tklauser/ps"
	"os"

	"path/filepath"
	"strings"
)

func (m *Manager) StartXray() bool {

	m.xrayHelperList = make([]*xray_helper.XrayHelper, 0)
	// 首先需要找到当前系统中残留的 xray 程序，结束它们
	m.KillAllXray()
	// 然后需要扫描一个连续的端口段，便于后续的分配
	alivePorts := pkg.ScanAlivePortList(m.AppSettings.XrayPortRange)
	if alivePorts == nil || len(alivePorts) == 0 {
		logger.Errorf("没有找到可用的端口段: %s", m.AppSettings.XrayPortRange)
		return false
	}
	// 默认只需要考虑 socks 的端口，如果需要同时开启 http 端口，则需要2倍
	needMinPortsCount := m.AppSettings.XrayInstanceCount
	if m.AppSettings.XrayOpenSocksAndHttp == true {
		needMinPortsCount = needMinPortsCount * 2
	}
	if len(alivePorts) < needMinPortsCount {
		logger.Errorf("没有找到足够的端口段: %s", m.AppSettings.XrayPortRange)
		return false
	}
	// 是否有足够的空闲、有效的节点
	m.NodesTCPing()

	// 开始启动 xray
	selectNodeIndex := 1
	alivePortIndex := 0
	startXrayCount := 0
	for {
		if startXrayCount >= m.AppSettings.XrayInstanceCount || selectNodeIndex > m.NodeLen() {
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

		if nowXrayHelper.Start(m.GetNode(selectNodeIndex), m.AppSettings.TestUrl, m.AppSettings.OneNodeTestTimeOut) == true {
			m.xrayHelperList = append(m.xrayHelperList, nowXrayHelper)
			startXrayCount++
		}
		selectNodeIndex++
	}

	logger.Info("批量启动 xray 完成")
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

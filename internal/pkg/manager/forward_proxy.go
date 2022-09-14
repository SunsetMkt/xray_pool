package manager

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
)

func (m *Manager) ForwardProxyStart() bool {

	m.forwardServerLocker.Lock()
	defer func() {
		m.forwardServerLocker.Unlock()
	}()

	if m.forwardServerRunning == true {
		logger.Debugln("Reverse Http Server is already running")
		return true
	}

	openResultList := m.GetOpenedProxyPorts()
	if len(openResultList) == 0 {
		logger.Errorf("ForwardProxyStart: no open ports to proxy")
		return false
	}

	if m.AppSettings.ManualLbPort == 0 {
		// 自动
		// 如果不满足，那么就再次扫描一个端口段，找到一个可用的端口给反向代理服务器
		alivePorts := pkg.ScanAlivePortList("63200-63400")
		if len(alivePorts) == 0 {
			logger.Errorf("自动选择负载均衡端口失败，搜索范围：63200-63400")
			return false
		} else {
			m.forwardServerHttpPort = alivePorts[0]
		}
	} else {
		alivePorts := pkg.ScanAlivePortList(fmt.Sprintf("%s", m.AppSettings.ManualLbPort))
		if len(alivePorts) == 0 {
			logger.Errorf("手动指定负载均衡端口:%d ,已经被占用！", m.AppSettings.ManualLbPort)
			return false
		} else {
			m.forwardServerHttpPort = m.AppSettings.ManualLbPort
		}
	}

	socksPorts := make([]int, 0)
	for _, result := range openResultList {
		socksPorts = append(socksPorts, result.SocksPort)
	}
	if m.gliderHelper.Check() == false {
		logger.Panic("ForwardProxyStart: gliderHelper Check Error")
		return false
	}
	err := m.gliderHelper.Start(
		m.AppSettings.HealthCheckUrl,
		m.AppSettings.HealthCheckInterval,
		m.forwardServerHttpPort,
		socksPorts, m.AppSettings.GliderStrategy)
	if err != nil {
		logger.Errorf("ForwardProxyStart: %s", err)
		return false
	}

	m.forwardServerRunning = true
	logger.Infof("ForwardProxyStart: http port %d", m.forwardServerHttpPort)

	logger.Infoln("代理池开启成功")

	return true
}

func (m *Manager) ForwardProxyStop() {

	m.forwardServerLocker.Lock()
	defer func() {
		m.forwardServerLocker.Unlock()
	}()
	if m.forwardServerRunning == false {
		logger.Debugln("Reverse Http Server is not running")
		return
	}

	if m.gliderHelper == nil {
		logger.Errorf("ForwardProxyStop: gliderHelper is nil")
		return
	}

	err := m.gliderHelper.Stop()
	if err != nil {
		logger.Errorf("ForwardProxyStop: %s", err)
		return
	}

	m.forwardServerRunning = false
}

func (m *Manager) ForwardProxyPort() int {
	return m.forwardServerHttpPort
}

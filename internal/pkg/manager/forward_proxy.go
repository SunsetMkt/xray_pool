package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
)

func (m *Manager) ForwardProxyStart() {

	m.forwardServerLocker.Lock()
	defer func() {
		m.forwardServerLocker.Unlock()
	}()

	if m.forwardServerRunning == true {
		logger.Debugln("Reverse Http Server is already running")
		return
	}

	socksPorts, httpPorts := m.GetOpenedProxyPorts()
	if len(socksPorts) == 0 && len(httpPorts) == 0 {
		logger.Panic("ForwardProxyStart: no open ports to proxy")
	}
	// 如果不满足，那么就再次扫描一个端口段，找到一个可用的端口给反向代理服务器
	alivePorts := pkg.ScanAlivePortList("63200-63400")
	if len(alivePorts) == 0 {
		logger.Panic("ForwardProxyStart: no open ports to proxy")
	} else {
		m.forwardServerHttpPort = alivePorts[0]
	}

	err := m.gliderHelper.Start(m.forwardServerHttpPort, socksPorts)
	if err != nil {
		logger.Panicf("ForwardProxyStart: %s", err)
		return
	}

	m.forwardServerRunning = true

	logger.Infof("ForwardProxyStart: http port %d", m.forwardServerHttpPort)
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
		logger.Panic("ForwardProxyStop: gliderHelper is nil")
		return
	}

	err := m.gliderHelper.Stop()
	if err != nil {
		logger.Panicf("ForwardProxyStop: %s", err)
		return
	}

	m.forwardServerRunning = false
}

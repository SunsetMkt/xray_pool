package xray_aio

import (
	"bufio"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/protocols"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/go-rod/rod"
	"os/exec"
	"time"
)

func NewXrayOne(index int, inNode *node.Node, appSettings *settings.AppSettings, ProxySettings settings.ProxySettings, route *routing.Routing, browser *rod.Browser) *XrayAIO {
	x := XrayAIO{
		index:            index,
		AppSettings:      appSettings,
		OneProxySettings: ProxySettings,
		route:            route,
		browser:          browser,
		startOneOrAll:    true,
		nodes:            make([]*node.Node, 0),
	}

	x.nodes = append(x.nodes, inNode)

	return &x
}

// StartOne 启动一个独立的 xray
func (x *XrayAIO) StartOne(testUrl string, testTimeOut int, skipSpeedTest bool) (bool, int) {

	nowNode := x.nodes[0]
	x.targetUrl = testUrl
	if x.runOne() == true {
		if x.OneProxySettings.HttpPort == 0 {
			logger.Infof("Xray -- %2d 启动成功, 监听 socks 端口: %d, 所选节点: %s",
				x.index,
				x.OneProxySettings.SocksPort, nowNode.GetName())
		} else {
			logger.Infof("Xray -- %2d 启动成功, 监听 socks/http 端口: %d/%d, 所选节点: %s",
				x.index,
				x.OneProxySettings.SocksPort, x.OneProxySettings.HttpPort, nowNode.GetName())
		}

		if skipSpeedTest == true {
			return true, 0
		}

		result := 0
		status := ""
		if x.AppSettings.TestUrlHardWay == false {
			result, status = x.TestNode(testUrl, x.OneProxySettings.SocksPort, testTimeOut)
		} else {
			result, status = x.TestNodeByRod(x.AppSettings, x.browser, testUrl, testTimeOut)
		}
		logger.Infof("Xray -- %2d %6s [ %s ] 延迟: %dms", x.index, status, testUrl, result)
		if result < 0 {
			x.Stop()
			logger.Infof("Xray -- %2d 当前节点: %v 访问 %v 失败, 将不再使用该节点", x.index, nowNode.GetName(), testUrl)
			return false, result
		}

		return true, result
	} else {
		return false, -1
	}
}

func (x *XrayAIO) runOne() bool {

	nowNodeProtocol := x.nodes[0].Protocol
	switch nowNodeProtocol.GetProtocolMode() {
	case protocols.ModeShadowSocks, protocols.ModeTrojan, protocols.ModeVMess, protocols.ModeSocks, protocols.ModeVLESS, protocols.ModeVMessAEAD:
		file := x.GenConfigOne()
		x.xrayCmd = exec.Command(x.xrayPath, "-c", file)
	default:
		logger.Errorf("Xray -- %2d 暂不支持%v协议", x.index, nowNodeProtocol.GetProtocolMode())
		return false
	}
	stdout, err := x.xrayCmd.StdoutPipe()
	if err != nil {
		logger.Errorf("Xray -- %2d 获取 xray 程序的 stdout 管道失败: %s", x.index, err.Error())
		return false
	}
	err = x.xrayCmd.Start()
	if err != nil {
		logger.Errorf("Xray -- %2d 启动 xray 程序失败: %s", x.index, err.Error())
		return false
	}
	r := bufio.NewReader(stdout)
	lines := new([]string)
	go readInfo(r, lines)
	status := make(chan struct{})
	go checkProc(x.xrayCmd, status)
	stopper := time.NewTimer(time.Millisecond * 300)
	select {
	case <-stopper.C:
		x.OneProxySettings.PID = x.xrayCmd.Process.Pid
		return true
	case <-status:
		logger.Error("Xray -- %2d 开启xray服务失败, 查看下面报错信息来检查出错问题", x.index)
		for _, x := range *lines {
			logger.Error(x)
		}
		return false
	}
}

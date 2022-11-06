package xray_aio

import (
	"bufio"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/protocols"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"os/exec"
	"time"
)

func NewXrayAIO(inNodes []*node.Node, appSettings *settings.AppSettings, route *routing.Routing, socksPorts, httpPorts []int) *XrayAIO {

	x := XrayAIO{
		AppSettings:   appSettings,
		route:         route,
		startOneOrAll: false,
		nodes:         make([]*node.Node, 0),
		socksPorts:    socksPorts,
		httpPorts:     httpPorts,
	}

	x.nodes = append(x.nodes, inNodes...)

	return &x
}

func (x *XrayAIO) StartMix() bool {

	x.runningLock.Lock()
	defer x.runningLock.Unlock()

	// 判断是否所有的 node 的协议都支持
	for _, nowNode := range x.nodes {
		nowNodeProtocol := nowNode
		switch nowNodeProtocol.GetProtocolMode() {
		case protocols.ModeShadowSocks, protocols.ModeShadowSocksR, protocols.ModeTrojan, protocols.ModeVMess, protocols.ModeSocks, protocols.ModeVLESS, protocols.ModeVMessAEAD:
		default:
			logger.Errorf("Xray -- %2d 暂不支持%v协议", x.index, nowNodeProtocol.GetProtocolMode())
			return false
		}
	}
	file := x.GenConfigMix()
	x.xrayCmd = exec.Command(x.xrayPath, "-c", file)

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

// GetOpenedProxyPorts 获取 Xray 开启的 socks 端口和 http 端口，是否有 http 端口需要看 AppSettings.XrayOpenSocksAndHttp 设置
func (x *XrayAIO) GetOpenedProxyPorts() []OpenResult {

	x.runningLock.Lock()
	defer x.runningLock.Unlock()

	openResultList := make([]OpenResult, 0)

	for i, nowNode := range x.nodes {
		now := OpenResult{}
		now.SocksPort = x.socksPorts[i]
		if x.AppSettings.XrayOpenSocksAndHttp == true {
			now.HttpPort = x.httpPorts[i]
		}

		now.Name = nowNode.GetName()
		now.ProtoModel = nowNode.GetProtocolMode().String()
		openResultList = append(openResultList, now)
	}

	return openResultList
}

type OpenResult struct {
	Name       string `json:"name"`
	ProtoModel string `json:"proto_model"`
	SocksPort  int    `json:"socks_port"`
	HttpPort   int    `json:"http_port"`
}

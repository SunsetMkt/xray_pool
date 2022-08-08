package v1

import (
	"github.com/WQGroup/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StartProxyPoolHandler 开启代理池
func (cb ControllerBase) StartProxyPoolHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "StartProxyPoolHandler", err)
	}()

	if cb.proxyPoolLocker.Lock() == false || cb.manager.XrayPoolRunning() == true {
		// 已经在执行，跳过
		c.JSON(http.StatusOK, ReplyProxyPool{Status: cb.proxyPoolRunningStatus})
		return
	}

	go func() {

		defer func() {
			cb.proxyPoolLocker.Unlock()
		}()
		cb.proxyPoolRunningStatus = "starting"
		// 检查可用的端口和可用的Node
		bok, aliveNodeIndexList, alivePorts := cb.manager.GetsValidNodesAndAlivePorts()
		if bok == false {
			cb.proxyPoolRunningStatus = "stopped"
			logger.Errorf("StartProxyPoolHandler: GetsValidNodesAndAlivePorts failed")
			return
		}
		// 开启本地的代理
		bok = cb.manager.StartXray(aliveNodeIndexList, alivePorts)
		if bok == false {
			cb.proxyPoolRunningStatus = "stopped"
			logger.Errorf("StartProxyPoolHandler: StartXray failed")
			return
		}

		cb.manager.ForwardProxyStart()
		// 开启反向代理
		cb.proxyPoolRunningStatus = "running"
	}()

	c.JSON(http.StatusOK, ReplyProxyPool{Status: "starting"})
}

// GetProxyListHandler 获取本地开启的代理列表
func (cb ControllerBase) GetProxyListHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "GetProxyListHandler", err)
	}()

	reply := ReplyProxyList{
		LBPort:    cb.manager.ForwardProxyPort(),
		SocksPots: make([]int, 0),
		HttpPots:  make([]int, 0),
	}
	SocksPots, HttpPots := cb.manager.GetOpenedProxyPorts()
	reply.SocksPots = append(reply.SocksPots, SocksPots...)
	reply.HttpPots = append(reply.HttpPots, HttpPots...)

	c.JSON(http.StatusOK, reply)
}

type ReplyProxyPool struct {
	Status string `json:"status"`
}

type ReplyProxyList struct {
	LBPort    int   `json:"lb_port"`
	SocksPots []int `json:"socks_pots"`
	HttpPots  []int `json:"http_pots"`
}

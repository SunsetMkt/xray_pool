package v1

import (
	"github.com/allanpk716/xray_pool/internal/pkg/xray_aio"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StartProxyPoolHandler 开启代理池
func (cb *ControllerBase) StartProxyPoolHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "StartProxyPoolHandler", err)
	}()

	startProxyPool := RequestStartProxyPool{}
	err = c.ShouldBindJSON(&startProxyPool)
	if err != nil {
		return
	}

	if cb.proxyPoolLocker.Lock() == false {
		// 已经在执行，跳过
		c.JSON(http.StatusOK, ReplyProxyPool{Status: cb.proxyPoolRunningStatus})
		return
	}

	if cb.manager.XrayPoolRunning() == true {
		// 已经在执行，跳过
		c.JSON(http.StatusOK, ReplyProxyPool{Status: cb.proxyPoolRunningStatus})
		return
	}

	go func() {

		defer func() {
			cb.proxyPoolLocker.Unlock()
		}()
		cb.proxyPoolRunningStatus = "starting"
		brunok := cb.manager.Start(startProxyPool.TargetSiteUrl)
		if brunok == true {
			// 开启反向代理
			cb.proxyPoolRunningStatus = "running"
		} else {
			cb.manager.Stop()
			cb.proxyPoolRunningStatus = "stopped"
		}
	}()

	c.JSON(http.StatusOK, ReplyProxyPool{Status: "starting"})
}

// StopProxyPoolHandler 关闭代理池
func (cb *ControllerBase) StopProxyPoolHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "StopProxyPoolHandler", err)
	}()

	if cb.manager.XrayPoolRunning() == false {
		// 已经在执行，跳过
		c.JSON(http.StatusOK, ReplyProxyPool{Status: cb.proxyPoolRunningStatus})
		return
	}

	cb.manager.Stop()

	cb.proxyPoolRunningStatus = "stopped"

	c.JSON(http.StatusOK, ReplyProxyPool{Status: cb.proxyPoolRunningStatus})
}

// GetProxyListHandler 获取本地开启的代理列表
func (cb *ControllerBase) GetProxyListHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "GetProxyListHandler", err)
	}()

	reply := ReplyProxyList{
		Status:         cb.proxyPoolRunningStatus,
		LBPort:         cb.manager.ForwardProxyPort(),
		OpenResultList: make([]xray_aio.OpenResult, 0),
	}
	openResultList := cb.manager.GetOpenedProxyPorts()
	reply.OpenResultList = append(reply.OpenResultList, openResultList...)

	c.JSON(http.StatusOK, reply)
}

type RequestStartProxyPool struct {
	TargetSiteUrl string `json:"target_site_url"`
}

type ReplyProxyPool struct {
	Status string `json:"status"`
}

type ReplyProxyList struct {
	Status         string                `json:"status"`
	LBPort         int                   `json:"lb_port"`
	OpenResultList []xray_aio.OpenResult `json:"open_result_list"`
}

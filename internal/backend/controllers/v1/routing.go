package v1

import (
	"fmt"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RoutingAddHandler 添加路由规则
func (cb *ControllerBase) RoutingAddHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "RoutingAddHandler", err)
	}()

	addRouting := RequestAddRouting{}
	err = c.ShouldBindJSON(&addRouting)
	if err != nil {
		return
	}

	cb.manager.AddRule(addRouting.RoutingType, addRouting.Rules...)

	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

// RoutingListHandler 列举路由规则
func (cb *ControllerBase) RoutingListHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "RoutingListHandler", err)
	}()

	routingList := ReplyRoutingList{
		BlockList:  RequestAddRouting{RoutingType: routing.TypeBlock, Rules: make([]string, 0)},
		DirectList: RequestAddRouting{RoutingType: routing.TypeDirect, Rules: make([]string, 0)},
		ProxyList:  RequestAddRouting{RoutingType: routing.TypeProxy, Rules: make([]string, 0)},
	}

	blockRuleList := cb.manager.GetRule(routing.TypeBlock, "all")
	for _, oneRule := range blockRuleList {
		if len(oneRule) == 3 {
			routingList.BlockList.Rules = append(routingList.BlockList.Rules, oneRule[2])
		}
	}

	directRuleList := cb.manager.GetRule(routing.TypeDirect, "all")
	for _, oneRule := range directRuleList {
		if len(oneRule) == 3 {
			routingList.DirectList.Rules = append(routingList.DirectList.Rules, oneRule[2])
		}
	}

	proxyRuleList := cb.manager.GetRule(routing.TypeProxy, "all")
	for _, oneRule := range proxyRuleList {
		if len(oneRule) == 3 {
			routingList.ProxyList.Rules = append(routingList.ProxyList.Rules, oneRule[2])
		}
	}

	c.JSON(http.StatusOK, routingList)
}

// RoutingDeleteHandler 删除路由规则
func (cb *ControllerBase) RoutingDeleteHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "RoutingDeleteHandler", err)
	}()

	delRouting := RequestDeleteRouting{}
	err = c.ShouldBindJSON(&delRouting)
	if err != nil {
		return
	}

	deleteList := ""
	for i, index := range delRouting.IndexList {
		deleteList += fmt.Sprintf("%d", index)
		if i < len(delRouting.IndexList)-1 {
			deleteList += ","
		}
	}
	cb.manager.DelRule(delRouting.RoutingType, deleteList)
}

type RequestAddRouting struct {
	RoutingType routing.Type `json:"routing_type"`
	Rules       []string     `json:"rules"`
}

type ReplyRoutingList struct {
	BlockList  RequestAddRouting `json:"block_list"`  // 屏蔽列表
	DirectList RequestAddRouting `json:"direct_list"` // 直接跳转列表
	ProxyList  RequestAddRouting `json:"proxy_list"`  // 代理列表
}

type RequestDeleteRouting struct {
	RoutingType routing.Type `json:"routing_type"`
	IndexList   []int        `json:"index_list"` // Index 从 1 开始，而不是0
}

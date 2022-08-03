package v1

import (
	"github.com/gin-gonic/gin"
)

// SubscribeListHandler 列举有那些订阅源
func (cb ControllerBase) SubscribeListHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeListHandler", err)
	}()

}

// SubscribeAddHandler 添加一个订阅源
func (cb ControllerBase) SubscribeAddHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeAddHandler", err)
	}()

}

// SubscribeUpdateNodesHandler 订阅源获取节点的逻辑
func (cb ControllerBase) SubscribeUpdateNodesHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeUpdateHandler", err)
	}()

}

// SubscribeUpdateHandler 订阅源更新的逻辑，比如修改订阅源的备注名称，是否启用等
func (cb ControllerBase) SubscribeUpdateHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeUpdateHandler", err)
	}()

}

// SubscribeDelHandler 订阅源删除的逻辑
func (cb ControllerBase) SubscribeDelHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeUpdateHandler", err)
	}()

}

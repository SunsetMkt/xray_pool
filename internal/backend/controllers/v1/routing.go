package v1

import (
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RoutingBlockListHandler 列举有那些禁止的路由规则
func (cb *ControllerBase) RoutingBlockListHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "RoutingBlockListHandler", err)
	}()

	subscribeList := ReplySubscribeList{}
	subscribeList.SubscribeList = make([]subscribe.Subscribe, 0)
	cb.manager.SubscribeForEach(func(index int, subscribe *subscribe.Subscribe) {
		subscribeList.SubscribeList = append(subscribeList.SubscribeList, *subscribe)
	})

	c.JSON(http.StatusOK, subscribeList)
}

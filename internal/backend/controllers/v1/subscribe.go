package v1

import (
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// SubscribeListHandler 列举有那些订阅源
func (cb ControllerBase) SubscribeListHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeListHandler", err)
	}()

	subscribeList := ReplySubscribeList{}
	subscribeList.SubscribeList = make([]subscribe.Subscribe, 0)
	cb.manager.SubscribeForEach(func(index int, subscribe *subscribe.Subscribe) {
		subscribeList.SubscribeList = append(subscribeList.SubscribeList, *subscribe)
	})

	c.JSON(http.StatusOK, subscribeList)
}

// SubscribeAddHandler 添加一个订阅源
func (cb ControllerBase) SubscribeAddHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeAddHandler", err)
	}()

	subscribeAdd := RequestSubscribeAdd{}
	err = c.ShouldBindJSON(&subscribeAdd)
	if err != nil {
		return
	}

	sub := subscribe.NewSubscribe(subscribeAdd.Url, subscribeAdd.Name)
	cb.manager.AddSubscribe(sub)

	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

// SubscribeUpdateNodesHandler 订阅源获取节点的逻辑
func (cb ControllerBase) SubscribeUpdateNodesHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeUpdateHandler", err)
	}()

	opt := subscribe.NewUpdateOption(subscribe.NONE, "", 0, 5*time.Second)
	cb.manager.UpdateNode(opt)

	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

// SubscribeUpdateHandler 订阅源更新的逻辑，比如修改订阅源的备注名称，是否启用等
func (cb ControllerBase) SubscribeUpdateHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeUpdateHandler", err)
	}()

	subscribeUpdate := RequestSubscribeUpdate{}
	err = c.ShouldBindJSON(&subscribeUpdate)
	if err != nil {
		return
	}
	cb.manager.SetSubscribe(
		subscribeUpdate.Index,
		strconv.FormatBool(subscribeUpdate.Using),
		subscribeUpdate.Url,
		subscribeUpdate.Name,
	)

	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

// SubscribeDelHandler 订阅源删除的逻辑
func (cb ControllerBase) SubscribeDelHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SubscribeUpdateHandler", err)
	}()

	subscribeDelete := RequestSubscribeDelete{}
	err = c.ShouldBindJSON(&subscribeDelete)
	if err != nil {
		return
	}
	cb.manager.DelSubscribe(subscribeDelete.Index)

	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

type ReplySubscribeList struct {
	SubscribeList []subscribe.Subscribe `json:"subscribe_list"`
}

type RequestSubscribeAdd struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type RequestSubscribeUpdate struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Using bool   `json:"using"`
	Index string `json:"index"` // 索引从 1 开始，而不是 0，需要从 List 列表中获取的时候 +1
}

type RequestSubscribeDelete struct {
	Index string `json:"index"` // 索引从 1 开始，而不是 0，需要从 List 列表中获取的时候 +1
}

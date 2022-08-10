package v1

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg/lock"
	"github.com/allanpk716/xray_pool/internal/pkg/manager"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerBase struct {
	manager                *manager.Manager
	proxyPoolLocker        lock.Lock
	proxyPoolRunningStatus string
	restartSignal          chan interface{}
	exitSignal             chan interface{}
}

func NewControllerBase(restartSignal, exitSignal chan interface{}) *ControllerBase {
	cb := &ControllerBase{
		restartSignal:          restartSignal,
		exitSignal:             exitSignal,
		manager:                manager.NewManager(),
		proxyPoolRunningStatus: "stopped",
		proxyPoolLocker:        lock.NewLock(),
	}

	return cb
}

func (cb *ControllerBase) GetVersion() string {
	return "v1"
}

func (cb *ControllerBase) GetAppStartPort() int {
	return cb.manager.AppSettings.AppStartPort
}

func (cb *ControllerBase) ErrorProcess(c *gin.Context, funcName string, err error) {
	if err != nil {
		logger.Errorln(funcName, err.Error())
		c.JSON(http.StatusInternalServerError, backend.ReplyCommon{Message: err.Error()})
	}
}

// Close 关闭 HTTP 服务器
func (cb *ControllerBase) Close() {
	if cb != nil {
		if cb.manager != nil {
			cb.manager.Stop()
		}
		cb.proxyPoolLocker.Close()
	}
}

// ExitHandler 退出 APP 的逻辑
func (cb *ControllerBase) ExitHandler(c *gin.Context) {
	cb.exitSignal <- true
}

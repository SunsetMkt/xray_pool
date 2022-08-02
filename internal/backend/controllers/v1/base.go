package v1

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerBase struct {
	restartSignal chan interface{}
}

func NewControllerBase(restartSignal chan interface{}) *ControllerBase {
	cb := &ControllerBase{
		restartSignal: restartSignal,
	}

	return cb
}

func (cb ControllerBase) GetVersion() string {
	return "v1"
}

func (cb *ControllerBase) ErrorProcess(c *gin.Context, funcName string, err error) {
	if err != nil {
		logger.Errorln(funcName, err.Error())
		c.JSON(http.StatusInternalServerError, backend.ReplyCommon{Message: err.Error()})
	}
}
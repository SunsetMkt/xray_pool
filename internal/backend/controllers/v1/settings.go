package v1

import (
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SettingsHandler 设置参数
func (cb *ControllerBase) SettingsHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SettingsHandler", err)
	}()

	switch c.Request.Method {
	case "GET":
		{
			// 回复没有密码的 settings
			c.JSON(http.StatusOK, cb.manager.AppSettings)
		}
	case "PUT":
		{
			appSettings := settings.AppSettings{}
			err = c.ShouldBindJSON(&appSettings)
			if err != nil {
				return
			}

			cb.manager.AppSettings = &appSettings
			cb.manager.Save()

			c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
		}
	}
}

func (cb *ControllerBase) DefSettingsHandler(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "DefSettingsHandler", err)
	}()

	c.JSON(http.StatusOK, settings.NewAppSettings())
}

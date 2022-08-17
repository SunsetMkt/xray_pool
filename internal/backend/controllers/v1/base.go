package v1

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/common"
	"github.com/allanpk716/xray_pool/internal/pkg/lock"
	"github.com/allanpk716/xray_pool/internal/pkg/manager"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"net/http"
	"os"
	"time"

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

	cb.manager.Save()

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

func (cb *ControllerBase) SystemStatus(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SystemStatus", err)
	}()

	replySystemStatus := ReplySystemStatus{}
	if cb.manager.AppSettings.UserName != "" && cb.manager.AppSettings.Password != "" {
		replySystemStatus.IsSetup = true
	} else {
		replySystemStatus.IsSetup = false
	}

	c.JSON(http.StatusOK, replySystemStatus)

}

func (cb *ControllerBase) SetUp(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "SetUp", err)
	}()

	setup := RequestSetup{}
	err = c.ShouldBindJSON(&setup)
	if err != nil {
		return
	}

	if cb.manager.AppSettings.UserName == "" && cb.manager.AppSettings.Password == "" {
		// 可以执行 Setup 流程
		cb.manager.AppSettings.UserName = setup.UserName
		cb.manager.AppSettings.Password = setup.Password
		cb.manager.Save()
		c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
		return
	} else {
		c.JSON(http.StatusOK, backend.ReplyCommon{Message: "already set up"})
	}
}

func (cb *ControllerBase) Login(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "Login", err)
	}()

	login := RequestLogin{}
	err = c.ShouldBindJSON(&login)
	if err != nil {
		return
	}

	if cb.manager.AppSettings.UserName == "" ||
		cb.manager.AppSettings.Password == "" {
		c.JSON(http.StatusOK, backend.ReplyCommon{Message: "username or password error"})
		return
	}

	if cb.manager.AppSettings.UserName == login.UserName &&
		cb.manager.AppSettings.Password == login.Password {
		// 登录成功
		nowToken := pkg.RandStringBytesMaskImprSrcSB(32)
		common.SetAccessToken(nowToken)
		c.JSON(http.StatusOK, ReplyLogin{Message: "ok", Token: nowToken})
		return
	} else {
		c.JSON(http.StatusOK, backend.ReplyCommon{Message: "username or password error"})
		return
	}
}

func (cb *ControllerBase) Logout(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "Logout", err)
	}()

	common.SetAccessToken("")
	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

func (cb *ControllerBase) ChangePWD(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "ChangePWD", err)
	}()

	changePWD := RequestChangePWD{}
	err = c.ShouldBindJSON(&changePWD)
	if err != nil {
		return
	}

	if changePWD.OldPassword == "" || changePWD.NewPassword == "" {
		c.JSON(http.StatusOK, backend.ReplyCommon{Message: "old password or new password is empty"})
		return
	}

	if changePWD.OldPassword != cb.manager.AppSettings.Password {
		c.JSON(http.StatusOK, backend.ReplyCommon{Message: "old password error"})
		return
	}

	cb.manager.AppSettings.Password = changePWD.NewPassword
	cb.manager.Save()

	// 需要重新登录
	common.SetAccessToken("")
	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
}

// ExitHandler 退出 APP 的逻辑
func (cb *ControllerBase) ExitHandler(c *gin.Context) {

	c.JSON(http.StatusOK, backend.ReplyCommon{Message: "ok"})
	time.Sleep(1 * time.Second)
	cb.exitSignal <- true
}

func (cb *ControllerBase) ClearTmpFolder(c *gin.Context) {

	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "ClearTmpFolder", err)
	}()

	if cb.manager.XrayPoolRunning() == true {
		// 已经在执行，跳过
		c.JSON(http.StatusOK, ReplyClearTmpFolder{
			Status:        cb.proxyPoolRunningStatus,
			TmpFolderPath: pkg.GetTmpFolderFPath(),
		})
		return
	}

	err = os.RemoveAll(pkg.GetTmpFolderFPath())
	if err != nil {
		err = fmt.Errorf("remove tmp folder error: %v", err)
		logger.Error(err)
		return
	}

	c.JSON(http.StatusOK, ReplyClearTmpFolder{
		Status:        "ok",
		TmpFolderPath: pkg.GetTmpFolderFPath(),
	})
}

type ReplySystemStatus struct {
	IsSetup bool `json:"is_setup"`
}

type RequestSetup struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type RequestLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type ReplyLogin struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RequestChangePWD struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ReplyClearTmpFolder struct {
	Status        string `json:"status"`
	TmpFolderPath string `json:"tmp_folder_path"`
}

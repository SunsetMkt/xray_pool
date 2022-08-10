package v1

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/lock"
	"github.com/allanpk716/xray_pool/internal/pkg/manager"
	"github.com/allanpk716/xray_pool/internal/pkg/types/backend"
	"net/http"
	"os"

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

type ReplyClearTmpFolder struct {
	Status        string `json:"status"`
	TmpFolderPath string `json:"tmp_folder_path"`
}

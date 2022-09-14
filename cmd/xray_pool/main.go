package main

import (
	_ "embed"
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/backend"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/common"
	"github.com/allanpk716/xray_pool/internal/pkg/logger_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/manager"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"os"
	"runtime"
)

var exitSignal = make(chan interface{}, 1)

func init() {
	common.SetAppVersion(AppVersion)
}

func main() {

	logger.Infoln("Start XrayPool...")

	go logger_helper.Listen()

	restartSignal := make(chan interface{}, 1)
	defer close(restartSignal)
	defer close(exitSignal)
	bend := backend.NewBackEnd(restartSignal, exitSignal)
	go bend.Restart()
	restartSignal <- 1

	systray.Run(onReady, onExit)
}

func onReady() {

	AppStartPort := common.DefAppStartPort
	{
		m := manager.NewManager()
		AppStartPort = m.AppSettings.AppStartPort
	}
	systray.SetIcon(mainICON)

	if runtime.GOOS != "darwin" {
		// macos 的时候，就不设置 title 了，不然太占位置了
		systray.SetTitle("XrayPool")
	}
	systray.SetTooltip("XrayPool - 代理池")
	mMainWindow := systray.AddMenuItem("主程序", "打开主程序窗体")
	mQuit := systray.AddMenuItem("退出", "退出程序，清理缓存")
	go func() {
		<-mQuit.ClickedCh
		exitSignal <- true
		systray.Quit()
	}()

	go func() {
		<-mMainWindow.ClickedCh
		err := open.Run(fmt.Sprintf("http://127.0.0.1:%d", AppStartPort))
		if err != nil {
			logger.Errorln("open.Run", err.Error())
		}
	}()

	// Sets the icon of a menu item. Only available on Mac and Windows.
	//mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
	_ = os.RemoveAll(pkg.GetTmpFolderFPath())
}

//go:embed icon/swimming_pool.ico
var mainICON []byte

/*
	使用 git tag 来做版本描述，然后在编译的时候传入版本号信息到这个变量上
*/
var AppVersion = "unknow"

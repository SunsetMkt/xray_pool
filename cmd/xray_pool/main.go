package main

import (
	_ "embed"
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
	"github.com/allanpk716/xray_pool/internal/backend"
	v1 "github.com/allanpk716/xray_pool/internal/backend/controllers/v1"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/common"
	"github.com/allanpk716/xray_pool/internal/pkg/logger_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/manager"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
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

	m := manager.NewManager()
	httpProxyUrl := ""
	if m.AppSettings.ProxyInfoSettings.Enable == true {
		httpProxyUrl = m.AppSettings.ProxyInfoSettings.GetHttpProxyUrl()
	}
	rod_helper.InitFakeUA(m.AppSettings.CachePath, httpProxyUrl)

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

	if pkg.IsFile(AutoStartPool) == true {
		// 需要自动启动代理池，模拟提交 http 启动请求
		go func() {
			ops := rod_helper.NewHttpClientOptions(15)
			httpClient, err := rod_helper.NewHttpClient(ops)
			if err != nil {
				logger.Panicln("NewHttpClient", err.Error())
			}
			post, err := httpClient.R().
				SetBody(map[string]interface{}{"target_site_url": settings.NewAppSettings().TestUrl}).
				SetResult(&v1.ReplyProxyPool{}).
				Post(fmt.Sprintf("http://127.0.0.1:%d/v1/start_proxy_pool", AppStartPort))
			if err != nil {
				logger.Errorln("Post", err.Error())
				return
			}
			logger.Infoln("AutoStartPool", post.String())
		}()
	}

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

const (
	AutoStartPool = "AutoStartPool"
)

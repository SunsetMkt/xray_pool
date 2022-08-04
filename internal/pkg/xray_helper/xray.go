package xray_helper

import (
	"bufio"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	coreSettings "github.com/allanpk716/xray_pool/internal/pkg/core/settings"
	"github.com/allanpk716/xray_pool/internal/pkg/protocols"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type XrayHelper struct {
	index         int                            // 第几个 xray 实例
	xrayCmd       *exec.Cmd                      // xray 程序的进程
	AppSettings   *settings.AppSettings          // 主程序的配置
	xrayPath      string                         // xray 程序的路径
	proxySettings *coreSettings.OneProxySettings // 代理的配置
	route         *routing.Routing               // 路由
}

func NewXrayHelper(index int, appSettings *settings.AppSettings, xrayPath string, proxySettings *coreSettings.OneProxySettings, route *routing.Routing) *XrayHelper {
	return &XrayHelper{index: index, AppSettings: appSettings, xrayPath: xrayPath, proxySettings: proxySettings, route: route}
}

// Check 检查 Xray 程序和需求的资源是否已经存在，不存在则需要提示用户去下载
func (x *XrayHelper) Check() bool {

	// 在这个目录下进行搜索是否存在 Xray 程序
	nowRootPath := pkg.GetBaseXrayFolderFPath()
	xrayExeName := XrayName
	sysType := runtime.GOOS
	if sysType == "windows" {
		xrayExeName += ".exe"
	}
	xrayExeFullPath := filepath.Join(nowRootPath, xrayExeName)
	if pkg.IsFile(xrayExeFullPath) == false {
		return false
	}
	// 检查 geoip.dat geosite.dat 是否存在
	geoIPResource := filepath.Join(nowRootPath, GEOIP_RESOURCE_NAME)
	geoSiteResource := filepath.Join(nowRootPath, GEOSite_RESOURCE_NAME)
	if pkg.IsFile(geoIPResource) == false || pkg.IsFile(geoSiteResource) == false {
		return false
	}

	x.xrayPath = xrayExeFullPath

	return true
}

func (x *XrayHelper) Start(key string) {
	testUrl := x.AppSettings.TestUrl
	testTimeout := x.AppSettings.OneNodeTestTimeOut
	manager := manage.Manager
	indexList := core.IndexList(key, manager.NodeLen())
	if len(indexList) == 0 {
		log.Warn("没有选取到节点")
	} else if len(indexList) == 1 {
		index := indexList[0]
		node := manager.GetNode(index)
		manager.SetSelectedIndex(index)
		manager.Save()
		exe := run(node.Protocol)
		if exe {
			if setting.Http() == 0 {
				log.Infof("启动成功, 监听socks端口: %d, 所选节点: %d", setting.Socks(), manager.SelectedIndex())
			} else {
				log.Infof("启动成功, 监听socks/http端口: %d/%d, 所选节点: %d", setting.Socks(), setting.Http(), manager.SelectedIndex())
			}
			result, status := TestNode(testUrl, setting.Socks(), testTimeout)
			log.Infof("%6s [ %s ] 延迟: %dms", status, testUrl, result)
		}
	}
}

func (x *XrayHelper) run(node protocols.Protocol) bool {
	x.Stop()
	switch node.GetProtocolMode() {
	case protocols.ModeShadowSocks, protocols.ModeTrojan, protocols.ModeVMess, protocols.ModeSocks, protocols.ModeVLESS, protocols.ModeVMessAEAD:
		file := x.GenConfig(node, x.proxySettings, x.route)
		x.xrayCmd = exec.Command(x.xrayPath, "-c", file)
	default:
		logger.Errorf("暂不支持%v协议", node.GetProtocolMode())
		return false
	}
	stdout, _ := x.xrayCmd.StdoutPipe()
	_ = x.xrayCmd.Start()
	r := bufio.NewReader(stdout)
	lines := new([]string)
	go readInfo(r, lines)
	status := make(chan struct{})
	go checkProc(x.xrayCmd, status)
	stopper := time.NewTimer(time.Millisecond * 300)
	select {
	case <-stopper.C:
		x.proxySettings.PID = x.xrayCmd.Process.Pid
		return true
	case <-status:
		logger.Error("开启xray服务失败, 查看下面报错信息来检查出错问题")
		for _, x := range *lines {
			logger.Error(x)
		}
		return false
	}
}

// Stop 停止服务
func (x *XrayHelper) Stop() {

	if x.xrayCmd != nil {
		err := x.xrayCmd.Process.Kill()
		if err != nil {
			logger.Errorf("停止xray服务失败: %v", err)
		}
		x.xrayCmd = nil
	}
	if x.proxySettings.PID != 0 {
		process, err := os.FindProcess(x.proxySettings.PID)
		if err == nil {
			err = process.Kill()
			if err != nil {
				logger.Errorf("停止xray服务失败: %v", err)
			}
		}
		x.proxySettings.PID = 0
	}
	// 日志文件过大清除
	file, _ := os.Stat(core.LogFile)
	if file != nil {
		fileSize := float64(file.Size()) / (1 << 20)
		if fileSize > 5 {
			err := os.Remove(core.LogFile)
			if err != nil {
				logger.Errorf("清除日志文件失败: %v", err)
			}
		}
	}
}

func readInfo(r *bufio.Reader, lines *[]string) {
	for i := 0; i < 20; i++ {
		line, _, _ := r.ReadLine()
		if len(string(line[:])) != 0 {
			*lines = append(*lines, string(line[:]))
		}
	}
}

// 检查进程状态
func checkProc(c *exec.Cmd, status chan struct{}) {
	_ = c.Wait()
	status <- struct{}{}
}

const (
	GEOIP_RESOURCE_NAME   = "geoip.dat"
	GEOSite_RESOURCE_NAME = "geosite.dat"
	XrayName              = "xray"
)

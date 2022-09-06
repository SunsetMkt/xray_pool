package xray_aio

import (
	"bufio"
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/go-rod/rod"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type XrayAIO struct {
	xrayCmd          *exec.Cmd              // xray 程序的进程
	xrayPath         string                 // xray 程序的路径
	AppSettings      *settings.AppSettings  // 主程序的设置信息
	OneProxySettings settings.ProxySettings // 代理的配置
	route            *routing.Routing       // 路由
	targetUrl        string                 // 目标 url
	browser          *rod.Browser           // 浏览器实例
	index            int                    // 如果是启动单个实例单个端口的情况下，需要填写这个 Index 字段
	startOneOrAll    bool                   // true 表示启动单个，false 表示启动全部
	nodes            []*node.Node           // 节点列表
	socksPorts       []int                  // socks 端口列表
	httpPorts        []int                  // http 端口列表
	runningLock      sync.Mutex
}

// Check 检查 Xray 程序和需求的资源是否已经存在，不存在则需要提示用户去下载
func (x *XrayAIO) Check() bool {

	// 在这个目录下进行搜索是否存在 Xray 程序
	nowRootPath := pkg.GetBaseThingsFolderFPath()
	xrayExeName := pkg.GetXrayExeName()
	xrayExeFullPath := filepath.Join(nowRootPath, xrayExeName)
	if pkg.IsFile(xrayExeFullPath) == false {
		logger.Panic(XrayDownloadInfo)
		return false
	}
	// 检查 geoip.dat geosite.dat 是否存在
	geoIPResource := filepath.Join(nowRootPath, GEOIP_RESOURCE_NAME)
	geoSiteResource := filepath.Join(nowRootPath, GEOSite_RESOURCE_NAME)
	if pkg.IsFile(geoIPResource) == false || pkg.IsFile(geoSiteResource) == false {
		logger.Panic(XrayDownloadInfo)
		return false
	}

	x.xrayPath = xrayExeFullPath

	return true
}

// Stop 停止服务
func (x *XrayAIO) Stop() {

	x.runningLock.Lock()
	defer x.runningLock.Unlock()

	x.targetUrl = ""
	if x.xrayCmd != nil {
		err := x.xrayCmd.Process.Kill()
		if err != nil {
			logger.Errorf("Xray -- %2d 停止xray服务失败: %v", x.index, err)
		}
		x.xrayCmd = nil
		x.OneProxySettings.PID = 0
	}
	// 日志文件过大清除
	if pkg.IsFile(x.GetLogFPath()) == true {
		file, err := os.Stat(x.GetLogFPath())
		if err != nil {
			logger.Errorf("Xray -- %2d os.Stat日志文件大小: %v", x.index, err.Error())
			return
		}
		if file != nil {
			fileSize := float64(file.Size()) / (1 << 20)
			if fileSize > 5 {
				err := os.Remove(x.GetLogFPath())
				if err != nil {
					logger.Errorf("Xray -- %2d 清除日志文件失败: %v", x.index, err)
				}
			}
		}
	}

	// 清理传递进来的节点信息
	x.nodes = make([]*node.Node, 0)
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
)

var (
	XrayDownloadInfo = errors.New(fmt.Sprintf("缺少 Xray 可执行程序或者资源，请去 https://github.com/XTLS/Xray-core/releases 下载对应平台的程序，解压放入 %v 文件夹中", pkg.GetBaseThingsFolderAbsFPath()))
)

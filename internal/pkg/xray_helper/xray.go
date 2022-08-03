package xray_helper

import (
	"github.com/allanpk716/xray_pool/internal/pkg"
	"path/filepath"
	"runtime"
)

type XrayHelper struct {
}

// Check 检查 Xray 程序和需求的资源是否已经存在，不存在则需要提示用户去下载
func (x XrayHelper) Check() bool {

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

	return true
}

const (
	GEOIP_RESOURCE_NAME   = "geoip.dat"
	GEOSite_RESOURCE_NAME = "geosite.dat"
	XrayName              = "xray"
)

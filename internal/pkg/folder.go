package pkg

import (
	"os"
	"runtime"
)

// GetConfigRootDirFPath 获取 Config 的根目录，不同系统不一样
func GetConfigRootDirFPath() string {

	nowConfigFPath := ""
	sysType := runtime.GOOS
	if sysType == "linux" {
		nowConfigFPath = configDirRootFPathLinux
	} else if sysType == "windows" {
		nowConfigFPath = configDirRootFPathWindows
	} else if sysType == "darwin" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("GetConfigRootDirFPath darwin get UserHomeDir, Error:" + err.Error())
		}
		nowConfigFPath = home + configDirRootFPathDarwin
	} else {
		panic("GetConfigRootDirFPath can't matched OSType: " + sysType + " ,You Should Implement It Yourself")
	}

	// 如果文件夹不存在则创建
	if _, err := os.Stat(nowConfigFPath); os.IsNotExist(err) {
		err = os.MkdirAll(nowConfigFPath, os.ModePerm)
		if err != nil {
			panic("GetConfigRootDirFPath mkdir, Error:" + err.Error())
		}
	}

	return nowConfigFPath
}

// GetBaseXrayFolderFPath 获取基础的 Xray 程序存放目录
func GetBaseXrayFolderFPath() string {

	nowPath := GetConfigRootDirFPath() + baseXrayFolderName
	if _, err := os.Stat(nowPath); os.IsNotExist(err) {
		err = os.MkdirAll(nowPath, os.ModePerm)
		if err != nil {
			panic("GetBaseXrayFolderFPath mkdir, Error:" + err.Error())
		}
	}
	return nowPath
}

// GetTmpFolderFPath 获取临时目录文件夹
func GetTmpFolderFPath() string {

	nowPath := GetConfigRootDirFPath() + tmpFolderName
	if _, err := os.Stat(nowPath); os.IsNotExist(err) {
		err = os.MkdirAll(nowPath, os.ModePerm)
		if err != nil {
			panic("GetTmpFolderFPath mkdir, Error:" + err.Error())
		}
	}
	return nowPath
}

// 配置文件的位置信息，这个会根据系统版本做区分
const (
	configDirRootFPathWindows = "."                  // Windows 就是在当前的程序目录
	configDirRootFPathLinux   = "/config"            // Linux 是在 /config 下
	configDirRootFPathDarwin  = "/.config/xray_pool" // Darwin 是在 os.UserHomeDir()/.config/xray_pool/ 下
)

const (
	baseXrayFolderName = "xray"
	tmpFolderName      = "tmp"
)

package pkg

import (
	"fmt"
	"os"
	"path/filepath"
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

// GetBaseThingsFolderFPath 获取基础的支持程序存放目录
func GetBaseThingsFolderFPath() string {

	nowPath := filepath.Join(GetConfigRootDirFPath(), baseFolderName)
	if _, err := os.Stat(nowPath); os.IsNotExist(err) {
		err = os.MkdirAll(nowPath, os.ModePerm)
		if err != nil {
			panic("GetBaseThingsFolderFPath mkdir, Error:" + err.Error())
		}
	}
	return nowPath
}

func GetBaseThingsFolderAbsFPath() string {
	absPath, _ := filepath.Abs(GetBaseThingsFolderFPath())
	return absPath
}

// GetIndexXrayFolderFPath 根据基础的 Xray 程序生成新的 Index 序列号的 Xray 程序存放目录
func GetIndexXrayFolderFPath(index int) string {

	nowPath := filepath.Join(GetTmpFolderFPath(), fmt.Sprintf("%s_%d", xrayFolderName, index))
	if _, err := os.Stat(nowPath); os.IsNotExist(err) {
		err = os.MkdirAll(nowPath, os.ModePerm)
		if err != nil {
			panic("GetIndexXrayFolderFPath mkdir, Error:" + err.Error())
		}
	}
	return nowPath
}

// GetTmpFolderFPath 获取临时目录文件夹
func GetTmpFolderFPath() string {

	nowPath := filepath.Join(GetConfigRootDirFPath(), tmpFolderName)
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
	baseFolderName = "base_things"
	xrayFolderName = "xray"
	tmpFolderName  = "tmp"
)

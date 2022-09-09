package common

import (
	"sync"
)

// SetAccessToken 设置 Web UI 访问的 Token
func SetAccessToken(newToken string) {

	defer mutexAccessToken.Unlock()
	mutexAccessToken.Lock()
	accessToken = newToken
}

// GetAccessToken 获取 Web UI 访问的 Token
func GetAccessToken() string {

	defer mutexAccessToken.Unlock()
	mutexAccessToken.Lock()
	return accessToken
}

func SetAppVersion(newVersion string) {
	appVersion = newVersion
}

func GetAppVersion() string {
	return appVersion
}

var (
	accessToken      = ""
	mutexAccessToken sync.Mutex
	appVersion       = "unknown"
)

const DefAppStartPort = 19038

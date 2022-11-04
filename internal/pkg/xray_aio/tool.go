package xray_aio

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/go-rod/rod"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// TestNode 获取节点代理访问外网的延迟
func TestNode(testUrl string, socks5Port int, timeout int) (int, string) {
	start := time.Now()
	res, e := GetBySocks5Proxy(testUrl, "127.0.0.1", socks5Port, time.Duration(timeout)*time.Second)
	elapsed := time.Since(start)
	if e != nil {
		logger.Warn(e)
		return -1, "Error"
	}
	result, status := int(float32(elapsed.Nanoseconds())/1e6), res.Status
	defer res.Body.Close()
	return result, status
}

// GetBySocks5Proxy 通过Socks5代理访问网站
func GetBySocks5Proxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("socks5://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeOut,
	}
	return client.Get(objUrl)
}

func TestNodeByRod(appSettings *settings.AppSettings,
	browser *rod.Browser,
	localProxyHttpPort int) (int, string) {

	start := time.Now()
	page, e, err := rod_helper.NewPageNavigateWithProxy(browser,
		fmt.Sprintf("http://127.0.0.1:%d", localProxyHttpPort),
		appSettings.TestUrl, time.Duration(appSettings.OneNodeTestTimeOut)*time.Second)
	defer func() {
		if page != nil {
			_ = page.Close()
		}
	}()
	if err != nil {
		return -1, "Error"
	}

	elapsed := time.Since(start)
	speedResult := int(float32(elapsed.Nanoseconds()) / 1e6)

	pageHtmlString, err := page.HTML()
	if err != nil {
		return -1, "Error"
	}

	if appSettings.TestUrlStatusCode != 0 {
		// 需要判断
		if e == nil || e.Response == nil || e.Response.Status != appSettings.TestUrlStatusCode {
			return -1, "Error statusCode"
		}
	}

	for _, word := range appSettings.TestUrlFailedWords {

		if strings.Contains(strings.ToLower(pageHtmlString), word) == true {
			return -1, "Error FailedWord " + word
		}
	}

	if appSettings.TestUrlFailedRegex != "" {
		FailedRegex := regexp.MustCompile(appSettings.TestUrlFailedRegex)
		matches := FailedRegex.FindAllString(pageHtmlString, -1)
		if matches == nil || len(matches) == 0 {
			// 没有找到匹配的内容，那么认为是成功的
		} else {
			// 匹配到了失败的内容，那么认为是失败的
			return -1, "Error FailedRegex"
		}
	}

	return speedResult, "OK"
}

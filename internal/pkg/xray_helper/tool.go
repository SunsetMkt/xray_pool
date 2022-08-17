package xray_helper

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg/rod_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/go-rod/rod"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// TestNode 获取节点代理访问外网的延迟
func (x *XrayHelper) TestNode(url string, port int, timeout int) (int, string) {
	start := time.Now()
	res, e := x.GetBySocks5Proxy(url, "127.0.0.1", port, time.Duration(timeout)*time.Second)
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
func (x *XrayHelper) GetBySocks5Proxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {
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

func (x *XrayHelper) TestNodeByRod(appSettings *settings.AppSettings,
	browser *rod.Browser,
	targetUrl string,
	timeout int) (int, string) {

	start := time.Now()
	page, statusCode, _, err := rod_helper.NewPageNavigate(browser,
		fmt.Sprintf("http://127.0.0.1:%d", x.ProxySettings.HttpPort),
		targetUrl, time.Duration(timeout)*time.Second)
	if err != nil {
		return -1, "Error"
	}
	defer func() {
		if page != nil {
			_ = page.Close()
		}
	}()

	elapsed := time.Since(start)
	speedResult := int(float32(elapsed.Nanoseconds()) / 1e6)

	pageHtmlString, err := page.HTML()
	if err != nil {
		return -1, "Error"
	}

	if appSettings.TestUrlStatusCode != 0 {
		// 需要判断
		if statusCode != appSettings.TestUrlStatusCode {
			return -1, "Error"
		}
	}

	for _, word := range appSettings.TestUrlFailedWords {

		if strings.Contains(strings.ToLower(pageHtmlString), word) == true {
			return -1, "Error"
		}
	}

	FailedRegex := regexp.MustCompile(appSettings.TestUrlFailedRegex)
	matches := FailedRegex.FindAllString(pageHtmlString, -1)
	if matches == nil || len(matches) == 0 {
		// 没有找到匹配的内容，那么认为是成功的
	} else {
		// 匹配到了失败的内容，那么认为是失败的
		return -1, "Error"
	}

	return speedResult, "OK"
}

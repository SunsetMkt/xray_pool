package xray_aio

import (
	"context"
	"fmt"
	"github.com/allanpk716/rod_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/go-rod/rod"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// TestNode 获取节点代理访问外网的延迟
func TestNode(appSettings *settings.AppSettings,
	testUrl string,
	succeedWords []string,
	socks5Port int,
	timeout int) (int, string) {
	start := time.Now()
	res, err := GetBySocks5Proxy(testUrl, "127.0.0.1", socks5Port, time.Duration(timeout)*time.Second)
	elapsed := time.Since(start)
	if err != nil {
		if len(appSettings.TestUrlSucceedWords) > 0 {
			// 无需判断超时，因为后续需要判断页面成功关键词
			if errors.Is(err, context.DeadlineExceeded) == false {
				// 不是超时错误，那么就返回错误，跳过
				return -1, "Error"
			}
		} else {
			// 因为没有设置成功关键词，那么就需要判断超时
			if errors.Is(err, context.DeadlineExceeded) == false {
				return -1, "Error"
			} else {
				// 超时了，那么就返回错误
				return -1, "time out"
			}
		}
	}
	speedResult, status := int(float32(elapsed.Nanoseconds())/1e6), res.Status
	defer res.Body.Close()

	if appSettings.TestUrlStatusCode != 0 {
		// 需要判断
		if res == nil || res.StatusCode != appSettings.TestUrlStatusCode {
			return -1, "Error statusCode"
		}
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, "Error res.Body"
	}

	pageHtmlString := string(content)

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

	if succeedWords != nil && len(succeedWords) > 0 {
		for _, word := range succeedWords {

			if strings.Contains(strings.ToLower(pageHtmlString), word) == true {
				return speedResult, "OK"
			}
		}
		return -1, "Error SucceedWord not found"
	}

	return speedResult, status
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
	testUrl string,
	succeedWords []string,
	browser *rod.Browser,
	localProxyHttpPort int) (int, string) {

	start := time.Now()
	page, e, err := rod_helper.NewPageNavigateWithProxy(browser,
		fmt.Sprintf("http://127.0.0.1:%d", localProxyHttpPort),
		testUrl, time.Duration(appSettings.OneNodeTestTimeOut)*time.Second)
	defer func() {
		if page != nil {
			_ = page.Close()
		}
	}()

	if err != nil {

		if len(appSettings.TestUrlSucceedWords) > 0 {
			// 无需判断超时，因为后续需要判断页面成功关键词
			if errors.Is(err, context.DeadlineExceeded) == false {
				// 不是超时错误，那么就返回错误，跳过
				return -1, "Error"
			}
		} else {
			// 因为没有设置成功关键词，那么就需要判断超时
			if errors.Is(err, context.DeadlineExceeded) == false {
				return -1, "Error"
			} else {
				// 超时了，那么就返回错误
				return -1, "time out"
			}
		}
	}

	elapsed := time.Since(start)
	speedResult := int(float32(elapsed.Nanoseconds()) / 1e6)

	pageHtmlString, err := page.HTML()
	if err != nil {
		return -1, "Error page.HTML " + err.Error()
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

	if succeedWords != nil && len(succeedWords) > 0 {
		for _, word := range succeedWords {

			if strings.Contains(strings.ToLower(pageHtmlString), word) == true {
				return speedResult, "OK"
			}
		}
		return -1, "Error SucceedWord not found"
	}

	return speedResult, "OK"
}

package xray_aio

import (
	"context"
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/csf-supplier/pkg"
	"github.com/allanpk716/rod_helper"
	"github.com/allanpk716/xray_pool/internal/pkg/settings"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pkg/errors"
	"github.com/ysmood/gson"
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
		if appSettings.TestUrlSucceedWordsEnable == true && len(appSettings.TestUrlSucceedWords) > 0 {
			// 无需判断超时，因为后续需要判断页面成功关键词
			if errors.Is(err, context.DeadlineExceeded) == false {
				// 不是超时错误，那么就返回错误，跳过
				return -1, "Error"
			}
			// 超时就通过，继续判断页面成功关键词
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

	if appSettings.TestUrlFailedWordsEnable == true {

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
	}

	if appSettings.TestUrlSucceedWordsEnable == true {

		if succeedWords != nil && len(succeedWords) > 0 {
			for _, word := range succeedWords {

				if strings.Contains(strings.ToLower(pageHtmlString), word) == true {
					return speedResult, "OK"
				}
			}
			return -1, "Error SucceedWord not found"
		}
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
	browserInfo *rod_helper.BrowserInfo,
	localProxyHttpPort int) (int, string) {

	start := time.Now()

	var err error
	var page *rod.Page
	nowHttpProxyUrl := fmt.Sprintf("http://127.0.0.1:%d", localProxyHttpPort)
	timeOut := time.Duration(appSettings.OneNodeTestTimeOut) * time.Second

	opt := pkg.NewHttpClientOptions(timeOut, nowHttpProxyUrl, "")
	client, err := pkg.NewHttpClient(opt)
	if err != nil {
		return -1, err.Error()
	}
	page, err = rod_helper.NewPage(browserInfo.Browser)
	if err != nil {
		return -1, err.Error()
	}
	defer func() {
		if page != nil {
			_ = page.Close()
		}
	}()
	err = page.SetWindow(&proto.BrowserBounds{
		Left:        gson.Int(0),
		Top:         gson.Int(50),
		Width:       gson.Int(900),
		Height:      gson.Int(900),
		WindowState: proto.BrowserWindowStateNormal,
	})
	router := rod_helper.NewPageHijackRouter(page, true, client.GetClient())
	defer func() {
		_ = router.Stop()
	}()
	go router.Run()
	page, e, err := rod_helper.PageNavigate(
		page, testUrl,
		timeOut,
	)
	if err != nil {
		// 这里可能会出现超时，但是实际上是成功的，所以这里不需要返回错误
		if errors.Is(err, context.DeadlineExceeded) == false {
			// 不是超时错误，那么就返回错误，跳过
			return -1, err.Error()
		}
	}
	err = page.Timeout(timeOut).WaitLoad()
	if err != nil {
		// 这里可能会出现超时，但是实际上是成功的，所以这里不需要返回错误
		if errors.Is(err, context.DeadlineExceeded) == false {
			// 不是超时错误，那么就返回错误，跳过
			return -1, err.Error()
		}
	}
	// ------------------判断返回值是否符合期望------------------
	logger.Infoln("PageStatusCodeCheck: ", testUrl)
	statusCode := rod_helper.StatusCodeInfo{
		Codes:          []int{403},
		Operator:       rod_helper.Match,
		WillDo:         rod_helper.Skip,
		NeedPunishment: false,
	}
	StatusCodeCheck, err := rod_helper.PageStatusCodeCheck(
		e,
		[]rod_helper.StatusCodeInfo{statusCode})
	if err != nil {
		return -1, err.Error()
	}
	switch StatusCodeCheck {
	case rod_helper.Skip:
		// 跳过后续的逻辑，不需要再次访问
		return -1, "StatusCodeCheck Error"
	}
	// 激活界面
	_, err = page.Activate()
	if err != nil {
		err = errors.New("Activate Error: " + err.Error())
		return -1, err.Error()
	}
	// ------------------会循环检测是否加载完毕，关键 Ele 出现即可------------------
	const manPageKeyWordXPath = "/html/body/div[2]/div/div/div/div[2]/div[1]/div[1]"
	logger.Infoln("HasPageLoaded: ", testUrl)
	pageLoaded := rod_helper.HasPageLoaded(page, manPageKeyWordXPath, appSettings.OneNodeTestTimeOut)
	logger.Infoln("HasPageLoaded: ", testUrl, pageLoaded)
	// 要在 StatusCode 检查之后再判断
	if pageLoaded == false {
		return -1, "PageLoaded Error"
	}

	elapsed := time.Since(start)
	speedResult := int(float32(elapsed.Nanoseconds()) / 1e6)
	// ------------------是否包含成功关键词------------------
	if appSettings.TestUrlSucceedWordsEnable == true {

		pageHtmlString, err := page.Timeout(timeOut).HTML()
		if err != nil {
			return -1, "Error page.HTML " + err.Error()
		}

		if succeedWords != nil && len(succeedWords) > 0 {
			for _, word := range succeedWords {

				if strings.Contains(strings.ToLower(pageHtmlString), word) == true {
					return speedResult, "OK"
				}
			}
			return -1, "Error SucceedWord not found"
		}
	}

	return speedResult, "OK"
}

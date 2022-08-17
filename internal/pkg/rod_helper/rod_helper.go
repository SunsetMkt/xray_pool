package rod_helper

import (
	"crypto/tls"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

func NewBrowser() *rod.Browser {

	nowUserData := filepath.Join(pkg.GetTmpFolderFPath(), pkg.RandStringBytesMaskImprSrcSB(20))
	purl := launcher.New().
		UserDataDir(nowUserData).
		MustLaunch()

	return rod.New().ControlURL(purl).MustConnect()
}

func NewPageNavigate(browser *rod.Browser, proxyUrl, desURL string, timeOut time.Duration) (*rod.Page, int, string, error) {

	page, err := newPage(browser)
	if err != nil {
		return nil, 0, "", err
	}

	return PageNavigate(page, proxyUrl, desURL, timeOut)
}

func PageNavigate(page *rod.Page, proxyUrl, desURL string, timeOut time.Duration) (*rod.Page, int, string, error) {

	router := page.HijackRequests()
	defer router.Stop()

	router.MustAdd("*", func(ctx *rod.Hijack) {
		px, _ := url.Parse(proxyUrl)
		err := ctx.LoadResponse(&http.Client{
			Transport: &http.Transport{
				Proxy:           http.ProxyURL(px),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}, true)
		if err != nil {
			return
		}
	})
	go router.Run()

	err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: pkg.RandomUserAgent(true),
	})
	if err != nil {
		if page != nil {
			page.Close()
		}
		return nil, 0, "", err
	}

	var e proto.NetworkResponseReceived
	wait := page.WaitEvent(&e)
	page = page.Timeout(timeOut)
	err = rod.Try(func() {
		page.MustNavigate(desURL)
		wait()
	})
	if err != nil {
		if page != nil {
			page.Close()
		}
		return nil, 0, "", err
	}

	// 出去前把 TimeOUt 取消了
	page = page.CancelTimeout()

	Status := e.Response.Status
	ResponseURL := e.Response.URL

	return page, Status, ResponseURL, nil
}

func newPage(browser *rod.Browser) (*rod.Page, error) {
	page, err := browser.Page(proto.TargetCreateTarget{URL: ""})
	if err != nil {
		return nil, err
	}
	return page, err
}

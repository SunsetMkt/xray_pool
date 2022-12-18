package subscribe

import (
	"fmt"
	"github.com/allanpk716/rod_helper"
	"github.com/go-resty/resty/v2"
	"time"
)

// GetByHTTPProxy 通过http代理访问网站
func GetByHTTPProxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*resty.Response, error) {

	opt := rod_helper.NewHttpClientOptions(timeOut)
	opt.SetHttpProxy(fmt.Sprintf("http://%s:%d", proxyAddress, proxyPort))
	client, err := rod_helper.NewHttpClient(opt)
	if err != nil {
		return nil, err
	}
	return client.R().Get(objUrl)
}

// GetBySocks5Proxy 通过Socks5代理访问网站
func GetBySocks5Proxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*resty.Response, error) {

	opt := rod_helper.NewHttpClientOptions(timeOut)
	opt.SetSocks5Proxy(fmt.Sprintf("socks5://%s:%d", proxyAddress, proxyPort))
	client, err := rod_helper.NewHttpClient(opt)
	if err != nil {
		return nil, err
	}
	return client.R().Get(objUrl)

}

// GetNoProxy 不通过代理访问网站
func GetNoProxy(url string, timeOut time.Duration) (*resty.Response, error) {

	opt := rod_helper.NewHttpClientOptions(timeOut)
	client, err := rod_helper.NewHttpClient(opt)
	if err != nil {
		return nil, err
	}
	return client.R().Get(url)
}

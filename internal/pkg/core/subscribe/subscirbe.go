package subscribe

import (
	"crypto/md5"
	"fmt"
	"github.com/WQGroup/logger"
	"net/http"
	"net/url"
)

type Subscribe struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Using bool   `json:"using"`
}

func NewSubscribe(inUrl, name string) *Subscribe {

	remarkName := ""
	if name == "" {
		u, err := url.Parse(inUrl)
		if err != nil {
			logger.Errorf("订阅[%s]解析失败:%s", inUrl, err.Error())
			remarkName = "remark"
		} else {
			remarkName = u.Host
		}
	} else {
		remarkName = name
	}

	return &Subscribe{Name: remarkName, Url: inUrl, Using: true}
}

// ID 返回订阅的ID，由 URL 的 MD5 值组成
func (s *Subscribe) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s.Url)))
}

func (s *Subscribe) UpdateNode(opt *UpdateOption) []string {
	var res *http.Response
	var err error
	switch opt.ProxyMode {
	case SOCKS:
		res, err = GetBySocks5Proxy(s.Url, opt.Addr, opt.Port, opt.Timeout)
	case HTTP:
		res, err = GetByHTTPProxy(s.Url, opt.Addr, opt.Port, opt.Timeout)
	default:
		res, err = GetNoProxy(s.Url, opt.Timeout)
	}
	if err != nil {
		logger.Error(err)
		return []string{}
	}
	logger.Info("访问 [", s.Url, "] -- ", res.Status)
	text := ReadDate(res)
	res.Body.Close()
	return Sub2links(text)
}

package subscribe

import (
	"crypto/md5"
	"fmt"
	"github.com/WQGroup/logger"
	"net/http"
)

type Subscribe struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Using bool   `json:"using"`
}

func NewSubscribe(url, name string) *Subscribe {
	return &Subscribe{Name: name, Url: url, Using: true}
}

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

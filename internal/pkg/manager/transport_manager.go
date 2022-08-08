package manager

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type TransportManager struct {
	index     int        // 内置的 transport index
	locker    sync.Mutex // 内置的锁
	httpPorts []int      // http 端口
}

func NewTransportManager(httpPorts []int) *TransportManager {
	t := &TransportManager{httpPorts: httpPorts}
	return t
}

func (t *TransportManager) ProxyMaker(req *http.Request) (*url.URL, error) {

	t.locker.Lock()
	defer func() {
		t.index++
		t.locker.Unlock()
	}()

	if t.index >= len(t.httpPorts) {
		t.index = 0
	}

	if !useProxy(canonicalAddr(req.URL)) {
		return nil, nil
	}
	proxy := fmt.Sprintf("http://127.0.0.1:%d", t.httpPorts[t.index])
	proxyURL, err := url.Parse(proxy)
	if err != nil || proxyURL.Scheme == "" {
		if u, err := url.Parse("http://" + proxy); err == nil {
			proxyURL = u
			err = nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	}
	return proxyURL, nil
}

// canonicalAddr returns url.Host but always with a ":port" suffix
func canonicalAddr(url *url.URL) string {
	addr := url.Host
	if !hasPort(addr) {
		return addr + ":" + portMap[url.Scheme]
	}
	return addr
}

func useProxy(addr string) bool {
	if len(addr) == 0 {
		return true
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return false
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() {
			return false
		}
	}

	//no_proxy := getenvEitherCase("NO_PROXY")
	//if no_proxy == "*" {
	//	return false
	//}
	no_proxy := ""

	addr = strings.ToLower(strings.TrimSpace(addr))
	if hasPort(addr) {
		addr = addr[:strings.LastIndex(addr, ":")]
	}

	for _, p := range strings.Split(no_proxy, ",") {
		p = strings.ToLower(strings.TrimSpace(p))
		if len(p) == 0 {
			continue
		}
		if hasPort(p) {
			p = p[:strings.LastIndex(p, ":")]
		}
		if addr == p || (p[0] == '.' && (strings.HasSuffix(addr, p) || addr == p[1:])) {
			return false
		}
	}
	return true
}

func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

var portMap = map[string]string{
	"http":  "80",
	"https": "443",
}

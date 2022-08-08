package manager

import (
	"context"
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/transport"
	"net/http"
	"time"
)

func (m *Manager) ReverseProxyStart() {

	m.reverseServerLocker.Lock()
	defer func() {
		m.reverseServerLocker.Unlock()
	}()

	if m.reverseServerRunning == true {
		logger.Debugln("Reverse Http Server is already running")
		return
	}
	m.reverseServerRunning = true

	socksPorts, httpPorts := m.GetOpenedProxyPorts()
	if len(socksPorts) == 0 && len(httpPorts) == 0 {
		logger.Panic("ReverseProxyStart: no open ports to proxy")
	}
	// 如果不满足，那么就再次扫描一个端口段，找到一个可用的端口给反向代理服务器
	alivePorts := pkg.ScanAlivePortList("63200-63400")
	if len(alivePorts) == 0 {
		logger.Panic("ReverseProxyStart: no open ports to proxy")
	} else {
		m.reverseServerHttpPort = alivePorts[0]
	}
	// 创建 transport 代理管理者
	m.transportManager = NewTransportManager(httpPorts)
	// 创建反向代理服务器
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {

			tr := transport.Transport{Proxy: m.transportManager.ProxyMaker}
			ctx.UserData, resp, err = tr.DetailedRoundTrip(req)
			return
		})
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		return resp
	})

	m.reverseServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.reverseServerHttpPort),
		Handler: proxy,
	}

	go func() {
		logger.Infoln("Try Start Reverse Http Server At Port", m.reverseServerHttpPort)
		if err := m.reverseServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorln("Start Reverse Server Error:", err)
		}
	}()
}

func (m *Manager) ReverseProxyStop() {

	m.reverseServerLocker.Lock()
	defer func() {
		m.reverseServerLocker.Unlock()
	}()
	if m.reverseServerRunning == false {
		logger.Debugln("Reverse Http Server is not running")
		return
	}
	m.reverseServerRunning = false

	exitOk := make(chan interface{}, 1)
	defer close(exitOk)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		if err := m.reverseServer.Shutdown(ctx); err != nil {
			logger.Errorln("Reverse Http Server Shutdown:", err)
		}
		exitOk <- true
	}()
	select {
	case <-ctx.Done():
		logger.Warningln("Reverse Http Server Shutdown timeout of 5 seconds.")
	case <-exitOk:
		logger.Infoln("Reverse Http Server Shutdown Successfully")
	}
	logger.Infoln("Reverse Http Server Shutdown Done.")
}

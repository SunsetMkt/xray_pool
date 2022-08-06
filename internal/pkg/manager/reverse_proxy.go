package manager

import (
	"github.com/WQGroup/logger"
	"github.com/valyala/fasthttp"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	"strconv"
)

func (m *Manager) ReverseProxyStart() {

	socksPorts, httpPorts := m.GetOpenProxyPorts()
	if len(socksPorts) == 0 && len(httpPorts) == 0 {
		logger.Panic("ReverseProxyStart: no open ports to proxy")
	}
	weights := make(map[string]proxy.Weight)
	for _, port := range httpPorts {
		weights["localhost:"+strconv.Itoa(port)] = 20
	}
	//weights := map[string]proxy.Weight{
	//	"localhost:9090": 20,
	//	"localhost:9091": 30,
	//	"localhost:9092": 50,
	//}
	m.reverseProxy = proxy.NewReverseProxy("", proxy.WithBalancer(weights))

	server := fasthttp.Server{
		Name:    "ReverseProxy Server",
		Handler: m.proxyHandler,
	}
	//server.Shutdown()

	if err := fasthttp.ListenAndServe(":8081", m.proxyHandler); err != nil {
		logger.Panicln("ReverseProxyStart:", err)
	}
	logger.Infoln("ReverseProxyStart", "success at:")
}

// proxyHandler ... fasthttp.RequestHandler func
func (m *Manager) proxyHandler(ctx *fasthttp.RequestCtx) {
	// all proxy to localhost
	m.reverseProxy.ServeHTTP(ctx)
}

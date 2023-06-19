package xray_aio

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core/routing"
	"github.com/allanpk716/xray_pool/internal/pkg/protocols"
	"github.com/allanpk716/xray_pool/internal/pkg/protocols/field"
	"path/filepath"
	"strings"
)

// GenConfigMix 构建 Xray 的运行配置，针对最终结果一次性启动所有节点，一个进程
func (x *XrayAIO) GenConfigMix() string {

	// 需要根据 base xray 的配置生成当前 Index xray 的配置
	path := filepath.Join(pkg.GetTmpFolderFPath(), configFileNameMix)

	var inbounds = make([]interface{}, 0)
	for i := range x.nodes {

		inputHttpPort := 0
		if x.AppSettings.XrayOpenSocksAndHttp == true {
			inputHttpPort = x.httpPorts[i]
		}
		one := x.inboundsConfig(i, x.socksPorts[i], inputHttpPort)
		inbounds = append(inbounds, one.([]interface{})...)
	}

	var conf = map[string]interface{}{
		//"log":       x.logConfig(),
		"inbounds":  inbounds,
		"outbounds": x.outboundConfig(),
		"policy":    x.policyConfig(),
		"dns":       x.dnsConfig(),
		"routing":   x.routingConfig(),
	}
	err := pkg.WriteJSON(conf, path)
	if err != nil {
		logger.Panicln("write config file error:", err.Error())
	}
	return path
}

func (x *XrayAIO) GetLogFPath() string {

	return filepath.Join(pkg.GetTmpFolderFPath(), xrayLogFileNameMix)
}

// logConfig 日志
func (x *XrayAIO) logConfig() interface{} {

	return map[string]string{
		"access":   x.GetLogFPath(),
		"loglevel": "warning",
	}
}

func (x *XrayAIO) getAllProxyTag() string {

	//tagStr := ""
	//for i := 0; i < len(x.nodes); i++ {
	//	tagStr += "proxy" + fmt.Sprintf("%d", i)
	//	if i != len(x.nodes)-1 {
	//		tagStr += ","
	//	}
	//}
	tagStr := fmt.Sprintf(outboundTag, 0)

	return tagStr
}

func (x *XrayAIO) getInboundTags(nodeIndex int) []interface{} {

	outInboundStr := make([]interface{}, 0)
	// 如果开启了 http ，那么就是两个 inbound，否则就是一个
	outInboundStr = append(outInboundStr, fmt.Sprintf(inboundSocksTag, nodeIndex))
	if x.AppSettings.XrayOpenSocksAndHttp == true {
		outInboundStr = append(outInboundStr, fmt.Sprintf(inboundHttpTag, nodeIndex))
	}
	return outInboundStr
}

func (x *XrayAIO) getOutboundTags(nodeIndex int) string {

	return fmt.Sprintf(outboundTag, nodeIndex)
}

// inboundsConfig 入站
func (x *XrayAIO) inboundsConfig(tagIndex int, SocksPort, HttpPort int) interface{} {
	listen := "127.0.0.1"
	if x.AppSettings.MainProxySettings.AllowLanConn {
		listen = "0.0.0.0"
	}
	data := []interface{}{
		map[string]interface{}{
			"tag":      fmt.Sprintf(inboundSocksTag, tagIndex),
			"port":     SocksPort,
			"listen":   listen,
			"protocol": "socks",
			"sniffing": map[string]interface{}{
				"enabled": x.AppSettings.MainProxySettings.Sniffing,
				"destOverride": []string{
					"http",
					"tls",
				},
			},
			"settings": map[string]interface{}{
				"auth":      "noauth",
				"udp":       x.AppSettings.MainProxySettings.RelayUDP,
				"userLevel": 0,
			},
		},
	}
	if HttpPort > 0 {
		data = append(data, map[string]interface{}{
			"tag":      fmt.Sprintf(inboundHttpTag, tagIndex),
			"port":     HttpPort,
			"listen":   listen,
			"protocol": "http",
			"settings": map[string]interface{}{
				"userLevel": 0,
			},
		})
	}
	if x.AppSettings.MainProxySettings.DNSPort > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "dns-in",
			"port":     x.AppSettings.MainProxySettings.DNSPort,
			"listen":   listen,
			"protocol": "dokodemo-door",
			"settings": map[string]interface{}{
				"userLevel": 0,
				"address":   x.AppSettings.MainProxySettings.DNSForeign,
				"network":   "tcp,udp",
				"port":      53,
			},
		})
	}
	return data
}

// policyConfig 本地策略
func (x *XrayAIO) policyConfig() interface{} {
	return map[string]interface{}{
		"levels": map[string]interface{}{
			"0": map[string]interface{}{
				"handshake":    4,
				"connIdle":     300,
				"uplinkOnly":   1,
				"downlinkOnly": 1,
				"bufferSize":   10240,
			},
		},
		"system": map[string]interface{}{
			"statsInboundUplink":   true,
			"statsInboundDownlink": true,
		},
	}
}

// dnsConfig DNS
func (x *XrayAIO) dnsConfig() interface{} {
	servers := make([]interface{}, 0)
	if x.AppSettings.MainProxySettings.DNSDomestic != "" {
		servers = append(servers, map[string]interface{}{
			"address": x.AppSettings.MainProxySettings.DNSDomestic,
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if x.AppSettings.MainProxySettings.DNSDomesticBackup != "" {
		servers = append(servers, map[string]interface{}{
			"address": x.AppSettings.MainProxySettings.DNSDomesticBackup,
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if x.AppSettings.MainProxySettings.DNSForeign != "" {
		servers = append(servers, map[string]interface{}{
			"address": x.AppSettings.MainProxySettings.DNSForeign,
			"port":    53,
			"domains": []interface{}{
				"geosite:geolocation-!cn",
				"geosite:speedtest",
			},
		})
	}
	return map[string]interface{}{
		"hosts": map[string]interface{}{
			"domain:googleapis.cn": "googleapis.com",
		},
		"servers": servers,
	}
}

// routingConfig 路由
func (x *XrayAIO) routingConfig() interface{} {
	rules := make([]interface{}, 0)

	// 生成对应的 inbound 和 outbound 的对应的路由关系
	for i := 0; i < len(x.nodes); i++ {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"inboundTag":  x.getInboundTags(i),
			"outboundTag": x.getOutboundTags(i),
		})
	}

	if x.AppSettings.MainProxySettings.DNSPort != 0 {
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"inboundTag": []interface{}{
				"dns-in",
			},
			"outboundTag": "dns-out",
		})
	}
	if x.AppSettings.MainProxySettings.DNSForeign != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": x.getAllProxyTag(),
			"ip": []string{
				x.AppSettings.MainProxySettings.DNSForeign,
			},
		})
	}
	if x.AppSettings.MainProxySettings.DNSDomestic != "" || x.AppSettings.MainProxySettings.DNSDomesticBackup != "" {
		var ip []string
		if x.AppSettings.MainProxySettings.DNSDomestic != "" {
			ip = append(ip, x.AppSettings.MainProxySettings.DNSDomestic)
		}
		if x.AppSettings.MainProxySettings.DNSDomesticBackup != "" {
			ip = append(ip, x.AppSettings.MainProxySettings.DNSDomesticBackup)
		}
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "direct",
			"ip":          ip,
		})
	}
	ips, domains := x.route.GetRulesGroupData(routing.TypeBlock)
	if len(ips) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "block",
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "block",
			"domain":      domains,
		})
	}
	ips, domains = x.route.GetRulesGroupData(routing.TypeDirect)
	if len(ips) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"domain":      domains,
		})
	}
	// 这里需要根据 xray 的 targetUrl 进行 临时的规则附加设置
	ips, domains = x.route.GetRulesGroupData(routing.TypeProxy)
	if x.targetUrl != "" {
		domains = append(domains, x.targetUrl)
	}
	if len(ips) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": x.getAllProxyTag(),
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": x.getAllProxyTag(),
			"domain":      domains,
		})
	}

	if x.AppSettings.MainProxySettings.BypassLANAndMainLand {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"ip": []string{
				"geoip:private",
				"geoip:cn",
			},
		})
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"domain": []string{
				"geosite:cn",
			},
		})
	}
	return map[string]interface{}{
		"domainStrategy": x.AppSettings.MainProxySettings.RoutingStrategy,
		"rules":          rules,
	}
}

// outboundConfig 出站
func (x *XrayAIO) outboundConfig() interface{} {

	out := make([]interface{}, 0)

	for i, nowNode := range x.nodes {

		p := nowNode.Protocol
		switch p.GetProtocolMode() {
		case protocols.ModeTrojan:
			t := p.(*protocols.Trojan)
			out = append(out, x.trojanOutbound(i, t))
		case protocols.ModeShadowSocks:
			ss := p.(*protocols.ShadowSocks)
			out = append(out, x.shadowSocksOutbound(i, ss))
		case protocols.ModeVMess:
			v := p.(*protocols.VMess)
			out = append(out, x.vMessOutbound(i, v))
		case protocols.ModeSocks:
			v := p.(*protocols.Socks)
			out = append(out, x.socksOutbound(i, v))
		case protocols.ModeVLESS:
			v := p.(*protocols.VLess)
			out = append(out, x.vLessOutbound(i, v))
		case protocols.ModeVMessAEAD:
			v := p.(*protocols.VMessAEAD)
			out = append(out, x.vMessAEADOutbound(i, v))
		}
	}

	out = append(out, map[string]interface{}{
		"tag":      "direct",
		"protocol": "freedom",
		"settings": map[string]interface{}{},
	})
	out = append(out, map[string]interface{}{
		"tag":      "block",
		"protocol": "blackhole",
		"settings": map[string]interface{}{
			"response": map[string]interface{}{
				"type": "http",
			},
		},
	})
	out = append(out, map[string]interface{}{
		"tag":      "dns-out",
		"protocol": "dns",
	})
	return out
}

// ShadowSocks 出站
func (x *XrayAIO) shadowSocksOutbound(tagIndex int, ss *protocols.ShadowSocks) interface{} {

	defaultOut := map[string]interface{}{
		"tag":      fmt.Sprintf(outboundTag, tagIndex),
		"protocol": "shadowsocks",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  ss.Address,
					"port":     ss.Port,
					"password": ss.Password,
					"method":   ss.Method,
					"level":    0,
				},
			},
		},
		"streamSettings": map[string]interface{}{
			"network": "tcp",
		},
	}
	// 如果 RequestHost 不为空
	if ss.RequestHost != "" {

		tcpSettings := map[string]interface{}{
			"header": map[string]interface{}{
				"type": "http",
				"request": map[string]interface{}{
					"version": "1.1",
					"method":  "GET",
					"path":    []string{"/"},
					"headers": map[string]interface{}{
						"Host": []string{
							ss.RequestHost,
						},
						"User-Agent": []string{
							"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
							"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/53.0.2785.109 Mobile/14A456 Safari/601.1.46",
						},
						"Accept-Encoding": []string{
							"gzip, deflate",
						},
						"Connection": []string{
							"keep-alive",
						},
						"Pragma": "no-cache",
					},
				},
			},
		}

		streamSettings := map[string]interface{}{
			"network":     "tcp",
			"tcpSettings": tcpSettings,
		}

		defaultOut["streamSettings"] = streamSettings
	}

	return defaultOut
}

// Trojan 出站
func (x *XrayAIO) trojanOutbound(tagIndex int, trojan *protocols.Trojan) interface{} {
	streamSettings := map[string]interface{}{
		"network":  "tcp",
		"security": "tls",
	}
	if trojan.Sni() != "" {
		streamSettings["tlsSettings"] = map[string]interface{}{
			"allowInsecure": false,
			"serverName":    trojan.Sni(),
		}
	}
	return map[string]interface{}{
		"tag":      fmt.Sprintf(outboundTag, tagIndex),
		"protocol": "trojan",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  trojan.Address,
					"port":     trojan.Port,
					"password": trojan.Password,
					"level":    0,
				},
			},
		},
		"streamSettings": streamSettings,
	}
}

// VMess 出站
func (x *XrayAIO) vMessOutbound(tagIndex int, vmess *protocols.VMess) interface{} {
	mux := x.AppSettings.MainProxySettings.Mux
	streamSettings := map[string]interface{}{
		"network":  vmess.Net,
		"security": vmess.Tls,
	}
	if vmess.Tls == "tls" {
		tlsSettings := map[string]interface{}{
			"allowInsecure": false,
		}
		if vmess.Sni != "" {
			tlsSettings["serverName"] = vmess.Sni
		}
		if vmess.Alpn != "" {
			tlsSettings["alpn"] = strings.Split(vmess.Alpn, ",")
		}
		streamSettings["tlsSettings"] = tlsSettings
	}
	switch vmess.Net {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
		if vmess.Type != "none" {
			kcpSettings["seed"] = vmess.Path
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vmess.Path,
			"headers": map[string]interface{}{
				"Host": vmess.Host,
			},
		}
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vmess.Host, ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vmess.Path,
			"host": host,
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vmess.Host,
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
		if vmess.Host != "none" {
			quicSettings["key"] = vmess.Path
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vmess.Path,
			"multiMode":   vmess.Type == "multi",
		}
	}
	return map[string]interface{}{
		"tag":      fmt.Sprintf(outboundTag, tagIndex),
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vmess.Add,
					"port":    vmess.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":       vmess.Id,
							"alterId":  vmess.Aid,
							"security": vmess.Scy,
							"level":    0,
						},
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// socks 出站
func (x *XrayAIO) socksOutbound(tagIndex int, socks *protocols.Socks) interface{} {
	user := map[string]interface{}{
		"address": socks.Address,
		"port":    socks.Port,
	}
	if socks.Username != "" || socks.Password != "" {
		user["users"] = map[string]interface{}{
			"user": socks.Username,
			"pass": socks.Password,
		}
	}
	return map[string]interface{}{
		"tag":      fmt.Sprintf(outboundTag, tagIndex),
		"protocol": "socks",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				user,
			},
		},
		"streamSettings": map[string]interface{}{
			"network": "tcp",
			"tcpSettings": map[string]interface{}{
				"header": map[string]interface{}{
					"type": "none",
				},
			},
		},
		"mux": map[string]interface{}{
			"enabled": false,
		},
	}
}

// VLESS 出站
func (x *XrayAIO) vLessOutbound(tagIndex int, vless *protocols.VLess) interface{} {
	mux := x.AppSettings.MainProxySettings.Mux
	security := vless.GetValue(field.Security)
	network := vless.GetValue(field.NetworkType)
	user := map[string]interface{}{
		"id":         vless.ID,
		"encryption": vless.GetValue(field.VLessEncryption),
		"level":      0,
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		tlsSettings := map[string]interface{}{
			"allowInsecure": false,
		}
		sni := vless.GetHostValue(field.SNI)
		alpn := vless.GetValue(field.Alpn)
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if alpn != "" {
			tlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		streamSettings["tlsSettings"] = tlsSettings
	case "xtls":
		xtlsSettings := map[string]interface{}{
			"allowInsecure": false,
		}
		sni := vless.GetHostValue(field.SNI)
		alpn := vless.GetValue(field.Alpn)
		if sni != "" {
			xtlsSettings["serverName"] = sni
		}
		if alpn != "" {
			xtlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		streamSettings["xtlsSettings"] = xtlsSettings
		user["flow"] = vless.GetValue(field.Flow)
		mux = false
	}
	switch network {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vless.GetValue(field.TCPHeaderType),
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vless.GetValue(field.MkcpHeaderType),
			},
		}
		if vless.Has(field.Seed.Key) {
			kcpSettings["seed"] = vless.GetValue(field.Seed)
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vless.GetHostValue(field.H2Host), ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vless.GetValue(field.H2Path),
			"host": host,
		}
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vless.GetValue(field.WsPath),
			"headers": map[string]interface{}{
				"Host": vless.GetValue(field.WsHost),
			},
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vless.GetValue(field.QuicSecurity),
			"header": map[string]interface{}{
				"type": vless.GetValue(field.QuicHeaderType),
			},
		}
		if vless.GetValue(field.QuicSecurity) != "none" {
			quicSettings["key"] = vless.GetValue(field.QuicKey)
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vless.GetValue(field.GrpcServiceName),
			"multiMode":   vless.GetValue(field.GrpcMode) == "multi",
		}
	}
	return map[string]interface{}{
		"tag":      fmt.Sprintf(outboundTag, tagIndex),
		"protocol": "vless",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vless.Address,
					"port":    vless.Port,
					"users": []interface{}{
						user,
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// VMessAEAD 出站
func (x *XrayAIO) vMessAEADOutbound(tagIndex int, vmess *protocols.VMessAEAD) interface{} {
	mux := x.AppSettings.MainProxySettings.Mux
	security := vmess.GetValue(field.Security)
	network := vmess.GetValue(field.NetworkType)
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		tlsSettings := map[string]interface{}{
			"allowInsecure": false,
		}
		sni := vmess.GetHostValue(field.SNI)
		alpn := vmess.GetValue(field.Alpn)
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if alpn != "" {
			tlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		streamSettings["tlsSettings"] = tlsSettings
	}
	switch network {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vmess.GetValue(field.TCPHeaderType),
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vmess.GetValue(field.MkcpHeaderType),
			},
		}
		if vmess.Has(field.Seed.Key) {
			kcpSettings["seed"] = vmess.GetValue(field.Seed)
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vmess.GetHostValue(field.H2Host), ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vmess.GetValue(field.H2Path),
			"host": host,
		}
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vmess.GetValue(field.WsPath),
			"headers": map[string]interface{}{
				"Host": vmess.GetValue(field.WsHost),
			},
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vmess.GetValue(field.QuicSecurity),
			"header": map[string]interface{}{
				"type": vmess.GetValue(field.QuicHeaderType),
			},
		}
		if vmess.GetValue(field.QuicSecurity) != "none" {
			quicSettings["key"] = vmess.GetValue(field.QuicKey)
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vmess.GetValue(field.GrpcServiceName),
			"multiMode":   vmess.GetValue(field.GrpcMode) == "multi",
		}
	}
	return map[string]interface{}{
		"tag":      fmt.Sprintf(outboundTag, tagIndex),
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vmess.Address,
					"port":    vmess.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":       vmess.ID,
							"security": vmess.GetValue(field.VMessEncryption),
							"level":    0,
						},
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

const (
	configFileName  = "xray_config_%d.json"
	xrayLogFileName = "xray_access_%d.log"

	configFileNameMix  = "xray_config_mix.json"
	xrayLogFileNameMix = "xray_access_mix.log"
)

const (
	inboundSocksTag = "in_socks_%d"
	inboundHttpTag  = "in_http_%d"
	outboundTag     = "out_bound_%d"
)

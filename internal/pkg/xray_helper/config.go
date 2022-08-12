package xray_helper

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

// GenConfig 构建 Xray 的运行配置
func (x *XrayHelper) GenConfig(node protocols.Protocol) string {

	// 需要根据 base xray 的配置生成当前 Index xray 的配置
	path := filepath.Join(pkg.GetTmpFolderFPath(), fmt.Sprintf(configFileName, x.index))
	var conf = map[string]interface{}{
		"log":       x.logConfig(),
		"inbounds":  x.inboundsConfig(),
		"outbounds": x.outboundConfig(node),
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

func (x *XrayHelper) GetLogFPath() string {
	path := filepath.Join(pkg.GetTmpFolderFPath(), fmt.Sprintf(xrayLogFileName, x.index))
	return path
}

// logConfig 日志
func (x *XrayHelper) logConfig() interface{} {

	return map[string]string{
		"access":   x.GetLogFPath(),
		"loglevel": "warning",
	}
}

// inboundsConfig 入站
func (x *XrayHelper) inboundsConfig() interface{} {
	listen := "127.0.0.1"
	if x.ProxySettings.AllowLanConn {
		listen = "0.0.0.0"
	}
	data := []interface{}{
		map[string]interface{}{
			"tag":      "proxy",
			"port":     x.ProxySettings.SocksPort,
			"listen":   listen,
			"protocol": "socks",
			"sniffing": map[string]interface{}{
				"enabled": x.ProxySettings.Sniffing,
				"destOverride": []string{
					"http",
					"tls",
				},
			},
			"settings": map[string]interface{}{
				"auth":      "noauth",
				"udp":       x.ProxySettings.RelayUDP,
				"userLevel": 0,
			},
		},
	}
	if x.ProxySettings.HttpPort > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "http",
			"port":     x.ProxySettings.HttpPort,
			"listen":   listen,
			"protocol": "http",
			"settings": map[string]interface{}{
				"userLevel": 0,
			},
		})
	}
	if x.ProxySettings.DNSPort > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "dns-in",
			"port":     x.ProxySettings.DNSPort,
			"listen":   listen,
			"protocol": "dokodemo-door",
			"settings": map[string]interface{}{
				"userLevel": 0,
				"address":   x.ProxySettings.DNSForeign,
				"network":   "tcp,udp",
				"port":      53,
			},
		})
	}
	return data
}

// policyConfig 本地策略
func (x *XrayHelper) policyConfig() interface{} {
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
func (x *XrayHelper) dnsConfig() interface{} {
	servers := make([]interface{}, 0)
	if x.ProxySettings.DNSDomestic != "" {
		servers = append(servers, map[string]interface{}{
			"address": x.ProxySettings.DNSDomestic,
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if x.ProxySettings.DNSDomesticBackup != "" {
		servers = append(servers, map[string]interface{}{
			"address": x.ProxySettings.DNSDomesticBackup,
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if x.ProxySettings.DNSForeign != "" {
		servers = append(servers, map[string]interface{}{
			"address": x.ProxySettings.DNSForeign,
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
func (x *XrayHelper) routingConfig() interface{} {
	rules := make([]interface{}, 0)
	if x.ProxySettings.DNSPort != 0 {
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"inboundTag": []interface{}{
				"dns-in",
			},
			"outboundTag": "dns-out",
		})
	}
	if x.ProxySettings.DNSForeign != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "proxy",
			"ip": []string{
				x.ProxySettings.DNSForeign,
			},
		})
	}
	if x.ProxySettings.DNSDomestic != "" || x.ProxySettings.DNSDomesticBackup != "" {
		var ip []string
		if x.ProxySettings.DNSDomestic != "" {
			ip = append(ip, x.ProxySettings.DNSDomestic)
		}
		if x.ProxySettings.DNSDomesticBackup != "" {
			ip = append(ip, x.ProxySettings.DNSDomesticBackup)
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
			"outboundTag": "proxy",
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "proxy",
			"domain":      domains,
		})
	}

	if x.ProxySettings.BypassLANAndMainLand {
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
		"domainStrategy": x.ProxySettings.RoutingStrategy,
		"rules":          rules,
	}
}

// outboundConfig 出站
func (x *XrayHelper) outboundConfig(n protocols.Protocol) interface{} {
	out := make([]interface{}, 0)
	switch n.GetProtocolMode() {
	case protocols.ModeTrojan:
		t := n.(*protocols.Trojan)
		out = append(out, x.trojanOutbound(t))
	case protocols.ModeShadowSocks:
		ss := n.(*protocols.ShadowSocks)
		out = append(out, x.shadowSocksOutbound(ss))
	case protocols.ModeVMess:
		v := n.(*protocols.VMess)
		out = append(out, x.vMessOutbound(v))
	case protocols.ModeSocks:
		v := n.(*protocols.Socks)
		out = append(out, x.socksOutbound(v))
	case protocols.ModeVLESS:
		v := n.(*protocols.VLess)
		out = append(out, x.vLessOutbound(v))
	case protocols.ModeVMessAEAD:
		v := n.(*protocols.VMessAEAD)
		out = append(out, x.vMessAEADOutbound(v))
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
func (x *XrayHelper) shadowSocksOutbound(ss *protocols.ShadowSocks) interface{} {
	return map[string]interface{}{
		"tag":      "proxy",
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
}

// Trojan 出站
func (x *XrayHelper) trojanOutbound(trojan *protocols.Trojan) interface{} {
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
		"tag":      "proxy",
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
func (x *XrayHelper) vMessOutbound(vmess *protocols.VMess) interface{} {
	mux := x.ProxySettings.Mux
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
		"tag":      "proxy",
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
func (x *XrayHelper) socksOutbound(socks *protocols.Socks) interface{} {
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
		"tag":      "proxy",
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
func (x *XrayHelper) vLessOutbound(vless *protocols.VLess) interface{} {
	mux := x.ProxySettings.Mux
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
		"tag":      "proxy",
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
func (x *XrayHelper) vMessAEADOutbound(vmess *protocols.VMessAEAD) interface{} {
	mux := x.ProxySettings.Mux
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
		"tag":      "proxy",
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
)

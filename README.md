# xray_pool

因为某些原因，需要在做爬虫任务的时候用到代理池（否则会在调用频率上容易被限制），显然是没打算额外去购买这样的服务，但是呢，鸡场一般是有的，也就是可以在 [XTLS/Xray-core](https://github.com/XTLS/Xray-core) 的核心基础上，进行一些使用上的调整，变相变成一个“伪代理池”来用。

本来是打算使用 [v2rayA](https://github.com/v2rayA/v2rayA) 的"负载均衡"来当“代理池”用的，但是测试下来，貌似没有能够真正的去均衡流量。又看到一个 [Txray](https://github.com/hsernos/Txray) 项目，应该就是想要的样子了，就打算根据这个项目进行魔改。

## 该软件适用于什么情况

1. 鸡场支持：VMess、Shadowsocks、Trojan、VLESS、VMessAEAD、Socks 协议
2. 鸡场不限制客户端数量，或者数量大于 5（太少没有意义）
3. 爬虫任务需要每一次 HTTP 请求就切换一个代理

## 软件界面

待实现

## 软件结构图

本程序可以理解为就是一个壳，调用两个核心程序来实现的（实现的代价最低，可维护性也最高）。

为了对外暴露一个负载均衡的 HTTP 代理端口（负载均衡策略是 round robin），所以使用了 [glider](https://github.com/nadoo/glider) 这个软件。

![image-20220809165339675](README.assets/image-20220809165339675.png)

## 如何使用

待实现

## API 接口

待实现

## Thanks

* [Txray](https://github.com/hsernos/Txray)
* [XTLS/Xray-core](https://github.com/XTLS/Xray-core)
* [glider](https://github.com/nadoo/glider)
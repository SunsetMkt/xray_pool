# xray_pool

因为某些原因，需要在做爬虫任务的时候用到代理池（否则会在调用频率上容易被限制），显然是没打算额外去购买这样的服务，但是呢，鸡场一般是有的，也就是可以在 [XTLS/Xray-core](https://github.com/XTLS/Xray-core) 的核心基础上，进行一些使用上的调整，变相变成一个“伪代理池”来用。

本来是打算使用 [v2rayA](https://github.com/v2rayA/v2rayA) 的"负载均衡"来当“代理池”用的，但是测试下来，貌似没有能够真正的去均衡流量。又看到一个 [Txray](https://github.com/hsernos/Txray) 项目，应该就是想要的样子了，就打算根据这个项目进行魔改。

## Thank

* [Txray](https://github.com/hsernos/Txray)
* [XTLS/Xray-core](https://github.com/XTLS/Xray-core)
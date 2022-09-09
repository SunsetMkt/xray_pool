package backend

import (
	"context"
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/frontend"
	v1 "github.com/allanpk716/xray_pool/internal/backend/controllers/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type BackEnd struct {
	running       bool
	srv           *http.Server
	locker        sync.Mutex
	restartSignal chan interface{}
	exitSignal    chan interface{}
}

func NewBackEnd(restartSignal, exitSignal chan interface{}) *BackEnd {
	return &BackEnd{restartSignal: restartSignal, exitSignal: exitSignal}
}

func (b *BackEnd) start() {

	defer b.locker.Unlock()
	b.locker.Lock()

	if b.running == true {
		logger.Debugln("Http Server is already running")
		return
	}
	b.running = true

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	engine := gin.Default()
	// 默认所有都通过
	engine.Use(cors.Default())

	// 本地启动前端调试的时候使用，替换为本地打开浏览器的地址即可
	//corsConfig := cors.DefaultConfig()
	//corsConfig.AllowOrigins = []string{"http://127.0.0.1:5173"}
	//corsConfig.AllowCredentials = true
	//engine.Use(cors.New(corsConfig))

	cbV1 := v1.NewControllerBase(b.restartSignal, b.exitSignal)
	// v1路由: /v1/xxx
	GroupV1 := engine.Group("/" + cbV1.GetVersion())
	{
		GroupV1.GET("/system-status", cbV1.SystemStatus)
		GroupV1.POST("/setup", cbV1.SetUp)
		GroupV1.POST("/login", cbV1.Login)
		GroupV1.GET("/proxy_list", cbV1.GetProxyListHandler)

		//GroupV1.Use(middle.CheckAuth())
		GroupV1.POST("/logout", cbV1.Logout)
		GroupV1.POST("/change_pwd", cbV1.ChangePWD)
		GroupV1.POST("/exit", cbV1.ExitHandler)
		GroupV1.POST("/clear_tmp_folder", cbV1.ClearTmpFolder)
		// 基础设置相关
		GroupV1.GET("/settings", cbV1.SettingsHandler)
		GroupV1.PUT("/settings", cbV1.SettingsHandler)
		GroupV1.GET("/def_settings", cbV1.DefSettingsHandler)
		// xray pool 启动、停止相关
		GroupV1.POST("/start_proxy_pool", cbV1.StartProxyPoolHandler)
		GroupV1.POST("/stop_proxy_pool", cbV1.StopProxyPoolHandler)
		// 订阅相关
		GroupV1.GET("/subscribe_list", cbV1.SubscribeListHandler)
		GroupV1.GET("/node_list", cbV1.NodesListHandler)
		GroupV1.POST("/add_subscribe", cbV1.SubscribeAddHandler)
		GroupV1.POST("/update_nodes", cbV1.SubscribeUpdateNodesHandler)
		GroupV1.POST("/update_subscribe", cbV1.SubscribeUpdateHandler)
		GroupV1.POST("/del_subscribe", cbV1.SubscribeDelHandler)
		// 路由规则相关
		GroupV1.POST("/routing_add", cbV1.RoutingAddHandler)
		GroupV1.GET("/routing_list", cbV1.RoutingListHandler)
		GroupV1.POST("/routing_delete", cbV1.RoutingDeleteHandler)
	}
	// -------------------------------------------------
	// 静态文件服务器，加载 html 页面
	engine.GET("/", func(c *gin.Context) {
		c.Header("content-type", "text/html;charset=utf-8")
		c.String(http.StatusOK, string(frontend.SpaIndexHtml))
	})
	engine.StaticFS(frontend.SpaRelativePath, frontend.Assets(frontend.SpaFolderName, frontend.SpaJS))
	// -------------------------------------------------
	engine.Any("/api", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	// -------------------------------------------------
	// listen and serve on 0.0.0.0:8080(default)
	b.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cbV1.GetAppStartPort()),
		Handler: engine,
	}
	go func() {
		logger.Infoln("Try Start Http Server At Port", cbV1.GetAppStartPort())
		if err := b.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorln("Start Server Error:", err)
		}

		defer func() {
			cbV1.Close()
		}()
	}()
}

func (b *BackEnd) Restart() {

	stopFunc := func() {

		b.locker.Lock()
		defer func() {
			b.locker.Unlock()
		}()
		if b.running == false {
			logger.Debugln("Http Server is not running")
			return
		}
		b.running = false

		exitOk := make(chan interface{}, 1)
		defer close(exitOk)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func() {
			if err := b.srv.Shutdown(ctx); err != nil {
				logger.Errorln("Http Server Shutdown:", err)
			}
			exitOk <- true
		}()
		select {
		case <-ctx.Done():
			logger.Warningln("Http Server Shutdown timeout of 5 seconds.")
		case <-exitOk:
			logger.Infoln("Http Server Shutdown Successfully")
		}
		logger.Infoln("Http Server Shutdown Done.")
	}

	for {
		select {
		case <-b.restartSignal:
			{
				stopFunc()
				b.start()
			}
		case <-b.exitSignal:
			{
				stopFunc()
				logger.Infoln("Http Server Exit.")
				//os.Exit(0)
				return
			}
		}
	}
}

func (b *BackEnd) Close() {
	defer b.locker.Unlock()
	b.locker.Lock()

	b.exitSignal <- true
}

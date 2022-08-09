package backend

import (
	"context"
	"fmt"
	"github.com/WQGroup/logger"
	v1 "github.com/allanpk716/xray_pool/internal/backend/controllers/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type BackEnd struct {
	httpPort int
	running  bool
	srv      *http.Server

	locker        sync.Mutex
	restartSignal chan interface{}
}

func NewBackEnd(httpPort int, restartSignal chan interface{}) *BackEnd {
	return &BackEnd{httpPort: httpPort, restartSignal: restartSignal}
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

	cbV1 := v1.NewControllerBase(b.restartSignal)
	// v1路由: /v1/xxx
	GroupV1 := engine.Group("/" + cbV1.GetVersion())
	{
		//GroupV1.Use(middle.CheckAuth())
		GroupV1.POST("/start_proxy_pool", cbV1.StartProxyPoolHandler)
		GroupV1.POST("/stop_proxy_pool", cbV1.StopProxyPoolHandler)
		GroupV1.GET("/proxy_list", cbV1.GetProxyListHandler)
	}
	// -------------------------------------------------
	// 静态文件服务器，加载 html 页面
	//engine.GET("/", func(c *gin.Context) {
	//	c.Header("content-type", "text/html;charset=utf-8")
	//	c.String(http.StatusOK, string(dist.SpaIndexHtml))
	//})
	//engine.StaticFS(dist.SpaFolderJS, dist.Assets(dist.SpaFolderName+dist.SpaFolderJS, dist.SpaJS))
	//engine.StaticFS(dist.SpaFolderCSS, dist.Assets(dist.SpaFolderName+dist.SpaFolderCSS, dist.SpaCSS))
	//engine.StaticFS(dist.SpaFolderFonts, dist.Assets(dist.SpaFolderName+dist.SpaFolderFonts, dist.SpaFonts))
	//engine.StaticFS(dist.SpaFolderIcons, dist.Assets(dist.SpaFolderName+dist.SpaFolderIcons, dist.SpaIcons))
	//engine.StaticFS(dist.SpaFolderImages, dist.Assets(dist.SpaFolderName+dist.SpaFolderImages, dist.SpaImages))
	// -------------------------------------------------
	engine.Any("/api", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	// -------------------------------------------------
	// listen and serve on 0.0.0.0:8080(default)
	b.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", b.httpPort),
		Handler: engine,
	}
	go func() {
		logger.Infoln("Try Start Http Server At Port", b.httpPort)
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
		}
	}
}

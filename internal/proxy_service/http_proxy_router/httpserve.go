package http_proxy_router

import (
	"context"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/cert_file"
	mylog "github.com/captainlee1024/go-gateway/internal/proxy_service/log"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	// 5. 注册路由
	r := setUp() // 6. 启动服务（开启平滑下线）
	HttpSrvHandler = &http.Server{
		Addr:           fmt.Sprintf(":%v", settings.ConfProxy.HttpConfig.HTTPPort),
		Handler:        r,
		ReadTimeout:    time.Duration(settings.ConfProxy.HttpConfig.ReadTime) * time.Second,
		WriteTimeout:   time.Duration(settings.ConfProxy.HttpConfig.WriteTime) * time.Second,
		MaxHeaderBytes: 1 << uint(settings.ConfProxy.HttpConfig.MaxHeaderBytes),
	}

	fmt.Printf("[INFO] http_proxy_run:%d\n", settings.ConfProxy.HttpConfig.HTTPPort)
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		mylog.Log.Fatal("listen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
			"err": err,
		})
	}
}

func HttpServerStop() {
	shutdownTrace := mylog.NewTrace()
	mylog.Log.Info("Shutdown", shutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": fmt.Sprintf("Shutdown httpServer%s", HttpSrvHandler.Addr),
	})

	// 创建一个 5 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		mylog.Log.Fatal("Shutdown", shutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
	}

	mylog.Log.Info("Server exiting", shutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": fmt.Sprintf("httpServer%s exiting", HttpSrvHandler.Addr),
	})
}

func HttpsServerRun() {
	// 5. 注册路由
	r := setUp() // 6. 启动服务（开启平滑下线）
	HttpsSrvHandler = &http.Server{
		Addr:           fmt.Sprintf(":%v", settings.ConfProxy.HttpsConfig.HTTPSPort),
		Handler:        r,
		ReadTimeout:    time.Duration(settings.ConfProxy.HttpsConfig.ReadTime) * time.Second,
		WriteTimeout:   time.Duration(settings.ConfProxy.HttpsConfig.WriteTime) * time.Second,
		MaxHeaderBytes: 1 << uint(settings.ConfProxy.HttpsConfig.MaxHeaderBytes),
	}

	fmt.Printf("[INFO] https_proxy_run:%d\n", settings.ConfProxy.HttpsConfig.HTTPSPort)
	if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"),
		cert_file.Path("server.key")); err != nil && err != http.ErrServerClosed {
		mylog.Log.Fatal("listen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
			"err": err,
		})
	}
}

func HttpsServerStop() {
	shutdownTrace := mylog.NewTrace()
	mylog.Log.Info("Shutdown", shutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": fmt.Sprintf("Shutdown httpsServer%s", HttpsSrvHandler.Addr),
	})

	// 创建一个 5 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		mylog.Log.Fatal("Shutdown", shutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
	}

	mylog.Log.Info("Server exiting", shutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": fmt.Sprintf("httpsServer%s exiting", HttpsSrvHandler.Addr),
	})
}

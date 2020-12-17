package router

import (
	"context"
	"fmt"
	mylog "github.com/captainlee1024/go-gateway/internal/gateway/log"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	// 5. 注册路由
	r := SetUp() // 6. 启动服务（开启平滑下线）
	HttpSrvHandler = &http.Server{
		Addr:           fmt.Sprintf(":%v", settings.ConfBase.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(settings.ConfBase.HttpConfig.ReadTime) * time.Second,
		WriteTimeout:   time.Duration(settings.ConfBase.HttpConfig.WriteTime) * time.Second,
		MaxHeaderBytes: 1 << uint(settings.ConfBase.HttpConfig.MaxHeaderBytes),
	}
	go func() {
		fmt.Printf("[INFO] HTTPServerRun:%d\n", settings.ConfBase.Port)
		if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mylog.Log.Fatal("listen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
				"err": err,
			})
		}
	}()
}

func HttpServerStop() {
	shoutdownTrace := mylog.NewTrace()
	mylog.Log.Info("Shoutdown", shoutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": "Shoutdown Server ...",
	})

	// 创建一个 5 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		mylog.Log.Fatal("Shoutdown", shoutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
	}

	mylog.Log.Info("Server exiting", shoutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": "Server exiting",
	})
}

package main

import (
	"flag"
	"fmt"
	mylog "github.com/captainlee1024/go-gateway/internal/gateway/log"
	"github.com/captainlee1024/go-gateway/internal/gateway/router"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/pkg/snowflake"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 启动参数
// 1. endpoint: dashboard(后台管理) server(代理服务器)
// 2. config: ./configs/dev/(对应配置文件及)
var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like: ./configs/dev/")
)

// @title Go-Gateway
// @version 1.0
// @description Go-Gateway 是基于 Go 语言实现的网关！
// @termsOfService http://swagger.io/terms/

// @contact.name CaptainLee1024
// @contact.url http://blog.leecoding.club
// @contact.email 644052732@qq.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	// 首先解析传入的命令行参数
	flag.Parse()
	if *endpoint == "" {
		// 解析不到参数，打印提示信息并退出
		flag.Usage()
		os.Exit(1)
	}
	if *endpoint == "" {
		// 解析不到参数，打印提示信息并退出
		flag.Usage()
		os.Exit(1)
	}

	// 根据参数启动后台管理服务或者代理服务
	if *endpoint == "dashboard" {
		// 1. 加载配置
		// 2. 初始化日志
		if err := settings.Init(*config); err != nil {
			// log.Fatal(err)
			panic(err)
		}

		// 释放 mysql 资源，并且刷新缓冲里的日志信息
		defer settings.Destroy()
		// 注册路由，开启服务
		router.HttpServerRun()

		// 等待中断信号来优雅关闭服务器，为关闭服务器操作设置一个5秒的延时
		quit := make(chan os.Signal, 1)
		// kill 默认会发送 syscall.SIGTERM 信号
		// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
		// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
		// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
		<-quit

		// 收到信号，开始平滑下线
		router.HttpServerStop()
	} else {
		// 1. 加载配置
		// 2. 初始化日志
		if err := settings.Init(*config); err != nil {
			// log.Fatal(err)
			panic(err)
		}

		trace := mylog.NewTrace()
		// 3. 初始化 MySQL 连接
		if err := settings.InitDBPool(); err != nil {
			mylog.Log.Error("mysql", trace, mylog.DLTagUndefind, map[string]interface{}{
				"error": err,
			})
		}
		// 释放 mysql 资源，并且刷新缓冲里的日志信息
		defer func() {
			log.Println("------------------------------------------------------------------------")
			log.Printf("[INFO] %s\n", " start destroy resources.")
			settings.Close()
			mylog.Log.L.Sync()
			log.Printf("[INFO] %s\n", " success destroy resources.")
		}()

		// 4. 初始化 Redis 连接
		defaultConn, err := settings.ConnFactory("default")
		if err != nil {
			mylog.Log.Error("redis", trace, mylog.DLTagUndefind, map[string]interface{}{
				"error": err,
			})
		}
		defer defaultConn.Close()

		// 初始化雪花算法
		if err := snowflake.Init(settings.ConfBase.StartTime, settings.ConfBase.MachineID); err != nil {
			mylog.Log.Error("initSnowflake", trace, mylog.DLTagUndefind, map[string]interface{}{
				"error": err,
			})
			return
		}

		// 注册路由，开启服务
		// 这里替换成启动代理服务的方法
		//router.HttpServerRun()
		fmt.Println("start proxyServer...")

		// 等待中断信号来优雅关闭服务器，为关闭服务器操作设置一个5秒的延时
		quit := make(chan os.Signal, 1)
		// kill 默认会发送 syscall.SIGTERM 信号
		// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
		// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
		// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
		<-quit

		// 收到信号，开始平滑下线
		//router.HttpServerStop()
	}

}

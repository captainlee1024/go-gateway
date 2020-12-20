package router

import (
	_ "github.com/captainlee1024/go-gateway/docs"
	"github.com/captainlee1024/go-gateway/internal/gateway/controller"
	mylog "github.com/captainlee1024/go-gateway/internal/gateway/log"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// SetUp 初始化路由
func SetUp() *gin.Engine {
	// 当系统设置为 relase 的时候，为发布模式，其他配置都将设置成 debug 模式
	if settings.ConfBase.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	adminLoginRouter := r.Group("/admin_login")
	// 参数分别是
	// 最大空闲连接数
	// 网络类型：tcp/udp
	// 地址：host:port
	// 连接密码
	// 秘钥
	store, err := sessions.NewRedisStore(10, "tcp", "127.0.0.1:6379",
		"644315", []byte("secret"))
	if err != nil {
		mylog.Log.Fatal("sessions.NewRedisStore", mylog.NewTrace(),
			mylog.DLTagUndefind, map[string]interface{}{"error": err})
	}
	adminLoginRouter.Use(
		middleware.GinRecovery(true),
		sessions.Sessions("mysession", store),
		middleware.RequestLog(),
		middleware.IPAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}

	adminRouter := r.Group("/admin")
	adminRouter.Use(
		middleware.GinRecovery(true),
		sessions.Sessions("mysession", store),
		middleware.RequestLog(),
		middleware.IPAuthMiddleware(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminRegister(adminRouter)
	}

	serviceRouter := r.Group("/service")
	serviceRouter.Use(
		middleware.GinRecovery(true),
		sessions.Sessions("mysession", store),
		middleware.RequestLog(),
		middleware.IPAuthMiddleware(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		controller.ServiceRegister(serviceRouter)
	}

	return r
}

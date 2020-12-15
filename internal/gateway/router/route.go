package router

import (
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/gin-gonic/gin"
)

// SetUp 初始化路由
func SetUp() *gin.Engine {
	// 当系统设置为 relase 的时候，为发布模式，其他配置都将设置成 debug 模式
	if settings.ConfBase.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	v1 := r.Group("/gateway/v1")
	v1.Use(
		middleware.RequestLog(),
		middleware.GinRecovery(true),
		middleware.TranslationMiddleware(),
		middleware.IPAuthMiddleware(),
	)
	{
		v1.GET("/", func(c *gin.Context) {
			middleware.ResponseSuccess(c, "welcome to gateway!")
		})

	}

	return r
}

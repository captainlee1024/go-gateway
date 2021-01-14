package http_proxy_router

import (
	_ "github.com/captainlee1024/go-gateway/docs"
	"github.com/captainlee1024/go-gateway/internal/gateway/controller"
	middleware2 "github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/http_middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
)

// setUp 初始化路由
func setUp() *gin.Engine {
	// 当系统设置为 relase 的时候，为发布模式，其他配置都将设置成 debug 模式
	if settings.ConfProxy.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.GinRecovery(true),
		middleware.RequestLog())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	oauth := r.Group("/oauth")
	oauth.Use(middleware2.TranslationMiddleware())
	{
		controller.OauthRegister(oauth)
	}
	r.Use(http_middleware.HTTPAccessModeMiddleware(),

		http_middleware.HTTPFlowCountMiddleware(),
		http_middleware.HTTPFlowLimitMiddleware(),

		http_middleware.HTTPJWTOauthTokenMiddleware(),
		http_middleware.HTTPJWTFlowCountMiddleware(),
		http_middleware.HTTPJWTFlowLimitMiddleware(),
		http_middleware.HTTPWhiteListMiddleware(),
		http_middleware.HTTPBlackListMiddleware(),

		http_middleware.HTTPHeaderTransferMiddleware(),
		http_middleware.HTTPStripURIMiddleware(),
		http_middleware.HTTPURLRewriteMiddleware(),

		http_middleware.HTTPReverseProxyMiddleware())

	return r
}

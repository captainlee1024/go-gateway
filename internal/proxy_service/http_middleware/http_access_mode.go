package http_middleware

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/gin-gonic/gin"
)

// HTTPAccessModeMiddleware 匹配请求接入方式
func HTTPAccessModeMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		service, err := po.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}

		// 匹配成功
		fmt.Println("matched service: ", public.ObjectToJson(service))
		c.Set("service", service)
		c.Next()
	}
}

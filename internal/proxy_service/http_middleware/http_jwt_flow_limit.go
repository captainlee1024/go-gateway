package http_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowlimit"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/gin-gonic/gin"
)

func HTTPJWTFlowLimitMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*po.App)

		// 客户端基于客户IP限流，每个IP都有一个限流器进行监视
		// 单机拦截并发数
		if appInfo.Qps > 0 {
			clientLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}

			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5004, errors.New(fmt.Sprintf(
					"%v flow limit %v", c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

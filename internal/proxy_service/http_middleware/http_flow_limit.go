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

func HTTPFlowLimitMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 5001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		// 服务限流，为０时不限制
		qps := float64(serviceDetail.AccessControl.ServiceFlowLimit)
		if qps > 0 {
			serviceLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName, qps)
			if err != nil {
				middleware.ResponseError(c, 5002, err)
				c.Abort()
				return
			}

			if !serviceLimiter.Allow() {
				middleware.ResponseError(c, 5003, errors.New(fmt.Sprintf(
					"service flow limit %v", qps)))
				c.Abort()
				return
			}
		}

		// 客户端基于客户IP限流
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				middleware.ResponseError(c, 5004, err)
				c.Abort()
				return
			}

			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5005, errors.New(fmt.Sprintf(
					"%v flow limit %v", c.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

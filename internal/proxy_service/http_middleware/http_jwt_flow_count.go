package http_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowcount"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/gin-gonic/gin"
)

func HTTPJWTFlowCountMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*po.App)

		// 1. 租户流量统计

		appCounter, err := flowcount.FlowCounterHandler.GetCounter(
			public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit: %v, current: %v",
				appInfo.Qpd, appCounter.TotalCount)))
			c.Abort()
			return
		}
		//dayServiceCount, _ := appCounter.GetDayData(time.Now())
		fmt.Printf("appCounter qps:%v, dayCount:%v", appCounter.QPS, appCounter.TotalCount)
		c.Next()
	}
}

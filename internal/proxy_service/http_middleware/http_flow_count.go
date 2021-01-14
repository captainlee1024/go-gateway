package http_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowcount"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/gin-gonic/gin"
	"time"
)

func HTTPFlowCountMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 4001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		// 1. 全站
		// 2. 服务统计
		totalCounter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		totalCounter.Increase()
		dayCount, _ := totalCounter.GetDayData(time.Now())
		fmt.Printf("totalCounter qps:%v, dayCOunt:%v", totalCounter.QPS, dayCount)

		serviceCounter, err := flowcount.FlowCounterHandler.GetCounter(
			public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()
		dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		fmt.Printf("serviceCOunter qps:%v, dayCount:%v", serviceCounter.QPS, dayServiceCount)

		c.Next()
	}
}

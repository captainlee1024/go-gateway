package tcp_proxy_middleware

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowcount"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"time"
)

func TCPFlowCountMiddleware() func(c *TCPSliceRouterContext) {
	return func(c *TCPSliceRouterContext) {
		serviceInterface := c.Get("service")
		if serviceInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		// 1. 全站
		// 2. 服务统计
		totalCounter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()
		dayCount, _ := totalCounter.GetDayData(time.Now())
		fmt.Printf("totalCounter qps:%v, dayCOunt:%v", totalCounter.QPS, dayCount)

		serviceCounter, err := flowcount.FlowCounterHandler.GetCounter(
			public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		fmt.Printf("serviceCOunter qps:%v, dayCount:%v", serviceCounter.QPS, dayServiceCount)

		c.Next()
	}
}

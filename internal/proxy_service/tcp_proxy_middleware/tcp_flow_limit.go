package tcp_proxy_middleware

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowlimit"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"strings"
)

func TCPFlowLimitMiddleware() func(c *TCPSliceRouterContext) {
	return func(c *TCPSliceRouterContext) {
		serviceInterface := c.Get("service")
		if serviceInterface == nil {
			c.conn.Write([]byte("get service empty"))
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
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}

			if !serviceLimiter.Allow() {
				//middleware.ResponseError(c, 5003, errors.New(fmt.Sprintf(
				//	"service flow limit %v", qps)))
				c.conn.Write([]byte(fmt.Sprintf("service flow limit %v", qps)))
				c.Abort()
				return
			}
		}

		// 客户端基于客户IP限流

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}

			if !clientLimiter.Allow() {
				//middleware.ResponseError(c, 5005, errors.New(fmt.Sprintf(
				//	"%v flow limit %v", c.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.conn.Write([]byte(fmt.Sprintf(
					"%v flow limit %v", clientIP, serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

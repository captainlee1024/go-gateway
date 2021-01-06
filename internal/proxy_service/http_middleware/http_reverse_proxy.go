package http_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/reverse_proxy"
	"github.com/gin-gonic/gin"
)

func HTTPReverseProxyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		lb, err := po.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
		}

		transport, err := po.TransportHandler.GetTransport(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
		}
		// 创建 reverseProxy
		fmt.Println("\n===>proxy")
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, transport)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

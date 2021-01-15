package tcp_proxy_middleware

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"strings"
)

// HTTPBlackListMiddleware ip黑名单
func TCPBlackListMiddleware() func(c *TCPSliceRouterContext) {
	return func(c *TCPSliceRouterContext) {
		serviceInterface := c.Get("service")
		if serviceInterface == nil {
			//middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		whiteIPList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIPList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		blackIPList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIPList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(whiteIPList) == 0 && len(blackIPList) > 0 {
			if public.InStringSlice(blackIPList, clientIP) {
				//middleware.ResponseError(c, 2001, errors.New(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				c.conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

package tcp_proxy_middleware

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"strings"
)

// HTTPWhiteListMiddleware ip白名单列表
func TCPWhiteListMiddleware() func(c *TCPSliceRouterContext) {
	return func(c *TCPSliceRouterContext) {
		serviceInterface := c.Get("service")
		if serviceInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		// 不能直接使用Split去拿，当为空的时候，它会返回一个空的字符串
		// 先初始化为空的字符切片
		ipList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			splits := strings.Split(c.conn.RemoteAddr().String(), ":")
			clientIP := ""
			if len(splits) == 2 {
				clientIP = splits[0]
			}
			if !public.InStringSlice(ipList, clientIP) {
				//middleware.ResponseError(c, 2001, errors.New(fmt.Sprintf(
				//	"%s not in white ip", c.ClientIP())))
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

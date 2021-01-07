package http_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/gin-gonic/gin"
	"strings"
)

// HTTPWhiteListMiddleware ip白名单列表
func HTTPWhiteListMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
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
			if !public.InStringSlice(ipList, c.ClientIP()) {
				middleware.ResponseError(c, 2001, errors.New(fmt.Sprintf(
					"%s not in white ip", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

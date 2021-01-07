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

// HTTPBlackListMiddleware ip黑名单
func HTTPBlackListMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
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

		if serviceDetail.AccessControl.OpenAuth == 1 && len(whiteIPList) == 0 && len(blackIPList) > 0 {
			if public.InStringSlice(blackIPList, c.ClientIP()) {
				middleware.ResponseError(c, 2001, errors.New(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

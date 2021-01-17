package http_middleware

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/gin-gonic/gin"
	"strings"
)

func HTTPJWTOauthTokenMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInfo, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInfo.(*po.ServiceDetail)

		// 获取claims
		// 从claims获取appID
		// 根据appID从appList中获取appInfo
		// AppInfo 放入上下文中
		auth := c.GetHeader("Authorization")
		token := strings.ReplaceAll(auth, "Bearer ", "")
		appMatched := false
		if token != "" {
			claims, err := public.JWTDeCode(token)
			if err != nil {
				middleware.ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			for _, appItem := range po.AppManagerHandler.AppSlice {
				if claims.Issuer == appItem.AppID {
					c.Set("app", appItem)
					appMatched = true
					break
				}
			}
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			middleware.ResponseError(c, 2003, errors.New("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}

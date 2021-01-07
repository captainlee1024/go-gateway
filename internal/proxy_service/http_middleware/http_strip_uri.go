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

// HTTPStripURIMiddleware stripURI
// http://127.0.0.1:8880/test/a
// http://127.0.0.1:2004/a
func HTTPStripURIMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*po.ServiceDetail)

		if serviceDetail.HTTPRule.NeedStripUri == 1 && serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			// 使用空格替换掉请求地址
			fmt.Println("before: ", c.Request.URL.Path)
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
			fmt.Println("after: ", c.Request.URL.Path)
		}
		c.Next()
	}
}

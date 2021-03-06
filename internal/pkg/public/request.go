package public

import (
	"errors"
	"strconv"

	mylog "github.com/captainlee1024/go-gateway/internal/gateway/log"
	proxylog "github.com/captainlee1024/go-gateway/internal/proxy_service/log"
	"github.com/gin-gonic/gin"
)

//const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// ProxyGetGinTraceContext 从gin的Context中获取数据
func ProxyGetGinTraceContext(c *gin.Context) *proxylog.TraceContext {
	// 防御
	if c == nil {
		return proxylog.NewTrace()
	}
	traceContext, exists := c.Get(ContextTrace)
	if exists {
		if tc, ok := traceContext.(*proxylog.TraceContext); ok {
			return tc
		}
	}
	return proxylog.NewTrace()
}

// GetGinTraceContext 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *mylog.TraceContext {
	// 防御
	if c == nil {
		return mylog.NewTrace()
	}
	traceContext, exists := c.Get(ContextTrace)
	if exists {
		if tc, ok := traceContext.(*mylog.TraceContext); ok {
			return tc
		}
	}
	return mylog.NewTrace()
}

// GetCurrentUserID 获取当前登录用户的 ID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// GetPageInfo 设置分页
func GetPageInfo(c *gin.Context) (int, int) {
	// 获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		size int64
		page int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil { // 页数默认值为 1
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil { // size 默认值为 10
		size = 10
	}
	//fmt.Printf("\n\n%d<%d>\n\n", page, size)

	return int(page), int(size)
}

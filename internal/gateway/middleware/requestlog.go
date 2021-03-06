package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/captainlee1024/go-gateway/internal/pkg/public"

	mylog "github.com/captainlee1024/go-gateway/internal/gateway/log"
	"github.com/gin-gonic/gin"
)

// RequestLog 请求日志中间件
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequestInLog(c)
		defer RequestOutLog(c)
		c.Next()
	}
}

// RequestInLog 请求进来的日志
func RequestInLog(c *gin.Context) {
	// 设置 traceID spanID cspanID 及开始时间
	traceContext := mylog.NewTrace()
	if traceID := c.Request.Header.Get(public.HeaderTraceID); traceID != "" {
		traceContext.TraceID = traceID
	}
	if spanID := c.Request.Header.Get(public.HeaderSpanID); spanID != "" {
		traceContext.SpandID = spanID
	}

	c.Set(public.ContextStartTime, time.Now())
	c.Set(public.ContextTrace, traceContext)

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	// write body back
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	mylog.Log.Info(c.Request.Method+"  "+c.Request.URL.Path, traceContext, mylog.DLTagRequestIn, map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"query":  c.Request.URL.RawQuery,
		"body":   string(bodyBytes),
		"from":   c.ClientIP(),
	})
}

// RequestOutLog 请求返回是的日志
func RequestOutLog(c *gin.Context) {
	endExecTime := time.Now()
	response, _ := c.Get(public.CtxResponseKey)
	st, _ := c.Get(public.ContextStartTime)
	startExecTime := st.(time.Time)
	traceContext := public.GetGinTraceContext(c)

	mylog.Log.Info(c.Request.Method+"  "+c.Request.URL.Path, traceContext, mylog.DLTagRequestOut, map[string]interface{}{
		"status":   c.Writer.Status(),
		"method":   c.Request.Method,
		"uri":      c.Request.RequestURI,
		"response": response,
		//"user-agent": c.Request.UserAgent(),
		"proc_time": endExecTime.Sub(startExecTime),
	})
}

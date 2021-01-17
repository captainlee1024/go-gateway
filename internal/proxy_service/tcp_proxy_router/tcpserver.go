package tcp_proxy_router

import (
	"context"
	"fmt"
	mylog "github.com/captainlee1024/go-gateway/internal/proxy_service/log"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/reverse_proxy"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/tcp_proxy_middleware"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/tcp_server"
	"net"
)

var tcpServerList []*tcp_server.TCPServer

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TCPServerRun() {
	serviceList := po.ServiceManagerHandler.GetTCPServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem

		fmt.Printf("[INFO] tcp_proxy_run:%d\n", tempItem.TCPRule.Port)
		//mylog.Log.Info("listen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
		//	"msg": "tcpServer" + fmt.Sprintf("%d", tempItem.TCPRule.Port) + " exiting",
		//})

		go func(serviceDetail *po.ServiceDetail) {

			//fmt.Printf("\n\n%#v%#v%#v%#v\n\n", serviceDetail.Info, serviceDetail.TCPRule,
			//	serviceDetail.LoadBalance, serviceDetail.AccessControl)
			// 基于 thrift 代理测试 tcp 中间件
			rb, err := po.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				mylog.Log.Fatal("tcpListen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
					"err": fmt.Sprintf("GetTCPLoadBalancer :%d err: %v", serviceDetail.TCPRule.Port, err),
				})
			}

			// 构建路由并设置中间件
			//counter, _ := public.NewFlowCountService("local_app", time.Second)
			router := tcp_proxy_middleware.NewTCPSliceRouter()
			router.Group("/").Use(
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				tcp_proxy_middleware.TCPFlowLimitMiddleware(),
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)

			// 构建回调 handler
			routerHandler := tcp_proxy_middleware.NewTCPSliceRouterHandler(
				func(c *tcp_proxy_middleware.TCPSliceRouterContext) tcp_server.TCPHandler {
					return reverse_proxy.NewTCPLoadBalanceReverseProxy(c, rb)
					//return proxy.NewTCPLoadBalanceReverseProxy(c, rb)
				}, router)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := tcp_server.TCPServer{
				Addr:    fmt.Sprintf(":%d", serviceDetail.TCPRule.Port),
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, &tcpServer)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				mylog.Log.Fatal("tcpListen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
					"err": err,
				})
			}
		}(tempItem)
	}
}

func TCPServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		mylog.Log.Info("Shutdown", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
			"msg": "tcpServer" + fmt.Sprintf("%s", tcpServer.Addr) + " exiting",
		})
	}
}

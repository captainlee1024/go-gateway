package grpc_proxy_router

import (
	"context"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/grpc_proxy_middleware"
	mylog "github.com/captainlee1024/go-gateway/internal/proxy_service/log"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/reverse_proxy"
	"github.com/captainlee1024/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"net"
)

var grpcServerList []*warpGRPCServer

type warpGRPCServer struct {
	Addr string
	*grpc.Server
}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func GRPCServerRun() {
	serviceList := po.ServiceManagerHandler.GetGRPCServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem

		go func(serviceDetail *po.ServiceDetail) {

			//fmt.Printf("\n\n%#v%#v%#v%#v\n\n", serviceDetail.Info, serviceDetail.GRPCRule,
			//	serviceDetail.LoadBalance, serviceDetail.AccessControl)
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)
			fmt.Printf("[INFO] grpc_proxy_run%s\n", addr)
			rb, err := po.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				mylog.Log.Fatal("grpcListen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
					"err": fmt.Sprintf("GetTCPLoadBalancer %s err: %v", addr, err),
				})
			}

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				mylog.Log.Fatal("grpcListen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
					"err": err,
				})
			}
			grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpc_proxy_middleware.GRPCFlowCountMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCFlowLimitMiddleware(serviceDetail),

					grpc_proxy_middleware.GRPCJWTOauthTokenMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCJWTFlowCountMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCJWTFlowLimitMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCWhiteListMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCBlackListMiddleware(serviceDetail),

					grpc_proxy_middleware.GRPCHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),         // 自定义 codec
				grpc.UnknownServiceHandler(grpcHandler)) // 自定义全局回调
			grpcServerList = append(grpcServerList, &warpGRPCServer{
				Addr:   addr,
				Server: s,
			})
			if err := s.Serve(lis); err != nil {
				mylog.Log.Fatal("grpcListen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
					"err": fmt.Sprintf("grpc run %s err: %v", addr, err),
				})
			}
		}(tempItem)
	}
}

func GRPCServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		mylog.Log.Info("Shutdown", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
			"msg": "grpcServer" + fmt.Sprintf("%s", grpcServer.Addr) + " exiting",
		})
	}
}

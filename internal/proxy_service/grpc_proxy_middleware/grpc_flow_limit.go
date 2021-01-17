package grpc_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowlimit"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

func GRPCFlowLimitMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// 服务限流，为０时不限制
		qps := float64(serviceDetail.AccessControl.ServiceFlowLimit)
		if qps > 0 {
			serviceLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName, qps)
			if err != nil {
				return err
			}

			if !serviceLimiter.Allow() {
				return errors.New(fmt.Sprintf("service flow limit %v", qps))
			}
		}

		// 客户端基于客户IP限流
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String() // ip:port
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		fmt.Println("ClientIP", clientIP)
		//splits := strings.Split(peerCtx.Addr.String(), ":")
		//clientIP := ""
		//if len(splits) == 2 {
		//	clientIP = splits[0]
		//}
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				return err
			}

			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf(
					"%v flow limit %v", clientIP, serviceDetail.AccessControl.ClientIPFlowLimit))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

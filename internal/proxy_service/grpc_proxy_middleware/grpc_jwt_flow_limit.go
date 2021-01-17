package grpc_proxy_middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowlimit"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

func GRPCJWTFlowLimitMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		appInterfaces := md.Get("app")
		if len(appInterfaces) == 0 {
			if err := handler(srv, ss); err != nil {
				log.Printf("RPC failed with error %v\n", err)
				return err
			}
			return nil
		}
		appInfo := &po.App{}
		if err := json.Unmarshal([]byte(appInterfaces[0]), appInfo); err != nil {
			return err
		}

		// 客户端基于客户IP限流，每个IP都有一个限流器进行监视
		// 单机拦截并发数
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String() // ip:port
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		//peerCtx.Addr.String() // ip:port
		//splits := strings.Split(peerCtx.Addr.String(), ":")
		//clientIP := ""
		//if len(splits) == 2 {
		//	clientIP = splits[0]
		//}

		if appInfo.Qps > 0 {
			clientLimiter, err := flowlimit.FlowLimiterHandler.GetLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+clientIP,
				float64(appInfo.Qps))
			if err != nil {
				return err
			}

			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf("%v flow limit %v", clientIP, appInfo.Qps))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

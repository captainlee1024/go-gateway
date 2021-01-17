package grpc_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

// GRPCWhiteListMiddleware ip白名单列表
func GRPCWhiteListMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// 不能直接使用Split去拿，当为空的时候，它会返回一个空的字符串
		// 先初始化为空的字符切片
		ipList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
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

		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, clientIP) {
				return errors.New(fmt.Sprintf(
					"%s not in white ip", clientIP))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

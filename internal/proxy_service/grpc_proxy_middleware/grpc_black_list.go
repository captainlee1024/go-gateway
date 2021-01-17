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

// GRPCBlackListMiddleware ip黑名单
func GRPCBlackListMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		whiteIPList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIPList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		blackIPList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIPList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
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
		//fmt.Println("\n\n\nblackMiddleware============>", whiteIPList, blackIPList)
		if serviceDetail.AccessControl.OpenAuth == 1 && len(whiteIPList) == 0 && len(blackIPList) > 0 {
			if public.InStringSlice(blackIPList, clientIP) {
				return errors.New(fmt.Sprintf("%s in black ip list", clientIP))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

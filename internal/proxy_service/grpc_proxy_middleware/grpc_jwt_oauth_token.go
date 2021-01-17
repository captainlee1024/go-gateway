package grpc_proxy_middleware

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

func GRPCJWTOauthTokenMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// 获取claims
		// 从claims获取appID
		// 根据appID从appList中获取appInfo
		// AppInfo 放入上下文中
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		authToken := ""
		auths := md.Get("Authorization")
		if len(auths) > 0 {
			authToken = auths[0]
		}
		token := strings.ReplaceAll(authToken, "Bearer ", "")
		appMatched := false
		if token != "" {
			claims, err := public.JWTDeCode(token)
			if err != nil {
				return err
			}
			for _, appItem := range po.AppManagerHandler.AppSlice {
				if claims.Issuer == appItem.AppID {
					md.Set("app", public.ObjectToJson(appItem))
					appMatched = true
					break
				}
			}
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			return errors.New("not match valid app")
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

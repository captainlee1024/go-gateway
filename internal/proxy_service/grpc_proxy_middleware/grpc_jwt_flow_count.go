package grpc_proxy_middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowcount"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func GRPCJWTFlowCountMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
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
		// 1. 租户流量统计

		appCounter, err := flowcount.FlowCounterHandler.GetCounter(
			public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			return err
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			return errors.New(fmt.Sprintf("租户日请求量限流 limit: %v, current: %v",
				appInfo.Qpd, appCounter.TotalCount))
		}
		//dayServiceCount, _ := appCounter.GetDayData(time.Now())
		fmt.Printf("appCounter qps:%v, dayCount:%v", appCounter.QPS, appCounter.TotalCount)
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

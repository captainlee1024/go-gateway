package grpc_proxy_middleware

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/flowcount"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"log"
	"time"
)

func GRPCFlowCountMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// 1. 全站
		// 2. 服务统计
		totalCounter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		dayCount, _ := totalCounter.GetDayData(time.Now())
		fmt.Printf("totalCounter qps:%v, dayCOunt:%v", totalCounter.QPS, dayCount)

		serviceCounter, err := flowcount.FlowCounterHandler.GetCounter(
			public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()
		dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		fmt.Printf("serviceCOunter qps:%v, dayCount:%v", serviceCounter.QPS, dayServiceCount)

		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

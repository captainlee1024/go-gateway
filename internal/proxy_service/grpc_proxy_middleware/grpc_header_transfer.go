package grpc_proxy_middleware

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

// GRPCHeaderTransferMiddleware header头转换
func GRPCHeaderTransferMiddleware(serviceDetail *po.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		for _, item := range strings.Split(serviceDetail.GRPCRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				md.Set(items[1], items[2])
			}
			if items[0] == "del" {
				delete(md, items[1])
			}
		}
		if err := ss.SetHeader(md); err != nil {
			return err
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}

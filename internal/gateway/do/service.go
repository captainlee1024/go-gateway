package do

import "github.com/captainlee1024/go-gateway/internal/gateway/po"

// ServiceListInput 查询信息 Do
type ServiceListInput struct {
	Info     string // 关键词
	PageNo   int    // 页数
	PageSize int    // 每页条数
}

// ServiceDetail 查询服务列表 Do
type ServiceDetail struct {
	Info          *po.ServiceInfo          `json:"info" description:"基本信息"`
	HTTPRule      *po.ServiceHTTPRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *po.ServiceTCPRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *po.ServiceGRPCRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *po.ServiceLoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *po.ServiceAccessControl `json:"access_control" description:"access_control"`
}

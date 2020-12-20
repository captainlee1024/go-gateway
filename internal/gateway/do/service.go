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
	Info          *po.ServiceInfo
	HTTPRule      *po.ServiceHTTPRule
	TCPRule       *po.ServiceTCPRule
	GRPCRule      *po.ServiceGRPCRule
	LoadBalance   *po.ServiceLoadBalance
	AccessControl *po.ServiceAccessControl
}

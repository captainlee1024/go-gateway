package dto

import (
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

type ServiceUpdateTcpInput struct {
	ID          int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`                                        // 服务ID
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`                    // 服务描述

	Port int `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"` // 端口，需要设置8001-8999范围内

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`                              // 是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`   // 白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`           // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`               // 服务端限流

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                            // 轮询策略
	IpList     string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`         // IP列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"` // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`            // 禁用IP列表
}

func (params *ServiceUpdateTcpInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceAddTcpInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`                    // 服务描述

	Port int `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"` // 端口，需要设置8001-8999范围内

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`                              // 是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`   // 白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`           // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`               // 服务端限流

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                            // 轮询策略
	IpList     string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`         // IP列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"` // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`            // 禁用IP列表
}

func (params *ServiceAddTcpInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceUpdateGrpcInput struct {
	ID          int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`                                        // 服务ID
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`                    // 服务描述

	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`         // 端口，需要设置8001-8999范围内
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"` // metadata转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`                              // 是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`   // 白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`           // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`               // 服务端限流

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                            // 轮询策略
	IpList     string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`         // IP列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"` // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`            // 禁用IP列表
}

func (params *ServiceUpdateGrpcInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceAddGrpcInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`                    // 服务描述

	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`         // 端口，需要设置8001-8999范围内
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"` // metadata转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`                              // 是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"` // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`   // 白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`           // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`               // 服务端限流

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                            // 轮询策略
	IpList     string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`         // IP列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"` // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`            // 禁用IP列表
}

func (params *ServiceAddGrpcInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" example:"" validate:""`         // 列表
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" example:"" validate:""` // 列表
}

type ServiceStatInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"1" validate:"required"` // 服务ID
}

func (params *ServiceStatInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceDetailInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"1" validate:"required"` // 服务ID
}

func (params *ServiceDetailInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceUpdateHTTPInput struct {
	ID          int64  `json:"id" form:"id" comment:"服务 ID" example:"66" validate:"required,min=1"`                               // 服务 ID
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`      // 服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                           // 接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或前缀" example:"" validate:"required,valid_rule"`                       // 接入路径：域名或前缀
	NeedHTTPs      int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                      // 支持https
	NeedStripURI   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`          // 启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"" validate:"max=1,min=0"`        // 是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"RUL重写功能" example:"" validate:"valid_url_rewrite"`              // RUL重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"Header转换" example:"" validate:"valid_header_transfor"` // Header转换

	OpenAuth         int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`            // 是否开启权限
	BlackList        string `json:"black_list" form:"black_list" comment:"黑名单IP" example:"" validate:""`                      // 黑名单IP
	WhiteList        string `json:"white_list" form:"white_list" comment:"白名单IP" example:"" validate:""`                      // 白名单IP
	ClientFlowLimit  int    `json:"client_flow_limit" form:"client_flow_limit" comment:"客户端IP限流" example:"" validate:"min=0"` // 客户端IP限流
	ServiceFlowLimit int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"" validate:"min=0"` // 服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`                           // 轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"IP列表" example:"" validate:"required,valid_ipportlist"`                   // IP列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"IP权重" example:"" validate:"required,valid_weightlist"`           // IP权重
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时" example:"" validate:"min=0"`   // 建立连接超时
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时" example:"" validate:"min=0"` // 获取header超时
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲连接时间" example:"" validate:"min=0"`     // 链接最大空闲连接时间
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲连接数" example:"" validate:"min=0"`                // 最大空闲连接数
}

func (params *ServiceUpdateHTTPInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceAddHTTPInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`      // 服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                           // 接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或前缀" example:"" validate:"required,valid_rule"`                       // 接入路径：域名或前缀
	NeedHTTPs      int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                      // 支持https
	NeedStripURI   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`          // 启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"" validate:"max=1,min=0"`        // 是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"RUL重写功能" example:"" validate:"valid_url_rewrite"`              // RUL重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"Header转换" example:"" validate:"valid_header_transfor"` // Header转换

	OpenAuth         int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`            // 是否开启权限
	BlackList        string `json:"black_list" form:"black_list" comment:"黑名单IP" example:"" validate:""`                      // 黑名单IP
	WhiteList        string `json:"white_list" form:"white_list" comment:"白名单IP" example:"" validate:""`                      // 白名单IP
	ClientFlowLimit  int    `json:"client_flow_limit" form:"client_flow_limit" comment:"客户端IP限流" example:"" validate:"min=0"` // 客户端IP限流
	ServiceFlowLimit int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"" validate:"min=0"` // 服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`                           // 轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"IP列表" example:"" validate:"required,valid_ipportlist"`                   // IP列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"IP权重" example:"" validate:"required,valid_weightlist"`           // IP权重
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时" example:"" validate:"min=0"`   // 建立连接超时
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时" example:"" validate:"min=0"` // 获取header超时
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲连接时间" example:"" validate:"min=0"`     // 链接最大空闲连接时间
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲连接数" example:"" validate:"min=0"`                // 最大空闲连接数
}

func (params *ServiceAddHTTPInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceDeleteInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"1" validate:"required"` // 服务ID
}

func (params *ServiceDeleteInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      // 关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        // 页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` // 每页条数
}

func (params *ServiceListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     // id
	ServiceName string `json:"service_name" form:"service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` // 服务描述
	ServiceAddr string `json:"service_addr" form:"service_addr"` // 服务地址
	LoadType    int    `json:"load_type" form:"load_type"`       // 服务类型
	Qps         int    `json:"qps" form:"qps"`                   // qps
	Qpd         int    `json:"qpd" form:"qpd"`                   // qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     // 节点数
}

type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数" example:"" validate:""` // 总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   // 列表
}

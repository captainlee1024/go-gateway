package dto

import (
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

type AppListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                                    // 关键词
	PageSize int    `json:"page_size" form:"page_size" comment:"每页大小" example:"10" validate:"required,min=1,max=999"` // 每页个数
	PageNo   int    `json:"page_no" form:"page_no" comment:"页码" example:"1" validate:"required,min=1,max=999"`        // 页码
}

func (params *AppListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AppListOutput AppList 输出结构体
type AppListOutput struct {
	List  []AppListItemOutput `json:"list" form:"list" comment:"租户列表"`   // 租户列表
	Total int64               `json:"total" form:"total" comment:"租户总数"` // 租户总数
}

// AppListItemOutput .
type AppListItemOutput struct {
	ID       int64  `json:"id" form:"" comment:"" example:"" validate:""`                             //
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" example:"" validate:""`                  // 每秒请求量限制
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" example:"" validate:""`                   // 日请求量限制
	RealQps  int64  `json:"real_qps" form:"real_qps" comment:"当前实际 QPS" example:"" validate:""`       // 当前实际 QPS
	RealQpd  int64  `json:"real_qpd" form:"real_qpd" comment:"当前实际 QPD" example:"" validate:""`       // 当前实际 QPD
	AppID    string `json:"app_id" form:"app_id" comment:"租户ID" example:"" validate:""`               // 租户ID
	Name     string `json:"name" form:"name" comment:"租户名称" example:"" validate:""`                   // 租户名称
	Secret   string `json:"secret" form:"secret" comment:"秘钥" example:"" validate:""`                 // 秘钥
	WhiteIPs string `json:"white_ips" form:"white_ips" comment:"IP白名单，支持前缀匹配" example:"" validate:""` // IP白名单，支持前缀匹配
}

type AppDetailInput struct {
	ID int64 `json:"id" form:"id" comment:"租户ID" example:"" validate:"required"` // 租户 ID
}

func (params *AppDetailInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type AppDetailOutput struct {
	ID       int64  `json:"id" form:"" comment:"" example:"" validate:""`                             //
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" example:"" validate:""`                  // 每秒请求量限制
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" example:"" validate:""`                   // 日请求量限制
	RealQps  int64  `json:"real_qps" form:"real_qps" comment:"当前实际 QPS" example:"" validate:""`       // 当前实际 QPS
	RealQpd  int64  `json:"real_qpd" form:"real_qpd" comment:"当前实际 QPD" example:"" validate:""`       // 当前实际 QPD
	AppID    string `json:"app_id" form:"app_id" comment:"租户ID" example:"" validate:""`               // 租户ID
	Name     string `json:"name" form:"name" comment:"租户名称" example:"" validate:""`                   // 租户名称
	Secret   string `json:"secret" form:"secret" comment:"秘钥" example:"" validate:""`                 // 秘钥
	WhiteIPs string `json:"white_ips" form:"white_ips" comment:"IP白名单，支持前缀匹配" example:"" validate:""` // IP白名单，支持前缀匹配
}

type AppDeleteInput struct {
	ID int64 `json:"id" form:"id" comment:"租户ID" example:"" validate:"required"` // 租户 ID
}

func (params *AppDeleteInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type AppAddInput struct {
	AppID    string `json:"app_id" form:"app_id" comment:"租户ID" example:"" validate:"required"`       // 租户ID
	Name     string `json:"name" form:"name" comment:"租户名称" example:"" validate:"required"`           // 租户名称
	Secret   string `json:"secret" form:"secret" comment:"秘钥" example:"" validate:""`                 // 秘钥
	WhiteIPs string `json:"white_ips" form:"white_ips" comment:"IP白名单，支持前缀匹配" example:"" validate:""` // IP白名单，支持前缀匹配
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" example:"" validate:""`                  // 每秒请求量限制
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" example:"" validate:""`                   // 日请求量限制
}

func (params *AppAddInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type AppUpdateInput struct {
	ID       int64  `json:"id" form:"id" comment:"租户ID" example:"" validate:"required"`               // 主键ID
	AppID    string `json:"app_id" form:"app_id" comment:"租户ID" example:"" validate:""`               // 租户ID
	Name     string `json:"name" form:"name" comment:"租户名称" example:"" validate:"required"`           // 租户名称
	Secret   string `json:"secret" form:"secret" comment:"秘钥" example:"" validate:"required"`         // 秘钥
	WhiteIPs string `json:"white_ips" form:"white_ips" comment:"IP白名单，支持前缀匹配" example:"" validate:""` // IP白名单，支持前缀匹配
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" example:"" validate:""`                  // 每秒请求量限制
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" example:"" validate:""`                   // 日请求量限制
}

func (params *AppUpdateInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type AppStatInput struct {
	ID int64 `json:"id" form:"id" comment:"租户ID" example:"" validate:"required"` // 租户 ID
}

func (params *AppStatInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type AppStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日统计" validate:"required"`         // 今日流量
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日统计" validate:"required"` // 昨日流量
}

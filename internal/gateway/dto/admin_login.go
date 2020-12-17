package dto

import (
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

/*
	json: 结构体序列化成 json
	form: json 反序列化到结构体
	comment: 对应字段的错误输出会使用 comment
	example: swagger 文档的默认值
	validate: 校验规则
	// : 注释会在 swagger 文档里显示出来
*/

// AdminLoginInput 登录接口输入参数
type AdminLoginInput struct {
	Username string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,valid_username"` // 管理员账户
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required,valid_password"` // 密码
}

func (params *AdminLoginInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AdminLoginOutput 登录接口输出参数
type AdminLoginOutput struct {
	Token string `json:"token" from:"token" comment:"token" example:"" validate:""` // 用户 token
}

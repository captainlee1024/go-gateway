package dto

import (
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminInfoOutput struct {
	ID           int       `json:"id"`           // 管理员 ID
	Username     string    `json:"user_name"`    // 登录管理员姓名
	LoginTime    time.Time `json:"login_time"`   // 登录时间
	Avatar       string    `json:"avatar"`       // 登录用户头像
	Introduction string    `json:"introduction"` // 介绍
	Roles        []string  `json:"roles"`        // .
}

type ChangePwdInput struct {
	Password string `json:"password" from:"password" comment:"新密码" example:"123456" validate:"required,valid_password"` // 新密码
}

func (params *ChangePwdInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

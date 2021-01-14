package dto

import (
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

type TokensInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"` // 授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围" example:"read_write" validate:"required"`                   // 权限范围
}

func (params *TokensInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type TokensOutput struct {
	AccessToken string `json:"access_token"` // Token
	TokenType   string `json:"token_type"`   // Token 类型
	Scope       string `json:"scope"`        // 权限范围
	ExpiresIn   int    `json:"expires_in"`   // 过期时间
}

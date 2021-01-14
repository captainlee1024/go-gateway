package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

type OauthServer interface {
	Tokens(ctx *gin.Context)
}

type UnimplementedOauthServer struct {
}

func (u *UnimplementedOauthServer) Tokens(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method tokens not implemented"))
}

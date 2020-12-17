package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

// AdminServer is the server API for Admin service
// All implementations embed UnimplementedAdminServer
// for forward compatibility
type AdminServer interface {
	AdminInfo(*gin.Context)
	ChangePwd(*gin.Context)
}

type UnimplementedAdminServer struct {
}

func (u UnimplementedAdminServer) AdminInfo(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method AdminInfo not implemented"))
}

func (u *UnimplementedAdminServer) ChangePwd(c *gin.Context) {
	middleware.ResponseError(c, 2002, errors.New("method ChangePwd not implemented"))
}

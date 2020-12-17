package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

// AdminLoginServer is the server IPA for AdminLogin service
// All implementations must embed UnimplementedAdminLoginServer
// for forward compatibility
type AdminLoginServer interface {
	// AdminLogin is the Login service
	AdminLogin(*gin.Context)
	// AdminLogOut is the LogOut service
	AdminLogOut(*gin.Context)
}

// UnimplementedAdminLoginServer must be embedded to have forward compatible implementations.
type UnimplementedAdminLoginServer struct {
}

func (u *UnimplementedAdminLoginServer) AdminLogin(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("meth AdminLogin not implemented"))
}

func (u *UnimplementedAdminLoginServer) AdminLogOut(c *gin.Context) {
	middleware.ResponseError(c, 2002, errors.New("meth AdminLogOut not implemented"))
}

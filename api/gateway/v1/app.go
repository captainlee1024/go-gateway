package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

// AppServer is the server IPA for AdminLogin service
// All implementations must embed UnimplementedAppServer
// for forward compatibility
type AppServer interface {
	AppList(*gin.Context)
	AppDetail(*gin.Context)
	AppStat(*gin.Context)
	AppAdd(*gin.Context)
	AppUpdate(*gin.Context)
	AppDelete(*gin.Context)
}

// UnimplementedAppServer must be embedded to have forward compatible implementations.
type UnimplementedAppServer struct{}

func (u *UnimplementedAppServer) AppList(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method AppList not implemented"))
}
func (u *UnimplementedAppServer) AppDetail(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method AppList not implemented"))
}
func (u *UnimplementedAppServer) AppStat(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method AppList not implemented"))
}
func (u *UnimplementedAppServer) AppAdd(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method AppList not implemented"))
}
func (u *UnimplementedAppServer) AppUpdate(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method AppList not implemented"))
}
func (u *UnimplementedAppServer) AppDelete(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method AppList not implemented"))
}

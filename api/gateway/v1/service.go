package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

type ServiceServer interface {
	ServiceList(*gin.Context)
	ServiceDelete(*gin.Context)
	ServiceAddHTTP(*gin.Context)
	ServiceUpdateHTTP(*gin.Context)
	ServiceDetail(*gin.Context)
	ServiceStat(*gin.Context)
}

type UnimplementedServiceServer struct{}

func (u *UnimplementedServiceServer) ServiceList(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("meth ServiceList not implemented"))
}

func (u *UnimplementedServiceServer) ServiceDelete(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method ServiceDelete not implemented"))
}

func (u *UnimplementedServiceServer) ServiceAddHTTP(c *gin.Context) {
	middleware.ResponseError(c, 2002, errors.New("method ServiceAddHTTP not implemented"))
}

func (u *UnimplementedServiceServer) ServiceUpdateHTTP(c *gin.Context) {
	middleware.ResponseError(c, 2003, errors.New("method ServiceUpdateHTTP not implemented"))
}

func (u *UnimplementedServiceServer) ServiceDetail(c *gin.Context) {
	middleware.ResponseError(c, 2004, errors.New("method ServiceDetail not implemented"))
}

func (u *UnimplementedServiceServer) ServiceStat(c *gin.Context) {
	middleware.ResponseError(c, 2004, errors.New("method ServiceStat not implemented"))
}

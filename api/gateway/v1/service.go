package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

type ServiceServer interface {
	ServiceList(*gin.Context)
}

type UnimplementedServiceServer struct{}

func (u *UnimplementedServiceServer) ServiceList(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("meth ServiceList not implemented"))
}

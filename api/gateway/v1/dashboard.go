package v1

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

// DashboardServer is the server IPA for AdminLogin service
// All implementations must embed UnimplementedDashboardServer
// for forward compatibility
type DashboardServer interface {
	PanelGroupData(*gin.Context)
	FlowStat(*gin.Context)
	DashboardServiceStat(*gin.Context)
}

// UnimplementedDashboardServer must be embedded to have forward compatible implementations.
type UnimplementedDashboardServer struct{}

func (u *UnimplementedDashboardServer) PanelGroupData(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("method PanelGroupData not implemented"))
}

func (u *UnimplementedDashboardServer) FlowStat(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method FlowStat not implemented"))
}

func (u *UnimplementedDashboardServer) DashboardServiceStat(c *gin.Context) {
	middleware.ResponseError(c, 2002, errors.New("method ServiceStat not implemented"))
}

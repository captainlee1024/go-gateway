package data

import (
	"database/sql"
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

var _ service.DashboardRepo = (service.DashboardRepo)(nil)

type dashboardRepo struct{}

func NewDashboardRepo() service.DashboardRepo {
	return &dashboardRepo{}
}

// GetServiceNum 获取总服务数
func (repo *dashboardRepo) GetServiceNum(c *gin.Context) (total int64, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return 0, err
	}
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT COUNT(*) FROM (
			SELECT * FROM gateway_service_info WHERE is_DELETE = 0) a`
	if err = settings.SqlxLogGet(trace, db, &total, sqlStr); err != nil {
		return 0, err
	}
	return total, nil
}

// GetAppNum 获取总租户数
func (repo *dashboardRepo) GetAppNum(c *gin.Context) (total int64, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return 0, err
	}
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT COUNT(*) FROM (
			SELECT * FROM gateway_app WHERE is_DELETE = 0) a`
	if err = settings.SqlxLogGet(trace, db, &total, sqlStr); err != nil {
		return 0, err
	}
	return total, nil
}

// GetServiceStat　获取类型及对应总数
func (repo *dashboardRepo) GetServiceStat(c *gin.Context) (serviceStatList []*do.DashboardServiceStat, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}
	trace := public.GetGinTraceContext(c)
	serviceStatList = make([]*do.DashboardServiceStat, 0, 2)

	sqlStr := `SELECT COUNT(*) as value, load_type as name FROM(
		SELECT * FROM gateway_service_info WHERE is_delete = 0
		) a
		GROUP BY a.load_type`
	if err = settings.SqlxLogSelect(trace, db, &serviceStatList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("没有数据，先去添加数据吧")
		}
		return nil, err
	}
	return serviceStatList, err
}

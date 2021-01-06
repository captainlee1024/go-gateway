package po

import (
	"database/sql"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
)

type ServiceAccessControl struct {
	ID                int64  `db:"id" json:"id"`
	ServiceID         int64  `db:"service_id" json:"service_id"`
	BlackList         string `db:"black_list" json:"black_list"`
	WhiteList         string `db:"white_list" json:"white_list"`
	WhiteHostName     string `db:"white_host_name" json:"white_host_name"`
	OpenAuth          int    `db:"open_auth" json:"open_auth"`
	ClientIPFlowLimit int    `db:"clientip_flow_limit" json:"clientip_flow_limit"`
	ServiceFlowLimit  int    `db:"service_flow_limit" json:"service_flow_limit"`
}

func GetServiceAccessControllerByID(ID int64, c *gin.Context,
) (accessController *ServiceAccessControl, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	accessController = new(ServiceAccessControl)
	trace := public.ProxyGetGinTraceContext(c)
	accessControlSqlStr := `SELECT * FROM gateway_service_access_control WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, accessController, accessControlSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return accessController, nil
}

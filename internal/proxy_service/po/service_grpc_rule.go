package po

import (
	"database/sql"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
)

type ServiceGRPCRule struct {
	ID             int64  `db:"id" json:"id"`
	ServiceID      int64  `db:"service_id" json:"service_id"`
	Port           int    `db:"port" json:"port"`
	HeaderTransfor string `db:"header_transfor" json:"header_transfor"`
}

func GetServiceGRPCRuleByID(ID int64, c *gin.Context,
) (grpcRule *ServiceGRPCRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	grpcRule = new(ServiceGRPCRule)
	trace := public.ProxyGetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_grpc_rule WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, grpcRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return grpcRule, nil
}

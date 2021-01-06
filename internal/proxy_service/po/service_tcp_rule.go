package po

import (
	"database/sql"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
)

type ServiceTCPRule struct {
	ID        int64 `db:"id" json:"id"`
	ServiceID int64 `db:"service_id" json:"service_id"`
	Port      int   `db:"port" json:"port"`
}

func GetServiceTCPRuleByID(ID int64, c *gin.Context,
) (tcpRule *ServiceTCPRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	tcpRule = new(ServiceTCPRule)
	trace := public.ProxyGetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_tcp_rule WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, tcpRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return tcpRule, nil
}

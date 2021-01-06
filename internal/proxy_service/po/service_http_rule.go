package po

import (
	"database/sql"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
)

type ServiceHTTPRule struct {
	ID             int64  `db:"id" json:"id"`
	ServiceID      int64  `db:"service_id" json:"service_id"`
	RuleType       int    `db:"rule_type" json:"rule_type"`
	NeedHTTPs      int    `db:"need_https" json:"need_https"`
	NeedStripUri   int    `db:"need_strip_uri" json:"need_strip_uri"`
	NeedWebsocket  int    `db:"need_websocket" json:"need_websocket"`
	Rule           string `db:"rule" json:"rule"`
	UrlRewrite     string `db:"url_rewrite" json:"url_rewrite"`
	HeaderTransfor string `db:"header_transfor" json:"header_transfor"`
}

func GetServiceHTTPRuleByID(ID int64, c *gin.Context,
) (httpRule *ServiceHTTPRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	httpRule = new(ServiceHTTPRule)
	trace := public.ProxyGetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_http_rule WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, httpRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return httpRule, nil
}

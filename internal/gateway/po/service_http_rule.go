package po

type ServiceHTTPRule struct {
	ID             int64  `db:"id"`
	ServiceID      int64  `db:"service_id"`
	RuleType       int    `db:"rule_type"`
	NeedHTTPs      int    `db:"need_https"`
	NeedStripUri   int    `db:"need_strip_uri"`
	NeedWebsocket  int    `db:"need_websocket"`
	Rule           string `db:"rule"`
	UrlRewrite     string `db:"url_rewrite"`
	HeaderTransfor string `db:"header_transfor"`
}

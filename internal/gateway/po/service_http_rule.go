package po

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

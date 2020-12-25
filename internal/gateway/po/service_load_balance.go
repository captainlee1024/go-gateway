package po

type ServiceLoadBalance struct {
	ID                     int64  `db:"id" json:"id"`
	ServiceID              int64  `db:"service_id" json:"service_id"`
	CheckMethod            int    `db:"check_method" json:"check_method"`
	CheckTimeout           int    `db:"check_timeout" json:"check_timeout"`
	CheckInterval          int    `db:"check_interval" json:"check_interval"`
	RoundType              int    `db:"round_type" json:"round_type"`
	IPList                 string `db:"ip_list" json:"ip_list"`
	WeightList             string `db:"weight_list" json:"weight_list"`
	ForbidList             string `db:"forbid_list" json:"forbid_list"`
	UpStreamConnectTimeout int    `db:"upstream_connect_timeout" json:"upstream_connect_timeout"`
	UpStreamHeaderTimeout  int    `db:"upstream_header_timeout" json:"upstream_header_timeout"`
	UpStreamIdleTimeout    int    `db:"upstream_idle_timeout" json:"upstream_idle_timeout"`
	UpStreamMaxIdle        int    `db:"upstream_max_idle" json:"upstream_max_idle"`
}

package po

type ServiceLoadBalance struct {
	ID                     int64  `db:"id"`
	ServiceID              int64  `db:"service_id"`
	CheckMethod            int    `db:"check_method"`
	CheckTimeout           int    `db:"check_timeout"`
	CheckInterval          int    `db:"check_interval"`
	RoundType              int    `db:"round_type"`
	IPList                 string `db:"ip_list"`
	WeightList             string `db:"weight_list"`
	ForbidList             string `db:"forbid_list"`
	UpStreamConnectTimeout int    `db:"upstream_connect_timeout"`
	UpStreamHeaderTimeout  int    `db:"upstream_header_timeout"`
	UpStreamIdleTimeout    int    `db:"upstream_idle_timeout"`
	UpStreamMaxIdle        int    `db:"upstream_max_idle"`
}

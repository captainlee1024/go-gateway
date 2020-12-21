package po

type ServiceAccessControl struct {
	ID                int64  `db:"id"`
	ServiceID         int64  `db:"service_id"`
	BlackList         string `db:"black_list"`
	WhiteList         string `db:"white_list"`
	WhiteHostName     string `db:"white_host_name"`
	OpenAuth          int    `db:"open_auth"`
	ClientIPFlowLimit int    `db:"clientip_flow_limit"`
	ServiceFlowLimit  int    `db:"service_flow_limit"`
}

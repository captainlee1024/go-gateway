package po

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

package po

type ServiceGRPCRule struct {
	ID             int64  `db:"id" json:"id"`
	ServiceID      int64  `db:"service_id" json:"service_id"`
	Port           int    `db:"port" json:"port"`
	HeaderTransfor string `db:"header_transfor" json:"header_transfor"`
}

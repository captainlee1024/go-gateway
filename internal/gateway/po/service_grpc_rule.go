package po

type ServiceGRPCRule struct {
	ID             int64  `db:"id"`
	ServiceID      int64  `db:"service_id"`
	Port           int    `db:"port"`
	HeaderTransfor string `db:"header_transfor"`
}

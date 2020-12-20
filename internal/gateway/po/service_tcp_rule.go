package po

type ServiceTCPRule struct {
	ID        int64 `db:"id"`
	ServiceID int64 `db:"service_id"`
	Port      int   `db:"port"`
}

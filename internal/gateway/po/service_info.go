package po

import "time"

type ServiceInfo struct {
	ID          int64     `db:"id"`
	LoadType    int       `db:"load_type"`
	IsDelete    int       `db:"is_delete"`
	ServiceName string    `db:"service_name"`
	ServiceDesc string    `db:"service_desc"`
	CreatedAt   time.Time `db:"create_at"'`
	UpdatedAt   time.Time `db:"update_at"`
}

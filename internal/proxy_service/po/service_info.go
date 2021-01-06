package po

import "time"

type ServiceInfo struct {
	ID          int64     `db:"id" json:"id"`
	LoadType    int       `db:"load_type" json:"load_type"`
	IsDelete    int       `db:"is_delete" json:"is_delete"`
	ServiceName string    `db:"service_name" json:"service_name"`
	ServiceDesc string    `db:"service_desc" json:"service_desc"`
	CreatedAt   time.Time `db:"create_at" json:"create_at"'`
	UpdatedAt   time.Time `db:"update_at" json:"update_at"`
}

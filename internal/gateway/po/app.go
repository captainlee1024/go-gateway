package po

import "time"

type App struct {
	ID        int64     `db:"id"`
	Qps       int64     `db:"qps"`
	Qpd       int64     `db:"qpd"`
	AppID     string    `db:"app_id"`
	Name      string    `db:"name"`
	Secret    string    `db:"secret"`
	WhiteIPs  string    `db:"white_ips"`
	CreatedAt time.Time `db:"create_at"`
	UpdatedAt time.Time `db:"update_at"`
	IsDelete  int       `db:"is_delete"`
}

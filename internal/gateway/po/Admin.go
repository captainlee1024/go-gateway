package po

import "time"

// Admin 管理员持久化实体
type Admin struct {
	ID         int       `db:"id"`
	IsDelete   int       `db:"is_delete"`
	Username   string    `db:"user_name"`
	Salt       string    `db:"salt"`
	Password   string    `db:"password"`
	CreateTime time.Time `db:"create_at"`
	UpdateTime time.Time `db:"update_at"`
}

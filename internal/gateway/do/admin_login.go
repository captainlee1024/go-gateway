package do

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// AdminLogin 业务逻辑实体
type AdminLogin struct {
	ID       int
	Username string
	Password string
	Salt     string
}

// PasswordCheck 密码校验
func (adminLogin *AdminLogin) PasswordCheck(salt, password string) bool {
	s1 := sha256.New()
	s1.Write([]byte(adminLogin.Password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	str2 := fmt.Sprintf("%x", s2.Sum(nil))
	return str2 == password
}

// AdminSessionInfo 管理员登录的 Session
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	Username  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}

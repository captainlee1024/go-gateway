package data

import (
	"database/sql"
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

// 在编译的时候可以知道这个对象是否实现了这个 interface{}
var _ service.AdminLoginRepo = (service.AdminLoginRepo)(nil)

func NewAdminLoginRepo() service.AdminLoginRepo {
	return new(adminLoginRepo)
}

type adminLoginRepo struct{}

func (admin *adminLoginRepo) GetAdmin(loginDo *do.AdminLogin, c *gin.Context) (adminDo *do.AdminLogin, err error) {
	// do -> po
	adminPo := &po.Admin{
		Username: loginDo.Username,
	}

	// 获取 db
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	// 查询
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT id, user_name, password, salt FROM gateway_admin WHERE user_name = ?`
	err = settings.SqlxLogGet(trace, db, adminPo, sqlStr, adminPo.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在！")
		}
		return nil, err
	}

	// po -> do
	adminDo = &do.AdminLogin{
		ID:       adminPo.ID,
		Username: adminPo.Username,
		Password: adminPo.Password,
		Salt:     adminPo.Salt,
	}

	// 返回结果
	return adminDo, nil
}

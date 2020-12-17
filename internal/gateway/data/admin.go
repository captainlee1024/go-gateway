package data

import (
	"database/sql"
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/data/mysql"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"time"
)

var _ service.AdminRepo = (service.AdminRepo)(nil)

type adminRepo struct{}

func NewAdminRepo() service.AdminRepo {
	return &adminRepo{}
}

func (a *adminRepo) GetAdmin(ID int, c *gin.Context) (adminInfoDo *do.AdminLogin, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}
	adminInfo := &po.Admin{
		ID: ID,
	}

	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT user_name, password, salt
			FROM gateway_admin
			WHERE id = ?`
	err = mysql.SqlxLogGet(trace, db, adminInfo, sqlStr, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	adminInfoDo = &do.AdminLogin{
		ID:       ID,
		Username: adminInfo.Username,
		Salt:     adminInfo.Salt,
	}
	return adminInfoDo, nil
}

func (a *adminRepo) UpdatePassword(changeDo *do.AdminLogin, c *gin.Context) (err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return err
	}

	currentTime := time.Now()
	trace := public.GetGinTraceContext(c)
	sqlStr := `UPDATE gateway_admin
			SET password=?, create_at=?
			WHERE id=?`
	_, err = mysql.SqlxLogExec(trace, db, sqlStr, changeDo.Password, currentTime, changeDo.ID)
	if err != nil {
		return err
	}
	return nil
}

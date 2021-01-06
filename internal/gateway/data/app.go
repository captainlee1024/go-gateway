package data

import (
	"database/sql"
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

var _ service.AppRepo = (service.AppRepo)(nil)

type appRepo struct{}

func NewAppRepo() service.AppRepo {
	return &appRepo{}
}

func (repo *appRepo) GetAppList(info string, page, size int, c *gin.Context) (appList []*po.App, total int64, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, 0, err
	}

	appList = make([]*po.App, 0, 2)
	trace := public.GetGinTraceContext(c)
	sqlSrt := `SELECT *
			FROM gateway_app
			WHERE (name LIKE ? OR app_id LIKE ?) AND is_delete = 0
			ORDER BY id DESC
			LIMIT ?,?`
	if err = settings.SqlxLogSelect(trace, db, &appList, sqlSrt, "%"+info+"%", "%"+info+"%", (page-1)*size, size); err != nil {
		if err == sql.ErrNoRows {
			appList = nil
			err = nil
		} else {
			return nil, 0, err
		}
	}

	countSqlStr := `SELECT COUNT(*) FROM(
				SELECT * FROM gateway_app
				WHERE (name LIKE ? OR app_id LIKE ?) AND is_delete = 0) a`
	if err = settings.SqlxLogGet(trace, db, &total, countSqlStr, "%"+info+"%", "%"+info+"%"); err != nil {
		if err == sql.ErrNoRows {
			return appList, 0, nil
		}
		return nil, 0, err
	}
	return appList, total, nil
}

func (repo *appRepo) GetAppDetailByID(ID int64, c *gin.Context) (appDetail *po.App, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	appDetail = new(po.App)
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT * FROM gateway_app WHERE is_delete = 0 AND id = ?`
	if err = settings.SqlxLogGet(trace, db, appDetail, sqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在！")
		}
		return nil, err
	}
	return appDetail, nil
}

func (repo *appRepo) DeleteAppByID(ID int64, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}
	trace := public.GetGinTraceContext(c)
	sqlStr := `UPDATE gateway_app SET is_delete = 1 WHERE id = ?`
	_, err = settings.SqlxLogExec(trace, db, sqlStr, ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *appRepo) InsertApp(app *po.App, c *gin.Context) (ID int64, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return 0, err
	}
	trace := public.GetGinTraceContext(c)
	sqlStr := `INSERT INTO gateway_app(
			app_id, name, secret, white_ips, qpd, qps, create_at, update_at, is_delete)
			VALUES(?,?,?,?,?,?,?,?,?)`
	ret, err := settings.SqlxLogExec(trace, db, sqlStr,
		app.AppID, app.Name, app.Secret, app.WhiteIPs, app.Qpd, app.Qps,
		app.CreatedAt, app.UpdatedAt, app.IsDelete)
	if err != nil {
		return 0, err
	}

	ID, err = ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (repo *appRepo) GetAppDetailByAppID(appID string, c *gin.Context) (appDetail *po.App, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	appDetail = new(po.App)
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT * FROM gateway_app WHERE is_delete = 0 AND app_id = ?`
	if err = settings.SqlxLogGet(trace, db, appDetail, sqlStr, appID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在！")
		}
		return nil, err
	}
	return appDetail, nil
}

func (repo *appRepo) UpdateApp(app *po.App, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	trace := public.GetGinTraceContext(c)
	sqlStr := `UPDATE gateway_app SET name=?, secret=?, white_ips=?, qpd=?, qps=?, update_at=?
			WHERE id = ?`
	_, err = settings.SqlxLogExec(trace, db, sqlStr,
		app.Name, app.Secret, app.WhiteIPs, app.Qpd, app.Qps, app.UpdatedAt, app.ID)
	return err
}

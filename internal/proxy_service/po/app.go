package po

import (
	"database/sql"
	"errors"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sync"
	"time"
)

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

var AppManagerHandler *AppManager

func init() {
	AppManagerHandler = NewAppManagerHandler()
}

type AppManager struct {
	AppMap   map[string]*App
	AppSlice []*App
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

func NewAppManagerHandler() *AppManager {
	return &AppManager{
		AppMap:   map[string]*App{},
		AppSlice: []*App{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

// LoadOnce 服务加载到内存
func (a *AppManager) LoadOnce() error {
	// 只执行一次，加载配置到 map 和 slice
	a.init.Do(func() {

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		list, _, err := getAppList(1, 9999, c)
		if err != nil {
			a.err = err
			return
		}

		a.Locker.Lock()
		defer a.Locker.Unlock()

		for _, listItem := range list {
			tmpItem := listItem
			appDetail, err := GetAppDetailByID(tmpItem.ID, c)
			if err != nil {
				a.err = err
				return
			}
			a.AppMap[listItem.Name] = appDetail
			a.AppSlice = append(a.AppSlice, appDetail)
		}
	})

	//fmt.Printf("=====>%#v", a.ServiceMap)
	//fmt.Printf("=====>%#v", a.ServiceSlice)
	return a.err
}

func (a *AppManager) GetAppList() []*App {
	return a.AppSlice
}

func getAppList(page, size int, c *gin.Context) (appList []*App, total int64, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, 0, err
	}

	appList = make([]*App, 0, 2)
	trace := public.ProxyGetGinTraceContext(c)
	sqlSrt := `SELECT *
			FROM gateway_app
			WHERE is_delete = 0
			ORDER BY id DESC
			LIMIT ?,?`
	if err = settings.SqlxLogSelect(trace, db, &appList, sqlSrt, (page-1)*size, size); err != nil {
		if err == sql.ErrNoRows {
			appList = nil
			err = nil
		} else {
			return nil, 0, err
		}
	}

	countSqlStr := `SELECT COUNT(*) FROM(
				SELECT * FROM gateway_app
				WHERE is_delete = 0) a`
	if err = settings.SqlxLogGet(trace, db, &total, countSqlStr); err != nil {
		if err == sql.ErrNoRows {
			return appList, 0, nil
		}
		return nil, 0, err
	}
	return appList, total, nil
}

func GetAppDetailByID(ID int64, c *gin.Context) (appDetail *App, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	appDetail = new(App)
	trace := public.ProxyGetGinTraceContext(c)
	sqlStr := `SELECT * FROM gateway_app WHERE is_delete = 0 AND id = ?`
	if err = settings.SqlxLogGet(trace, db, appDetail, sqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在！")
		}
		return nil, err
	}
	return appDetail, nil
}

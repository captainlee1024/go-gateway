package po

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"sync"
)

// ServiceListInput 查询信息 Do
type ServiceListInput struct {
	Info     string // 关键词
	PageNo   int    // 页数
	PageSize int    // 每页条数
}

// ServiceDetail 查询服务列表 Do
type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info" description:"基本信息"`
	HTTPRule      *ServiceHTTPRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *ServiceTCPRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *ServiceGRPCRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *ServiceAccessControl `json:"access_control" description:"access_control"`
}

var ServiceManagerHandler *ServiceManager

// 在加载的时候初始化，加载配置到内存
func init() {
	ServiceManagerHandler = NewServiceManager()
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex // 锁
	init         sync.Once    // 初始化
	err          error
}

// NewServiceManager xx
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

// HTTPAccessMode 匹配接入
func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {
	// 1. 前缀匹配 /abc -> serviceSlice.rule
	// 2. 域名匹配 www.text.com -> serviceSlice.rule

	// host c.Request.Host www.text.com:8080
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	fmt.Println("host", host)

	// path c.Request.URL.Path /abc/get?xxx=xxx
	path := c.Request.URL.Path
	fmt.Println("path", path)

	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}

		// 域名匹配
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}

		// 前缀匹配
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}

// LoadOnce 服务加载到内存
func (s *ServiceManager) LoadOnce() error {
	// 只执行一次，加载配置到 map 和 slice
	s.init.Do(func() {

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		list, _, err := GetServiceInfoList(1, 9999, c)
		if err != nil {
			s.err = err
			return
		}

		s.Locker.Lock()
		defer s.Locker.Unlock()

		for _, listItem := range list {
			tmpItem := listItem
			serviceDetail, err := GetServiceDetail(tmpItem, c)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})

	//fmt.Printf("=====>%#v", s.ServiceMap)
	//fmt.Printf("=====>%#v", s.ServiceSlice)
	return s.err
}

func GetServiceInfoList(page, size int, c *gin.Context,
) (serviceInfoList []*ServiceInfo, total int64, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, 0, err
	}

	trace := public.ProxyGetGinTraceContext(c)
	//sqlStr := `SELECT id, load_type, service_name, service_desc, create_at, update_at, is_delete
	sqlStr := `SELECT *
			FROM gateway_service_info
			WHERE is_delete = 0
			ORDER BY id DESC
			LIMIT ?,?`
	serviceInfoList = make([]*ServiceInfo, 0, 2)
	err = settings.SqlxLogSelect(trace, db, &serviceInfoList, sqlStr, (page-1)*size, size)
	if err != nil {
		if err == sql.ErrNoRows {
			serviceInfoList = nil
			err = nil
		} else {
			return nil, 0, err
		}
	}

	//countSqlStr := `SELECT COUNT(*) FROM (SELECT * FROM gateway_service_info
	//		WHERE (service_name LIKE ? OR service_desc LIKE ?) AND is_delete = 0
	//		LIMIT ?,?) a`
	countSqlStr := `SELECT COUNT(*) FROM (SELECT * FROM gateway_service_info
			WHERE is_delete = 0) a`
	err = settings.SqlxLogGet(trace, db, &total, countSqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			total = 0
			return serviceInfoList, total, nil
		}
		return nil, 0, err
	}

	return serviceInfoList, total, nil
}

func GetServiceDetail(serviceInfo *ServiceInfo, c *gin.Context) (detail *ServiceDetail, err error) {
	// page size info -> 查询各个表组装数据
	// gateway_service_info -> id, service_name, service_desc, load_type
	// service_addr
	// gateway_service_tcp_rule -> service_addr
	// gateway_service_grpc_rule -> service_addr
	// gateway_service_http_rule -> service_addr
	// gateway_service_load_balance -> total_node
	// 中间件 -> qps, qpd

	detail = &ServiceDetail{Info: serviceInfo}

	if serviceInfo.LoadType == public.LoadTypeHTTP {
		// 查询 HTTPRule
		httpRule := &ServiceHTTPRule{}
		httpRule, err = GetServiceHTTPRuleByID(serviceInfo.ID, c)
		if err != nil {
			return nil, err
		}
		detail.HTTPRule = httpRule
	} else if serviceInfo.LoadType == public.LoadTypeTCP {
		// 查询 TCPRule
		tcpRule := &ServiceTCPRule{}
		tcpRule, err = GetServiceTCPRuleByID(serviceInfo.ID, c)
		if err != nil {
			return nil, err
		}
		detail.TCPRule = tcpRule
	} else if serviceInfo.LoadType == public.LoadTypeGRPC {
		// 查询 GRPCRule
		grpcRule := &ServiceGRPCRule{}
		grpcRule, err = GetServiceGRPCRuleByID(serviceInfo.ID, c)
		if err != nil {
			return nil, err
		}
		detail.GRPCRule = grpcRule
	}

	// 查询 LoadBalance
	loadBalance := &ServiceLoadBalance{}
	loadBalance, err = GetServiceLoadBalanceByID(serviceInfo.ID, c)
	if err != nil {
		return nil, err
	}
	detail.LoadBalance = loadBalance

	// 查询 AccessControl
	accessControl := &ServiceAccessControl{}
	accessControl, err = GetServiceAccessControllerByID(serviceInfo.ID, c)
	if err != nil {
		return nil, err
	}
	detail.AccessControl = accessControl

	return detail, nil
}

package data

import (
	"database/sql"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/gateway/data/mysql"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

var _ service.ServiceRepo = (service.ServiceRepo)(nil)

type serviceRepo struct{}

func NewServiceRepo() service.ServiceRepo {
	return new(serviceRepo)
}

// GetServiceInfoList 根据条件获取服务信息列表
func (repo *serviceRepo) GetServiceDetail(serviceInfo *po.ServiceInfo, c *gin.Context) (detail *do.ServiceDetail, err error) {
	// page size info -> 查询各个表组装数据
	// gateway_service_info -> id, service_name, service_desc, load_type
	// service_addr
	// gateway_service_tcp_rule -> service_addr
	// gateway_service_grpc_rule -> service_addr
	// gateway_service_http_rule -> service_addr
	// gateway_service_load_balance -> total_node
	// 中间件 -> qps, qpd

	detail = &do.ServiceDetail{Info: serviceInfo}

	if serviceInfo.LoadType == public.LoadTypeHTTP {
		// 查询 HTTPRule
		httpRule := &po.ServiceHTTPRule{}
		httpRule, err = repo.GetServiceHTTPRuleByID(serviceInfo.ID, c)
		if err != nil {
			return nil, err
		}
		detail.HTTPRule = httpRule
	} else if serviceInfo.LoadType == public.LoadTypeTCP {
		// 查询 TCPRule
		tcpRule := &po.ServiceTCPRule{}
		tcpRule, err = repo.GetServiceTCPRuleByID(serviceInfo.ID, c)
		if err != nil {
			return nil, err
		}
		detail.TCPRule = tcpRule
	} else if serviceInfo.LoadType == public.LoadTypeGRPC {
		// 查询 GRPCRule
		grpcRule := &po.ServiceGRPCRule{}
		grpcRule, err = repo.GetServiceGRPCRuleByID(serviceInfo.ID, c)
		if err != nil {
			return nil, err
		}
		detail.GRPCRule = grpcRule
	}

	// 查询 LoadBalance
	loadBalance := &po.ServiceLoadBalance{}
	loadBalance, err = repo.GetServiceLoadBalanceByID(serviceInfo.ID, c)
	if err != nil {
		return nil, err
	}
	detail.LoadBalance = loadBalance

	// 查询 AccessControl
	accessControl := &po.ServiceAccessControl{}
	accessControl, err = repo.GetServiceAccessControllerByID(serviceInfo.ID, c)
	if err != nil {
		return nil, err
	}
	detail.AccessControl = accessControl

	return detail, nil
}

// GetServiceInfoList 获取列表并返回查询的条数
func (repo *serviceRepo) GetServiceInfoList(info string, page, size int, c *gin.Context,
) (serviceInfoList []*po.ServiceInfo, total int64, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, 0, err
	}

	trace := public.GetGinTraceContext(c)
	//sqlStr := `SELECT id, load_type, service_name, service_desc, create_at, update_at, is_delete
	sqlStr := `SELECT *
			FROM gateway_service_info
			WHERE (service_name LIKE ? OR service_desc LIKE ?) AND is_delete = 0
			ORDER BY id DESC
			LIMIT ?,?`
	serviceInfoList = make([]*po.ServiceInfo, 0, 2)
	err = mysql.SqlxLogSelect(trace, db, &serviceInfoList, sqlStr, "%"+info+"%", "%"+info+"%", (page-1)*size, size)
	fmt.Printf("%#v\n", serviceInfoList)
	if err != nil {
		if err == sql.ErrNoRows {
			serviceInfoList = nil
			err = nil
		} else {
			return nil, 0, err
		}
	}

	countSqlStr := `SELECT COUNT(*) FROM (SELECT * FROM gateway_service_info
			WHERE (service_name LIKE ? OR service_desc LIKE ?) AND is_delete = 0
			LIMIT ?,?) a`
	err = mysql.SqlxLogGet(trace, db, &total, countSqlStr, "%"+info+"%", "%"+info+"%", (page-1)*size, size)
	fmt.Printf("%d\n", total)
	if err != nil {
		if err == sql.ErrNoRows {
			total = 0
			return serviceInfoList, total, nil
		}
		return nil, 0, err
	}

	return serviceInfoList, total, nil
}

// GetServiceHTTPRuleByID 根据 ID 查询一条 service_http_rule 数据
func (repo *serviceRepo) GetServiceHTTPRuleByID(ID int64, c *gin.Context,
) (httpRule *po.ServiceHTTPRule, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	httpRule = new(po.ServiceHTTPRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_http_rule WHERE service_id=?`
	if err = mysql.SqlxLogGet(trace, db, httpRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return httpRule, nil
}

// GetServiceTCPRuleByID 根据 ID 查询一条 service_tcp_rule 记录
func (repo *serviceRepo) GetServiceTCPRuleByID(ID int64, c *gin.Context,
) (tcpRule *po.ServiceTCPRule, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	tcpRule = new(po.ServiceTCPRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_tcp_rule WHERE service_id=?`
	if err = mysql.SqlxLogGet(trace, db, tcpRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return tcpRule, nil
}

// GetServiceGRPCRuleByID 根据 ID 查询一条 Service
func (repo *serviceRepo) GetServiceGRPCRuleByID(ID int64, c *gin.Context,
) (grpcRule *po.ServiceGRPCRule, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	grpcRule = new(po.ServiceGRPCRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_grpc_rule WHERE service_id=?`
	if err = mysql.SqlxLogGet(trace, db, grpcRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return grpcRule, nil
}

// GetServiceLoadBalanceByID 根据 ID 查询一条 service_load_balance　记录
func (repo *serviceRepo) GetServiceLoadBalanceByID(ID int64, c *gin.Context,
) (loadBalance *po.ServiceLoadBalance, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	loadBalance = new(po.ServiceLoadBalance)
	trace := public.GetGinTraceContext(c)
	loadBalanceSqlStr := `SELECT * FROM gateway_service_load_balance WHERE service_id=?`
	if err = mysql.SqlxLogGet(trace, db, loadBalance, loadBalanceSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return loadBalance, nil
}

// GetServiceAccessControllerByID 根据 ID 查询一条 service_access_control 记录
func (repo *serviceRepo) GetServiceAccessControllerByID(ID int64, c *gin.Context,
) (accessController *po.ServiceAccessControl, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	accessController = new(po.ServiceAccessControl)
	trace := public.GetGinTraceContext(c)
	accessControlSqlStr := `SELECT * FROM gateway_service_access_control WHERE service_id=?`
	if err = mysql.SqlxLogGet(trace, db, accessController, accessControlSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return accessController, nil
}

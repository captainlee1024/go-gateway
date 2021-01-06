package data

import (
	"database/sql"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"time"
)

var _ service.ServiceRepo = (service.ServiceRepo)(nil)

type serviceRepo struct{}

func NewServiceRepo() service.ServiceRepo {
	return new(serviceRepo)
}

// ServiceInfo 相关持久化接口实现
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
	db, err := settings.GetDBPool("default")
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
	err = settings.SqlxLogSelect(trace, db, &serviceInfoList, sqlStr, "%"+info+"%", "%"+info+"%", (page-1)*size, size)
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
			WHERE (service_name LIKE ? OR service_desc LIKE ?) AND is_delete = 0) a`
	err = settings.SqlxLogGet(trace, db, &total, countSqlStr, "%"+info+"%", "%"+info+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			total = 0
			return serviceInfoList, total, nil
		}
		return nil, 0, err
	}

	return serviceInfoList, total, nil
}

func (repo *serviceRepo) GetServiceInfoByID(ID int64, c *gin.Context) (serviceInfoPo *po.ServiceInfo, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	serviceInfoPo = &po.ServiceInfo{}
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT *
			FROM gateway_service_info
			WHERE is_delete = 0
			AND id = ?`
	if err = settings.SqlxLogGet(trace, db, serviceInfoPo, sqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrServiceNotExit
		} else {
			return nil, err
		}
	}

	return serviceInfoPo, nil
}

func (repo *serviceRepo) GetServiceInfoByName(serviceName string, c *gin.Context) (serviceInfoPo *po.ServiceInfo, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	serviceInfoPo = &po.ServiceInfo{}
	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT *
			FROM gateway_service_info
			WHERE is_delete = 0
			AND service_name = ?`
	if err = settings.SqlxLogGet(trace, db, serviceInfoPo, sqlStr, serviceName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrServiceNotExit
		} else {
			return nil, err
		}
	}
	return serviceInfoPo, nil
}

func (repo *serviceRepo) InsertServiceInfo(tx *sqlx.Tx, serviceInfo *po.ServiceInfo, c *gin.Context) (ID int64, err error) {
	sqlStr := `INSERT INTO gateway_service_info(
			load_type, service_name, service_desc, create_at, update_at, is_delete)
			values(?,?,?,?,?,?)`
	trace := public.GetGinTraceContext(c)
	ret, err := settings.SqlxLogTxExec(trace, tx, sqlStr,
		serviceInfo.LoadType,
		serviceInfo.ServiceName,
		serviceInfo.ServiceDesc,
		serviceInfo.CreatedAt,
		serviceInfo.UpdatedAt,
		serviceInfo.IsDelete)
	if err != nil {
		return 0, err
	}

	ID, err = ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	return ID, err
}

func (repo *serviceRepo) DeleteServiceInfo(serviceInfoPo *po.ServiceInfo, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	trace := public.GetGinTraceContext(c)
	sqlStr := `UPDATE gateway_service_info
			SET is_delete = 1
			WHERE id = ?`
	if _, err = settings.SqlxLogExec(trace, db, sqlStr, serviceInfoPo.ID); err != nil {
		return err
	}
	return nil
}

func (repo *serviceRepo) UpdateServiceInfo(tx *sqlx.Tx, serviceInfo *po.ServiceInfo, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_info
			SET service_desc = ?, update_at = ?
			WHERE id = ?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr, serviceInfo.ServiceDesc, time.Now(), serviceInfo.ID)
	return err
}

// HTTPRule 相关持久化接口实现
// GetServiceHTTPRuleByID 根据 ID 查询一条 service_http_rule 数据
func (repo *serviceRepo) GetServiceHTTPRuleByID(ID int64, c *gin.Context,
) (httpRule *po.ServiceHTTPRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	httpRule = new(po.ServiceHTTPRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_http_rule WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, httpRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return httpRule, nil
}

func (repo *serviceRepo) GetServiceHTTPRuleByRule(ruleType int, rule string, c *gin.Context) (httpRule *po.ServiceHTTPRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	httpRule = new(po.ServiceHTTPRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_http_rule WHERE rule_type = ? AND rule=?`
	if err = settings.SqlxLogGet(trace, db, httpRule, httpRuleSqlStr, ruleType, rule); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return httpRule, nil
}

func (repo *serviceRepo) AddHTTPDetail(httpDetail *do.ServiceDetail, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// 插入 serviceInfo，并返回 ID
	ID, err := repo.InsertServiceInfo(tx, httpDetail.Info, c)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 获取 ServiceID
	//httpDetail.Info, err = repo.GetServiceInfoByName(httpDetail.Info.ServiceName, c)
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}
	httpDetail.Info.ID = ID

	// 使用 ServiceID 插入其他表
	httpDetail.HTTPRule.ServiceID = httpDetail.Info.ID
	httpDetail.LoadBalance.ServiceID = httpDetail.Info.ID
	httpDetail.AccessControl.ServiceID = httpDetail.Info.ID

	// HTTPRule
	if err = repo.InsertServiceHTTPRule(tx, httpDetail.HTTPRule, c); err != nil {
		tx.Rollback()
		return err
	}

	// LoadBalance
	if err = repo.InsertServiceHTTPLoadBalance(tx, httpDetail.LoadBalance, c); err != nil {
		tx.Rollback()
		return err
	}

	// AccessControl
	if err = repo.InsertServiceHTTPAccessControl(tx, httpDetail.AccessControl, c); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// InsertServiceHTTPRule 事务插入一条 http_rule 数据
func (repo *serviceRepo) InsertServiceHTTPRule(tx *sqlx.Tx, httpRule *po.ServiceHTTPRule, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_http_rule(
			service_id, rule_type, rule, need_https, need_strip_uri,
			need_websocket, url_rewrite, header_transfor)
			VALUES(?,?,?,?,?,?,?,?)`
	trace := public.GetGinTraceContext(c)

	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		httpRule.ServiceID,
		httpRule.RuleType,
		httpRule.Rule,
		httpRule.NeedHTTPs,
		httpRule.NeedStripUri,
		httpRule.NeedWebsocket,
		httpRule.UrlRewrite,
		httpRule.HeaderTransfor); err != nil {
		return err
	}
	return nil
}

func (repo *serviceRepo) UpdateHTTPDetail(httpDetail *do.ServiceDetail, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// 更新 ServiceInfo
	if err = repo.UpdateServiceInfo(tx, httpDetail.Info, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 HTTPRule
	if err = repo.UpdateHTTPRule(tx, httpDetail.HTTPRule, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 LoadBalance
	if err = repo.UpdateHTTPLoadBalance(tx, httpDetail.LoadBalance, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 AccessControl
	if err = repo.UpdateHTTPAccessControl(tx, httpDetail.AccessControl, c); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (repo *serviceRepo) UpdateHTTPRule(tx *sqlx.Tx, serviceHTTPRule *po.ServiceHTTPRule, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_http_rule
			SET need_https=?, need_strip_uri=?, need_websocket=?, url_rewrite=?, header_transfor=?
			WHERE service_id=?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		serviceHTTPRule.NeedHTTPs,
		serviceHTTPRule.NeedStripUri,
		serviceHTTPRule.NeedWebsocket,
		serviceHTTPRule.UrlRewrite,
		serviceHTTPRule.HeaderTransfor,
		serviceHTTPRule.ServiceID)
	return err
}

// TCPRule 相关持久化接口实现
// GetServiceTCPRuleByID 根据 ID 查询一条 service_tcp_rule 记录
func (repo *serviceRepo) GetServiceTCPRuleByID(ID int64, c *gin.Context,
) (tcpRule *po.ServiceTCPRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	tcpRule = new(po.ServiceTCPRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_tcp_rule WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, tcpRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return tcpRule, nil
}

func (repo *serviceRepo) InsertServiceTCPRule(tx *sqlx.Tx, tcpRule *po.ServiceTCPRule, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_tcp_rule(service_id, port)
			VALUES(?,?)`

	trace := public.GetGinTraceContext(c)
	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		tcpRule.ServiceID,
		tcpRule.Port); err != nil {
		return err
	}
	return nil
}

func (repo *serviceRepo) AddTCPDetail(tcpDetail *do.ServiceDetail, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// 插入 ServiceInfo 并返回 ID
	ID, err := repo.InsertServiceInfo(tx, tcpDetail.Info, c)
	if err != nil {
		tx.Rollback()
		return err
	}

	tcpDetail.Info.ID = ID
	tcpDetail.TCPRule.ServiceID = ID
	tcpDetail.LoadBalance.ServiceID = ID
	tcpDetail.AccessControl.ServiceID = ID

	// 插入其他表
	// TCPRule
	if err = repo.InsertServiceTCPRule(tx, tcpDetail.TCPRule, c); err != nil {
		tx.Rollback()
		return err
	}

	// LoadBalance
	if err = repo.InsertServiceGRPCTCPLoadBalance(tx, tcpDetail.LoadBalance, c); err != nil {
		tx.Rollback()
		return err
	}

	// AccessControl
	if err = repo.InsertServiceGRPCTCPAccessControl(tx, tcpDetail.AccessControl, c); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (repo *serviceRepo) UpdateTCPDetail(tcpDetail *do.ServiceDetail, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// 更新 ServiceInfo
	if err = repo.UpdateServiceInfo(tx, tcpDetail.Info, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 TCPRule -> 没有可变更选项
	//if err = repo.UpdateTCPRule(tx, tcpDetail.TCPRule, c); err != nil {
	//	tx.Rollback()
	//	return err
	//}

	// 更新 LoadBalance
	if err = repo.UpdateGRPCTCPLoadBalance(tx, tcpDetail.LoadBalance, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 AccessControl
	if err = repo.UpdateGRPCTCPAccessControl(tx, tcpDetail.AccessControl, c); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (repo *serviceRepo) UpdateTCPRule(tx *sqlx.Tx, serviceTCPRule *po.ServiceTCPRule, c *gin.Context) (err error) {
	// TCPRule 里的数据都是不允许变更的
	return nil
}

func (repo *serviceRepo) GetServiceTCPRuleByPort(port int, c *gin.Context) (serviceTCPRule *po.ServiceTCPRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT * FROM gateway_service_tcp_rule
			WHERE port = ?`
	serviceTCPRule = &po.ServiceTCPRule{}
	if err = settings.SqlxLogGet(trace, db, serviceTCPRule, sqlStr, port); err != nil {
		return nil, err
	}

	return serviceTCPRule, nil
}

// GRPCRule 相关持久化接口实现
// GetServiceGRPCRuleByID 根据 ID 查询一条 Service
func (repo *serviceRepo) GetServiceGRPCRuleByID(ID int64, c *gin.Context,
) (grpcRule *po.ServiceGRPCRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	grpcRule = new(po.ServiceGRPCRule)
	trace := public.GetGinTraceContext(c)
	httpRuleSqlStr := `SELECT * FROM gateway_service_grpc_rule WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, grpcRule, httpRuleSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return grpcRule, nil
}

func (repo *serviceRepo) InsertServiceGRPCRule(tx *sqlx.Tx, grpcRule *po.ServiceGRPCRule, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_grpc_rule(
			service_id, port, header_transfor)
			VALUES(?,?,?)`
	trace := public.GetGinTraceContext(c)
	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		grpcRule.ServiceID,
		grpcRule.Port,
		grpcRule.HeaderTransfor); err != nil {
		return err
	}
	return nil
}

func (repo *serviceRepo) AddGRPCDetail(grpcDetail *do.ServiceDetail, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// 插入 ServiceInfo 并返回 ID
	ID, err := repo.InsertServiceInfo(tx, grpcDetail.Info, c)
	if err != nil {
		tx.Rollback()
		return err
	}

	grpcDetail.Info.ID = ID
	grpcDetail.GRPCRule.ServiceID = ID
	grpcDetail.LoadBalance.ServiceID = ID
	grpcDetail.AccessControl.ServiceID = ID

	// 插入其他表
	// GRPCRule
	if err = repo.InsertServiceGRPCRule(tx, grpcDetail.GRPCRule, c); err != nil {
		tx.Rollback()
		return err
	}

	// LoadBalance
	if err = repo.InsertServiceGRPCTCPLoadBalance(tx, grpcDetail.LoadBalance, c); err != nil {
		tx.Rollback()
		return err
	}

	// AccessControl
	if err = repo.InsertServiceGRPCTCPAccessControl(tx, grpcDetail.AccessControl, c); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (repo *serviceRepo) UpdateGrpcDetail(grpcDetail *do.ServiceDetail, c *gin.Context) (err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// 更新 ServiceInfo
	if err = repo.UpdateServiceInfo(tx, grpcDetail.Info, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 GRPCRule
	if err = repo.UpdateGRPCRule(tx, grpcDetail.GRPCRule, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 LoadBalance
	if err = repo.UpdateGRPCTCPLoadBalance(tx, grpcDetail.LoadBalance, c); err != nil {
		tx.Rollback()
		return err
	}

	// 更新 AccessControl
	if err = repo.UpdateGRPCTCPAccessControl(tx, grpcDetail.AccessControl, c); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (repo *serviceRepo) UpdateGRPCRule(tx *sqlx.Tx, serviceGRPCRule *po.ServiceGRPCRule, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_grpc_rule
			SET header_transfor=?
			WHERE service_id=?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr, serviceGRPCRule.HeaderTransfor, serviceGRPCRule.ServiceID)
	return err
}

func (repo *serviceRepo) GetServiceGRPCRuleByPort(port int, c *gin.Context) (serviceGRPCRUle *po.ServiceGRPCRule, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	trace := public.GetGinTraceContext(c)
	sqlStr := `SELECT * FROM gateway_service_grpc_rule
			WHERE port = ?`
	serviceGRPCRUle = &po.ServiceGRPCRule{}
	if err = settings.SqlxLogGet(trace, db, serviceGRPCRUle, sqlStr, port); err != nil {
		return nil, err
	}

	return serviceGRPCRUle, nil
}

// LoadBalance 相关持久化接口实现
// GetServiceLoadBalanceByID 根据 ID 查询一条 service_load_balance　记录
func (repo *serviceRepo) GetServiceLoadBalanceByID(ID int64, c *gin.Context,
) (loadBalance *po.ServiceLoadBalance, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	loadBalance = new(po.ServiceLoadBalance)
	trace := public.GetGinTraceContext(c)
	loadBalanceSqlStr := `SELECT * FROM gateway_service_load_balance WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, loadBalance, loadBalanceSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return loadBalance, nil
}

func (repo *serviceRepo) InsertServiceHTTPLoadBalance(tx *sqlx.Tx, loadBalance *po.ServiceLoadBalance, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_load_balance(
			service_id, round_type, ip_list, weight_list,
			upstream_connect_timeout, upstream_header_timeout,
			upstream_idle_timeout, upstream_max_idle)
			VALUES(?,?,?,?,?,?,?,?)`
	trace := public.GetGinTraceContext(c)
	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		loadBalance.ServiceID,
		loadBalance.RoundType,
		loadBalance.IPList,
		loadBalance.WeightList,
		loadBalance.UpStreamConnectTimeout,
		loadBalance.UpStreamHeaderTimeout,
		loadBalance.UpStreamIdleTimeout,
		loadBalance.UpStreamMaxIdle); err != nil {
		return err
	}
	return nil
}

func (repo *serviceRepo) UpdateHTTPLoadBalance(tx *sqlx.Tx, loadBalance *po.ServiceLoadBalance, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_load_balance
			SET round_type=?, ip_list=?, weight_list=?,
			upstream_connect_timeout=?, upstream_header_timeout=?,
			upstream_idle_timeout=?, upstream_max_idle=?
			WHERE service_id = ?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		loadBalance.RoundType,
		loadBalance.IPList,
		loadBalance.WeightList,
		loadBalance.UpStreamConnectTimeout,
		loadBalance.UpStreamHeaderTimeout,
		loadBalance.UpStreamIdleTimeout,
		loadBalance.UpStreamMaxIdle,
		loadBalance.ServiceID)

	return err
}

func (repo *serviceRepo) InsertServiceGRPCTCPLoadBalance(tx *sqlx.Tx, loadBalance *po.ServiceLoadBalance, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_load_balance(
			service_id, round_type, ip_list, weight_list, forbid_list)
			VALUES(?,?,?,?,?)`
	trace := public.GetGinTraceContext(c)
	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		loadBalance.ServiceID,
		loadBalance.RoundType,
		loadBalance.IPList,
		loadBalance.WeightList,
		loadBalance.ForbidList); err != nil {
		return err
	}

	return nil
}
func (repo *serviceRepo) UpdateGRPCTCPLoadBalance(tx *sqlx.Tx, loadBalance *po.ServiceLoadBalance, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_load_balance
			SET round_type=?, ip_list=?, weight_list=?, forbid_list=?
			WHERE service_id = ?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		loadBalance.RoundType,
		loadBalance.IPList,
		loadBalance.WeightList,
		loadBalance.ForbidList,
		loadBalance.ServiceID)

	return err
}

// AccessControl 相关持久化接口实现
// GetServiceAccessControllerByID 根据 ID 查询一条 service_access_control 记录
func (repo *serviceRepo) GetServiceAccessControllerByID(ID int64, c *gin.Context,
) (accessController *po.ServiceAccessControl, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	accessController = new(po.ServiceAccessControl)
	trace := public.GetGinTraceContext(c)
	accessControlSqlStr := `SELECT * FROM gateway_service_access_control WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, accessController, accessControlSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return accessController, nil
}

func (repo *serviceRepo) InsertServiceHTTPAccessControl(tx *sqlx.Tx, accessControl *po.ServiceAccessControl, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_access_control(
			service_id, open_auth, black_list, white_list, clientip_flow_limit, service_flow_limit)
			VALUES(?,?,?,?,?,?)`
	trace := public.GetGinTraceContext(c)
	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		accessControl.ServiceID,
		accessControl.OpenAuth,
		accessControl.BlackList,
		accessControl.WhiteList,
		accessControl.ClientIPFlowLimit,
		accessControl.ServiceFlowLimit); err != nil {
		return err
	}
	return nil
}

func (repo *serviceRepo) UpdateHTTPAccessControl(tx *sqlx.Tx, accessControl *po.ServiceAccessControl, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_access_control
			SET open_auth=?, black_list=?, white_list=?, clientip_flow_limit=?, service_flow_limit=?
			WHERE service_id = ?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		accessControl.OpenAuth,
		accessControl.BlackList,
		accessControl.WhiteList,
		accessControl.ClientIPFlowLimit,
		accessControl.ServiceFlowLimit,
		accessControl.ServiceID)

	return err
}

func (repo *serviceRepo) InsertServiceGRPCTCPAccessControl(tx *sqlx.Tx, accessControl *po.ServiceAccessControl, c *gin.Context) (err error) {
	sqlStr := `INSERT INTO gateway_service_access_control(
			service_id, open_auth, black_list, white_list, clientip_flow_limit, service_flow_limit, white_host_name)
			VALUES(?,?,?,?,?,?,?)`
	trace := public.GetGinTraceContext(c)
	if _, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		accessControl.ServiceID,
		accessControl.OpenAuth,
		accessControl.BlackList,
		accessControl.WhiteList,
		accessControl.ClientIPFlowLimit,
		accessControl.ServiceFlowLimit,
		accessControl.WhiteHostName); err != nil {
		return err
	}
	return nil
}
func (repo *serviceRepo) UpdateGRPCTCPAccessControl(tx *sqlx.Tx, accessControl *po.ServiceAccessControl, c *gin.Context) (err error) {
	sqlStr := `UPDATE gateway_service_access_control
			SET open_auth=?, black_list=?, white_list=?, clientip_flow_limit=?, service_flow_limit=?, white_host_name=?
			WHERE service_id = ?`
	trace := public.GetGinTraceContext(c)
	_, err = settings.SqlxLogTxExec(trace, tx, sqlStr,
		accessControl.OpenAuth,
		accessControl.BlackList,
		accessControl.WhiteList,
		accessControl.ClientIPFlowLimit,
		accessControl.ServiceFlowLimit,
		accessControl.WhiteHostName,
		accessControl.ServiceID)

	return err
}

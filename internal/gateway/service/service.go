package service

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/gateway/data/flowcount"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type ServiceRepo interface {
	GetServiceDetail(*po.ServiceInfo, *gin.Context) (*do.ServiceDetail, error)
	GetServiceInfoByID(int64, *gin.Context) (*po.ServiceInfo, error)
	GetServiceInfoByName(string, *gin.Context) (*po.ServiceInfo, error)
	DeleteServiceInfo(*po.ServiceInfo, *gin.Context) error
	GetServiceInfoList(string, int, int, *gin.Context) ([]*po.ServiceInfo, int64, error)
	InsertServiceInfo(*sqlx.Tx, *po.ServiceInfo, *gin.Context) (int64, error)
	UpdateServiceInfo(*sqlx.Tx, *po.ServiceInfo, *gin.Context) error

	AddHTTPDetail(*do.ServiceDetail, *gin.Context) error
	UpdateHTTPDetail(*do.ServiceDetail, *gin.Context) error
	UpdateHTTPRule(*sqlx.Tx, *po.ServiceHTTPRule, *gin.Context) error
	GetServiceHTTPRuleByID(int64, *gin.Context) (*po.ServiceHTTPRule, error)
	GetServiceHTTPRuleByRule(int, string, *gin.Context) (*po.ServiceHTTPRule, error)
	InsertServiceHTTPRule(*sqlx.Tx, *po.ServiceHTTPRule, *gin.Context) error

	GetServiceTCPRuleByID(int64, *gin.Context) (*po.ServiceTCPRule, error)
	InsertServiceTCPRule(*sqlx.Tx, *po.ServiceTCPRule, *gin.Context) error
	UpdateTCPRule(*sqlx.Tx, *po.ServiceTCPRule, *gin.Context) error
	GetServiceTCPRuleByPort(int, *gin.Context) (*po.ServiceTCPRule, error)
	AddTCPDetail(*do.ServiceDetail, *gin.Context) error
	UpdateTCPDetail(*do.ServiceDetail, *gin.Context) error

	GetServiceGRPCRuleByID(int64, *gin.Context) (*po.ServiceGRPCRule, error)
	InsertServiceGRPCRule(*sqlx.Tx, *po.ServiceGRPCRule, *gin.Context) error
	UpdateGRPCRule(*sqlx.Tx, *po.ServiceGRPCRule, *gin.Context) error
	GetServiceGRPCRuleByPort(int, *gin.Context) (*po.ServiceGRPCRule, error)
	AddGRPCDetail(*do.ServiceDetail, *gin.Context) error
	UpdateGrpcDetail(*do.ServiceDetail, *gin.Context) error

	GetServiceLoadBalanceByID(int64, *gin.Context) (*po.ServiceLoadBalance, error)
	InsertServiceHTTPLoadBalance(*sqlx.Tx, *po.ServiceLoadBalance, *gin.Context) error
	UpdateHTTPLoadBalance(*sqlx.Tx, *po.ServiceLoadBalance, *gin.Context) error
	InsertServiceGRPCTCPLoadBalance(*sqlx.Tx, *po.ServiceLoadBalance, *gin.Context) error
	UpdateGRPCTCPLoadBalance(*sqlx.Tx, *po.ServiceLoadBalance, *gin.Context) error

	GetServiceAccessControllerByID(int64, *gin.Context) (*po.ServiceAccessControl, error)
	InsertServiceHTTPAccessControl(*sqlx.Tx, *po.ServiceAccessControl, *gin.Context) error
	UpdateHTTPAccessControl(*sqlx.Tx, *po.ServiceAccessControl, *gin.Context) error
	InsertServiceGRPCTCPAccessControl(*sqlx.Tx, *po.ServiceAccessControl, *gin.Context) error
	UpdateGRPCTCPAccessControl(*sqlx.Tx, *po.ServiceAccessControl, *gin.Context) error
}

type ServiceUseCase struct {
	repo ServiceRepo
}

// NewServiceUseCase 创建一个 Service 用例
func NewServiceUseCase(repo ServiceRepo) *ServiceUseCase {
	return &ServiceUseCase{repo: repo}
}

// GetServiceList 获取服务列表
func (useCase *ServiceUseCase) GetServiceList(serviceDo *do.ServiceListInput, c *gin.Context,
) (serviceListOutput *dto.ServiceListOutput, err error) {
	list, total, err := useCase.repo.GetServiceInfoList(serviceDo.Info,
		serviceDo.PageNo, serviceDo.PageSize, c)

	var outPutList []dto.ServiceListItemOutput
	for _, listItem := range list {
		serviceDetail, err := useCase.repo.GetServiceDetail(listItem, c)
		if err != nil {
			return nil, err
		}

		// service 的三种接入方式
		// 1. http 后缀接入 clusterIP + clusterPort + path
		// 2. http 域名接入
		// 3. tcp、grpc 接入 clusterIP + servicePort
		serviceAddr := "unKnow" // 找不到时的默认提示值
		// HTTP前缀接入
		// 启用 https
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHTTPs == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", settings.ConfBase.ClusterIP,
				settings.ConfBase.ClusterSslPort, serviceDetail.HTTPRule.Rule)
		}
		// 不启用 https
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHTTPs == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", settings.ConfBase.ClusterIP,
				settings.ConfBase.ClusterPost, serviceDetail.HTTPRule.Rule)
		}
		// HTTP 域名接入
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		// TCP
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", settings.ConfBase.ClusterIP,
				serviceDetail.TCPRule.Port)
		}
		// GRPC
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", settings.ConfBase.ClusterIP,
				serviceDetail.GRPCRule.Port)
		}

		// 获取的 ip 字符串用 , 分割，取出里面的 IP 列表
		ipList := strings.Split(serviceDetail.LoadBalance.IPList, ",")
		// 获取流量信息
		serviceCounter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowServicePrefix + listItem.ServiceName)
		if err != nil {
			return nil, err
		}
		outPutItem := dto.ServiceListItemOutput{
			ID:          serviceDetail.Info.ID,
			ServiceName: serviceDetail.Info.ServiceName,
			ServiceDesc: serviceDetail.Info.ServiceDesc,
			ServiceAddr: serviceAddr,
			LoadType:    listItem.LoadType,
			Qps:         int(serviceCounter.QPS),
			Qpd:         int(serviceCounter.TotalCount),
			TotalNode:   len(ipList),
		}
		outPutList = append(outPutList, outPutItem)
	}
	serviceListOutput = &dto.ServiceListOutput{List: outPutList, Total: total}
	return serviceListOutput, nil
}

// DeleteServiceInfo 删除一个服务
func (useCase *ServiceUseCase) DeleteServiceInfo(serviceDeleteInput *dto.ServiceDeleteInput, c *gin.Context) (err error) {

	serviceInfoPo := &po.ServiceInfo{ID: serviceDeleteInput.ID}

	// 根据 ID 查询
	serviceInfoPo, err = useCase.repo.GetServiceInfoByID(serviceInfoPo.ID, c)
	if err != nil {
		return err
	}

	// 修改 id_delete 选项，更新到数据库
	if err = useCase.repo.DeleteServiceInfo(serviceInfoPo, c); err != nil {
		return err
	}

	return nil
}

// ServiceStat 获取流量统计信息
func (useCase *ServiceUseCase) ServiceStat(serviceID int64, c *gin.Context) (todayList, yesterdayList []int64, err error) {
	serviceInfo, err := useCase.repo.GetServiceInfoByID(serviceID, c)
	if err != nil {
		return nil, nil, err
	}
	serviceDetail, err := useCase.repo.GetServiceDetail(serviceInfo, c)
	if err != nil {
		return nil, nil, err
	}
	serviceCounter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
	if err != nil {
		return nil, nil, err
	}
	todayList = make([]int64, 0, 2)
	currentTime := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for i := 0; i <= currentTime.Hour(); i++ {
		dataTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := serviceCounter.GetHourData(dataTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList = make([]int64, 0, 2)
	yesterdayTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	//for i := 0; i <= currentTime.Hour(); i++ {
	for i := 0; i <= 23; i++ {
		dataTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := serviceCounter.GetHourData(dataTime)
		yesterdayList = append(yesterdayList, hourData)
	}
	return
}

// AddHTTP 添加 HTTP 服务逻辑处理
func (useCase *ServiceUseCase) AddHTTP(addHTTP *dto.ServiceAddHTTPInput, c *gin.Context) (err error) {

	// （服务名称不能重复）检查服务在数据库中是否已存在
	if _, err = useCase.repo.GetServiceInfoByName(addHTTP.ServiceName, c); err == nil {
		return errors.New("服务已存在，请更换服务名称")
	}

	// 接入前缀或域名不能重复
	if _, err = useCase.repo.GetServiceHTTPRuleByRule(addHTTP.RuleType, addHTTP.Rule, c); err == nil {
		return errors.New("服务接入前缀或域名已存在")
	}

	// 分别入库
	//useCase.repo.AddHTTPDetail()
	currentTime := time.Now()
	serviceInfo := &po.ServiceInfo{
		ServiceName: addHTTP.ServiceName,
		ServiceDesc: addHTTP.ServiceDesc,

		IsDelete:  0,
		LoadType:  public.LoadTypeHTTP,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	serviceHTTPRule := &po.ServiceHTTPRule{
		RuleType:       addHTTP.RuleType,
		Rule:           addHTTP.Rule,
		NeedHTTPs:      addHTTP.NeedHTTPs,
		NeedStripUri:   addHTTP.NeedStripURI,
		NeedWebsocket:  addHTTP.NeedWebsocket,
		UrlRewrite:     addHTTP.UrlRewrite,
		HeaderTransfor: addHTTP.HeaderTransfor,
	}

	accessControl := &po.ServiceAccessControl{
		OpenAuth:          addHTTP.OpenAuth,
		BlackList:         addHTTP.BlackList,
		WhiteList:         addHTTP.WhiteList,
		ClientIPFlowLimit: addHTTP.ClientFlowLimit,
		ServiceFlowLimit:  addHTTP.ServiceFlowLimit,
	}

	loadBalance := &po.ServiceLoadBalance{
		IPList:                 addHTTP.IpList,
		WeightList:             addHTTP.WeightList,
		UpStreamConnectTimeout: addHTTP.UpstreamConnectTimeout,
		UpStreamHeaderTimeout:  addHTTP.UpstreamHeaderTimeout,
		UpStreamIdleTimeout:    addHTTP.UpstreamIdleTimeout,
		UpStreamMaxIdle:        addHTTP.UpstreamMaxIdle,
	}
	httpDetail := &do.ServiceDetail{
		Info:          serviceInfo,
		HTTPRule:      serviceHTTPRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	// 出现错误在内部回滚并返回错误
	if err = useCase.repo.AddHTTPDetail(httpDetail, c); err != nil {
		return err
	}
	return nil
}

// UpdateHTTP 更新 HTTP 服务
func (useCase *ServiceUseCase) UpdateHTTP(updateHTTP *dto.ServiceUpdateHTTPInput, c *gin.Context) (err error) {

	// 查询服务是否存在
	httpInfo, err := useCase.repo.GetServiceInfoByID(updateHTTP.ID, c)
	if err != nil {
		return errors.New("服务不存在")
	}

	httpDetail, err := useCase.repo.GetServiceDetail(httpInfo, c)
	if err != nil {
		return errors.New("服务不存在")
	}

	// 数据库原始数据
	httpDetail.Info = httpInfo

	// 赋值变更的数据
	// serviceInfo
	httpDetail.Info.ServiceDesc = updateHTTP.ServiceDesc

	// HTTPRule
	httpDetail.HTTPRule.NeedHTTPs = updateHTTP.NeedHTTPs
	httpDetail.HTTPRule.NeedStripUri = updateHTTP.NeedStripURI
	httpDetail.HTTPRule.NeedWebsocket = updateHTTP.NeedWebsocket
	httpDetail.HTTPRule.UrlRewrite = updateHTTP.UrlRewrite
	httpDetail.HTTPRule.HeaderTransfor = updateHTTP.HeaderTransfor

	// LoadBalance
	httpDetail.LoadBalance.RoundType = updateHTTP.RoundType
	httpDetail.LoadBalance.IPList = updateHTTP.IpList
	httpDetail.LoadBalance.WeightList = updateHTTP.WeightList
	httpDetail.LoadBalance.UpStreamConnectTimeout = updateHTTP.UpstreamConnectTimeout
	httpDetail.LoadBalance.UpStreamHeaderTimeout = updateHTTP.UpstreamHeaderTimeout
	httpDetail.LoadBalance.UpStreamIdleTimeout = updateHTTP.UpstreamIdleTimeout
	httpDetail.LoadBalance.UpStreamMaxIdle = updateHTTP.UpstreamMaxIdle

	// AccessControl
	httpDetail.AccessControl.OpenAuth = updateHTTP.OpenAuth
	httpDetail.AccessControl.BlackList = updateHTTP.BlackList
	httpDetail.AccessControl.WhiteList = updateHTTP.WhiteList
	httpDetail.AccessControl.ClientIPFlowLimit = updateHTTP.ClientFlowLimit
	httpDetail.AccessControl.ServiceFlowLimit = updateHTTP.ServiceFlowLimit

	return useCase.repo.UpdateHTTPDetail(httpDetail, c)
}

// GetServiceDetail 根据 ID 获取服务详情
func (useCase *ServiceUseCase) GetServiceDetail(ID int64, c *gin.Context) (serviceDetail *do.ServiceDetail, err error) {
	serviceInfo, err := useCase.repo.GetServiceInfoByID(ID, c)
	if err != nil {
		return nil, err
	}
	serviceDetail = &do.ServiceDetail{Info: serviceInfo}
	serviceDetail, err = useCase.repo.GetServiceDetail(serviceInfo, c)
	if err != nil {
		return nil, err
	}
	return serviceDetail, nil
}

func (useCase *ServiceUseCase) AddGRPC(addGRPC *dto.ServiceAddGrpcInput, c *gin.Context) (err error) {

	// 校验服务名称是否存在
	if _, err = useCase.repo.GetServiceInfoByName(addGRPC.ServiceName, c); err == nil {
		return errors.New("服务已存在，请更换服务名称")
	}

	// 校验端口是否占用
	if _, err = useCase.repo.GetServiceGRPCRuleByPort(addGRPC.Port, c); err == nil {
		return errors.New("端口已占用")
	}

	if _, err = useCase.repo.GetServiceTCPRuleByPort(addGRPC.Port, c); err == nil {
		return errors.New("端口已占用")
	}

	// 入库
	currentTime := time.Now()
	grpcInfo := &po.ServiceInfo{
		ServiceName: addGRPC.ServiceName,
		ServiceDesc: addGRPC.ServiceDesc,
		LoadType:    public.LoadTypeGRPC,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
		IsDelete:    0,
	}

	grpcRule := &po.ServiceGRPCRule{
		Port:           addGRPC.Port,
		HeaderTransfor: addGRPC.HeaderTransfor,
	}

	loadBalance := &po.ServiceLoadBalance{
		RoundType:  addGRPC.RoundType,
		IPList:     addGRPC.IpList,
		WeightList: addGRPC.WeightList,

		ForbidList: addGRPC.ForbidList,
	}

	accessControl := &po.ServiceAccessControl{
		OpenAuth:          addGRPC.OpenAuth,
		WhiteList:         addGRPC.WhiteList,
		BlackList:         addGRPC.BlackList,
		ClientIPFlowLimit: addGRPC.ClientIPFlowLimit,
		ServiceFlowLimit:  addGRPC.ServiceFlowLimit,

		WhiteHostName: addGRPC.WhiteHostName,
	}

	grpcDetail := &do.ServiceDetail{
		Info:          grpcInfo,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	return useCase.repo.AddGRPCDetail(grpcDetail, c)
}

func (useCase *ServiceUseCase) UpdateGRPC(updateGRPC *dto.ServiceUpdateGrpcInput, c *gin.Context) (err error) {
	// 校验服务是否存在
	grpcInfo, err := useCase.repo.GetServiceInfoByID(updateGRPC.ID, c)
	if err != nil {
		return errors.New("服务不存在")
	}

	// 校验服务是否存在
	grpcDetail, err := useCase.repo.GetServiceDetail(grpcInfo, c)
	if err != nil {
		return errors.New("服务不存在")
	}

	// 数据库原始数据
	grpcDetail.Info = grpcInfo

	// 复制变更的数据
	grpcDetail.Info.ServiceDesc = updateGRPC.ServiceDesc

	// GRPCRule
	grpcDetail.GRPCRule.HeaderTransfor = updateGRPC.HeaderTransfor

	// LoadBalance
	grpcDetail.LoadBalance.RoundType = updateGRPC.RoundType
	grpcDetail.LoadBalance.IPList = updateGRPC.IpList
	grpcDetail.LoadBalance.WeightList = updateGRPC.WeightList
	grpcDetail.LoadBalance.ForbidList = updateGRPC.ForbidList

	// AccessControl
	grpcDetail.AccessControl.OpenAuth = updateGRPC.OpenAuth
	grpcDetail.AccessControl.WhiteHostName = updateGRPC.WhiteHostName
	grpcDetail.AccessControl.BlackList = updateGRPC.BlackList
	grpcDetail.AccessControl.WhiteList = updateGRPC.WhiteList
	grpcDetail.AccessControl.ClientIPFlowLimit = updateGRPC.ClientIPFlowLimit
	grpcDetail.AccessControl.ServiceFlowLimit = updateGRPC.ServiceFlowLimit

	return useCase.repo.UpdateGrpcDetail(grpcDetail, c)
}

func (useCase *ServiceUseCase) AddTCP(addTCP *dto.ServiceAddTcpInput, c *gin.Context) (err error) {
	// 校验服务名是否已存在
	_, err = useCase.repo.GetServiceInfoByName(addTCP.ServiceName, c)
	if err == nil {
		return errors.New("服务已存在，请更换服务名称")
	}

	// 校验端口是否被占用
	if _, err = useCase.repo.GetServiceGRPCRuleByPort(addTCP.Port, c); err == nil {
		return errors.New("端口已占用")
	}

	if _, err = useCase.repo.GetServiceTCPRuleByPort(addTCP.Port, c); err == nil {
		return errors.New("端口已占用")
	}

	// 入库
	currentTime := time.Now()
	tcpInfo := &po.ServiceInfo{
		ServiceName: addTCP.ServiceName,
		ServiceDesc: addTCP.ServiceDesc,
		LoadType:    public.LoadTypeTCP,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
		IsDelete:    0,
	}

	tcpRule := &po.ServiceTCPRule{
		Port: addTCP.Port,
	}

	loadBalance := &po.ServiceLoadBalance{
		RoundType:  addTCP.RoundType,
		IPList:     addTCP.IpList,
		WeightList: addTCP.WeightList,

		ForbidList: addTCP.ForbidList,
	}

	accessControl := &po.ServiceAccessControl{
		OpenAuth:          addTCP.OpenAuth,
		WhiteList:         addTCP.WhiteList,
		BlackList:         addTCP.BlackList,
		ClientIPFlowLimit: addTCP.ClientIPFlowLimit,
		ServiceFlowLimit:  addTCP.ServiceFlowLimit,

		WhiteHostName: addTCP.WhiteHostName,
	}

	tcpDetail := &do.ServiceDetail{
		Info:          tcpInfo,
		TCPRule:       tcpRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	return useCase.repo.AddTCPDetail(tcpDetail, c)
}
func (useCase *ServiceUseCase) UpdateTCP(updateTCP *dto.ServiceUpdateTcpInput, c *gin.Context) (err error) {
	// 校验服务是否存在
	tcpInfo, err := useCase.repo.GetServiceInfoByID(updateTCP.ID, c)
	if err != nil {
		return errors.New("服务不存在")
	}

	// 校验服务是否存在
	tcpDetail, err := useCase.repo.GetServiceDetail(tcpInfo, c)
	if err != nil {
		return errors.New("服务不存在")
	}

	// 数据库原始数据
	tcpDetail.Info = tcpInfo

	// 复制变更的数据
	tcpDetail.Info.ServiceDesc = updateTCP.ServiceDesc

	// TCPRule -> 没有可更改的数据

	// LoadBalance
	tcpDetail.LoadBalance.RoundType = updateTCP.RoundType
	tcpDetail.LoadBalance.IPList = updateTCP.IpList
	tcpDetail.LoadBalance.WeightList = updateTCP.WeightList
	tcpDetail.LoadBalance.ForbidList = updateTCP.ForbidList

	// AccessControl
	tcpDetail.AccessControl.OpenAuth = updateTCP.OpenAuth
	tcpDetail.AccessControl.WhiteHostName = updateTCP.WhiteHostName
	tcpDetail.AccessControl.BlackList = updateTCP.BlackList
	tcpDetail.AccessControl.WhiteList = updateTCP.WhiteList
	tcpDetail.AccessControl.ClientIPFlowLimit = updateTCP.ClientIPFlowLimit
	tcpDetail.AccessControl.ServiceFlowLimit = updateTCP.ServiceFlowLimit

	return useCase.repo.UpdateTCPDetail(tcpDetail, c)
}

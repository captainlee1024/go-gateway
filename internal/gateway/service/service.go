package service

import (
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/gateway/settings"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"strings"
)

type ServiceRepo interface {
	GetServiceDetail(*po.ServiceInfo, *gin.Context) (*do.ServiceDetail, error)
	GetServiceInfoList(string, int, int, *gin.Context) ([]*po.ServiceInfo, int64, error)
	GetServiceHTTPRuleByID(int64, *gin.Context) (*po.ServiceHTTPRule, error)
	GetServiceTCPRuleByID(int64, *gin.Context) (*po.ServiceTCPRule, error)
	GetServiceGRPCRuleByID(int64, *gin.Context) (*po.ServiceGRPCRule, error)
	GetServiceLoadBalanceByID(int64, *gin.Context) (*po.ServiceLoadBalance, error)
	GetServiceAccessControllerByID(int64, *gin.Context) (*po.ServiceAccessControl, error)
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
		outPutItem := dto.ServiceListItemOutput{
			ID:          serviceDetail.Info.ID,
			ServiceName: serviceDetail.Info.ServiceName,
			ServiceDesc: serviceDetail.Info.ServiceDesc,
			ServiceAddr: serviceAddr,
			LoadType:    listItem.LoadType,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(ipList),
		}
		outPutList = append(outPutList, outPutItem)
	}
	serviceListOutput = &dto.ServiceListOutput{List: outPutList, Total: total}
	return serviceListOutput, nil
}

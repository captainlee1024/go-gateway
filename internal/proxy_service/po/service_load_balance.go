package po

import (
	"database/sql"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/reverse_proxy/loadbalance"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ServiceLoadBalance struct {
	ID                     int64  `db:"id" json:"id"`
	ServiceID              int64  `db:"service_id" json:"service_id"`
	CheckMethod            int    `db:"check_method" json:"check_method"`
	CheckTimeout           int    `db:"check_timeout" json:"check_timeout"`
	CheckInterval          int    `db:"check_interval" json:"check_interval"`
	RoundType              int    `db:"round_type" json:"round_type"`
	IPList                 string `db:"ip_list" json:"ip_list"`
	WeightList             string `db:"weight_list" json:"weight_list"`
	ForbidList             string `db:"forbid_list" json:"forbid_list"`
	UpStreamConnectTimeout int    `db:"upstream_connect_timeout" json:"upstream_connect_timeout"`
	UpStreamHeaderTimeout  int    `db:"upstream_header_timeout" json:"upstream_header_timeout"`
	UpStreamIdleTimeout    int    `db:"upstream_idle_timeout" json:"upstream_idle_timeout"`
	UpStreamMaxIdle        int    `db:"upstream_max_idle" json:"upstream_max_idle"`
}

func GetServiceLoadBalanceByID(ID int64, c *gin.Context,
) (loadBalance *ServiceLoadBalance, err error) {
	db, err := settings.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	loadBalance = new(ServiceLoadBalance)
	trace := public.ProxyGetGinTraceContext(c)
	loadBalanceSqlStr := `SELECT * FROM gateway_service_load_balance WHERE service_id=?`
	if err = settings.SqlxLogGet(trace, db, loadBalance, loadBalanceSqlStr, ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return loadBalance, nil
}

var LoadBalancerHandler *LoadBalancer

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   map[string]*LoadBalancerItem{},
		LoadBalanceSlice: []*LoadBalancerItem{},
		Locker:           sync.RWMutex{},
	}
}

type LoadBalancer struct {
	// 当服务多的时候使用map方便，但是需要加锁
	LoadBalanceMap map[string]*LoadBalancerItem
	// 当服务比较少的时候使用 slice 便利，减少锁的开销
	LoadBalanceSlice []*LoadBalancerItem
	Locker           sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBalance loadbalance.LoadBalance
	ServiceName string
}

func init() {
	LoadBalancerHandler = NewLoadBalancer()
}

func (lbr *LoadBalancer) GetLoadBalancer(service *ServiceDetail) (loadbalance.LoadBalance, error) {
	for _, lbrItem := range lbr.LoadBalanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName {
			return lbrItem.LoadBalance, nil
		}
	}

	schema := "http" // 默认 http
	if service.HTTPRule.NeedHTTPs == 1 {
		schema = "https" // 如果开启 https 支持，设置协议为 https
	}

	//prefix := ""
	//if service.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
	//	prefix = service.HTTPRule.Rule
	//}
	ipList := service.LoadBalance.GetIpListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	ipConf := map[string]string{}
	for ipIndex, ipItem := range ipList {
		ipConf[ipItem] = weightList[ipIndex]
	}
	mConf, err := loadbalance.NewCheckConf(
		// 基于 http rule 设置 schema 和接入类型
		// 负载的 IP和权重
		fmt.Sprintf("%s://%s", schema, "%s"), ipConf)
	// nil)
	if err != nil {
		return nil, err
	}

	//lb := loadbalance.FactorWithConf(loadbalance.LbWeightRoundRobin, mConf)
	lb := loadbalance.FactorWithConf(loadbalance.LbType(service.LoadBalance.RoundType), mConf)

	// save to map and slice
	loadbalanceItem := &LoadBalancerItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	}
	lbr.LoadBalanceSlice = append(lbr.LoadBalanceSlice, loadbalanceItem)
	lbr.Locker.Lock()
	defer lbr.Locker.Unlock()
	lbr.LoadBalanceMap[service.Info.ServiceName] = loadbalanceItem
	return lb, nil
}

func (l *ServiceLoadBalance) GetIpListByModel() []string {
	return strings.Split(l.IPList, ",")
}

func (l *ServiceLoadBalance) GetWeightListByModel() []string {
	return strings.Split(l.WeightList, ",")
}

var TransportHandler *Transport

type TransportItem struct {
	Transport   *http.Transport
	ServiceName string
}

type Transport struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

func NewTransport() *Transport {
	return &Transport{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

func init() {
	TransportHandler = NewTransport()
}

func (t *Transport) GetTransport(service *ServiceDetail) (*http.Transport, error) {
	// 判断该服务的连接池是否已经存在了
	for _, transportItem := range t.TransportSlice {
		if transportItem.ServiceName == service.Info.ServiceName {
			return transportItem.Transport, nil
		}
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(service.LoadBalance.UpStreamConnectTimeout) * time.Second, // 连接超时时间
		}).DialContext,
		MaxIdleConns:          service.LoadBalance.UpStreamMaxIdle,                                  // 最大空闲连接
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpStreamIdleTimeout) * time.Second, // 空闲超时时间
		TLSHandshakeTimeout:   time.Duration(service.LoadBalance.CheckTimeout) * time.Second,        // tls 握手超时时间
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpStreamHeaderTimeout) * time.Second,
	}

	// save to map and slice
	transportItem := &TransportItem{
		ServiceName: service.Info.ServiceName,
		Transport:   transport,
	}
	t.TransportSlice = append(t.TransportSlice, transportItem)
	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransportMap[service.Info.ServiceName] = transportItem
	return transport, nil
}

package flowlimit

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimiterHandler *FlowLimiter

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Locker:           sync.RWMutex{},
	}
}

type FlowLimiter struct {
	// 当服务多的时候使用map方便，但是需要加锁
	FlowLimiterMap map[string]*FlowLimiterItem
	// 当服务比较少的时候使用 slice 便利，减少锁的开销
	FlowLimiterSlice []*FlowLimiterItem
	Locker           sync.RWMutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

func (counter *FlowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	for _, item := range counter.FlowLimiterSlice {
		if item.ServiceName == serviceName {
			return item.Limiter, nil
		}
	}

	// 如果不存在，初始化
	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))

	item := &FlowLimiterItem{
		ServiceName: serviceName,
		Limiter:     newLimiter,
	}
	counter.FlowLimiterSlice = append(counter.FlowLimiterSlice, item)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.FlowLimiterMap[serviceName] = item
	return newLimiter, nil
}

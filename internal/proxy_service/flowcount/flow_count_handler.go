package flowcount

import (
	"sync"
	"time"
)

var FlowCounterHandler *FlowCounter

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

type FlowCounter struct {
	// 当服务多的时候使用map方便，但是需要加锁
	RedisFlowCountMap map[string]*RedisFlowCountService
	// 当服务比较少的时候使用 slice 便利，减少锁的开销
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

//type FlowCounterItem struct {
//	LoadBalance loadbalance.LoadBalance
//	ServiceName string
//}

func init() {
	FlowCounterHandler = NewFlowCounter()
}

func (counter *FlowCounter) GetCounter(serviceName string) (*RedisFlowCountService, error) {
	for _, item := range counter.RedisFlowCountSlice {
		if item.AppID == serviceName {
			return item, nil
		}
	}

	// 如果不存在，初始化
	newCounter := NewRedisFlowCountService(serviceName, time.Second)

	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCountMap[serviceName] = newCounter
	return newCounter, nil
}

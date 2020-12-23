package dto

type PanelGroupDataOutput struct {
	ServiceNum      int64 `json:"service_num"`       // 服务总数
	AppNumber       int64 `json:"app_num"`           // 租户总数
	TodayRequestNum int64 `json:"today_request_num"` // 今日请求总数
	CurrentQps      int64 `json:"current_qps"`       // 当前总 QPS
}

type FlowStatOutput struct {
	Today     []int64 `json:"today"`     // 今日请求数
	Yesterday []int64 `json:"yesterday"` // 昨日请求数
}

type DashboardServiceStatOutput struct {
	Legend []string                         `json:"legend"` // 服务类型列表
	Data   []DashboardServiceStatItemOutput `json:"data"`   // 各服务类型数据列表
}

type DashboardServiceStatItemOutput struct {
	Name  string `json:"name"`  // 服务类型
	Value string `json:"value"` // 服务数量
}

package controller

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/data/flowcount"
	"time"

	v1 "github.com/captainlee1024/go-gateway/api/gateway/v1"
	"github.com/captainlee1024/go-gateway/internal/gateway/data"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

func DashboardRegister(group *gin.RouterGroup) {
	repo := data.NewDashboardRepo()
	useCase := service.NewDashboardUseCase(repo)
	dashboardController := NewDashboardController(useCase)

	group.GET("/panel_group_data", dashboardController.PanelGroupData)
	group.GET("/flow_stat", dashboardController.FlowStat)
	group.GET("/service_stat", dashboardController.DashboardServiceStat)
}

type DashboardController struct {
	v1.UnimplementedDashboardServer
	useCase *service.DashboardUseCase
}

func NewDashboardController(useCase *service.DashboardUseCase) *DashboardController {
	return &DashboardController{useCase: useCase}
}

// AppStat godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [GET]
func (d *DashboardController) PanelGroupData(c *gin.Context) {

	// 查询数据
	serviceNum, appNum, err := d.useCase.PanelGroupData(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	totalCounter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// 组装返回数据
	outPut := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNumber:       appNum,
		TodayRequestNum: totalCounter.TotalCount, // 实现代理的时候接入数据
		CurrentQps:      totalCounter.QPS,        // 实现代理的时候接入数据
	}
	middleware.ResponseSuccess(c, outPut)
}

// FlowStat godoc
// @Summary 总流量统计
// @Description 总流量统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.FlowStatOutput} "success"
// @Router /dashboard/flow_stat [GET]
func (d *DashboardController) FlowStat(c *gin.Context) {
	//var todayList []int64
	//for i := 0; i <= time.Now().Hour(); i++ {
	//	todayList = append(todayList, int64(i*300))
	//}
	//
	//var yesterdayList []int64
	//for i := 0; i <= 23; i++ {
	//	yesterdayList = append(yesterdayList, int64(rand.Intn(10)*600))
	//}

	totalCOunter, err := flowcount.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	var todayList []int64
	currentTime := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for i := 0; i <= currentTime.Hour(); i++ {
		dataTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := totalCOunter.GetHourData(dataTime)
		todayList = append(todayList, hourData)
	}

	var yesterdayList []int64
	yesterdayTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= currentTime.Hour(); i++ {
		dataTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := totalCOunter.GetHourData(dataTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	outPut := &dto.FlowStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}
	middleware.ResponseSuccess(c, outPut)
}

// FlowStat godoc
// @Summary 服务占比统计
// @Description 服务占比统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashboardServiceStatOutput} "success"
// @Router /dashboard/service_stat [GET]
func (d *DashboardController) DashboardServiceStat(c *gin.Context) {

	statListDo, err := d.useCase.GetServiceStat(c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	outPut := &dto.DashboardServiceStatOutput{}

	for _, item := range statListDo {
		name, ok := public.LoadTypeMap[item.Name]
		if !ok {
			middleware.ResponseError(c, 2002, errors.New("load_type not found"))
			return
		}
		outItem := dto.DashboardServiceStatItemOutput{
			Name:     name,
			LoadType: item.Name,
			Value:    item.Value,
		}
		outPut.Legend = append(outPut.Legend, outItem.Name)
		outPut.Data = append(outPut.Data, outItem)
	}

	middleware.ResponseSuccess(c, outPut)
}

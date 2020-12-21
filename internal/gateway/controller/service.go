package controller

import (
	"errors"
	v1 "github.com/captainlee1024/go-gateway/api/gateway/v1"
	"github.com/captainlee1024/go-gateway/internal/gateway/data"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func ServiceRegister(group *gin.RouterGroup) {
	repo := data.NewServiceRepo()
	useCase := service.NewServiceUseCase(repo)
	serviceController := NewServiceController(useCase)

	group.GET("/service_list", serviceController.ServiceList)
	group.GET("/service_delete", serviceController.ServiceDelete)
	group.GET("/service_detail", serviceController.ServiceDetail)
	group.GET("/service_stat", serviceController.ServiceStat)
	group.POST("/service_add_http", serviceController.ServiceAddHTTP)
	group.POST("/service_update_http", serviceController.ServiceUpdateHTTP)
}

type ServiceController struct {
	v1.UnimplementedServiceServer
	useCase *service.ServiceUseCase
}

func NewServiceController(useCase *service.ServiceUseCase) *ServiceController {
	return &ServiceController{useCase: useCase}
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// dto -> do
	serviceListDo := &do.ServiceListInput{
		Info:     params.Info,
		PageNo:   params.PageNo,
		PageSize: params.PageSize,
	}

	serviceListOutput := &dto.ServiceListOutput{}
	serviceListOutput, err := s.useCase.GetServiceList(serviceListDo, c)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, serviceListOutput)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept json
// @Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [GET]
func (s *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	err := s.useCase.DeleteServiceInfo(params, c)
	if err != nil {
		// 如果是找不到，可以返回提示信息，找不到对应的 Service
		if errors.Is(err, data.ErrServiceNotExit) {
			middleware.ResponseError(c, 2004, data.ErrServiceNotExit)
			return
		}
		// TODO: 工程化 －> 统一风格的返回状态码和对应提示信息
		// 如果是其他错误，最好不要直接对外暴露错误，最好是返回统一的服务繁忙之类的
		middleware.ResponseError(c, 2005, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}

// ServiceAddHTTP godoc
// @Summary 添加 HTTP 服务
// @Description 添加 HTTP 服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept json
// @Produce json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "Success"
// @Router /service/service_add_http [POST]
func (s *ServiceController) ServiceAddHTTP(c *gin.Context) {

	// 参数基本校验
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 参数校验
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2007, errors.New("IP列表与权重列表数量不匹配"))
		return
	}

	// 调用 service 进行逻辑处理
	if err := s.useCase.AddHTTP(params, c); err != nil {
		middleware.ResponseError(c, 2008, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}

// ServiceUpdateHTTP godoc
// @Summary 修改 HTTP 服务
// @Description 修改 HTTP 服务
// @Tags 服务管理
// @ID /service/service_update_http
// @Accept json
// @Produce json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "Success"
// @Router /service/service_update_http [POST]
func (s *ServiceController) ServiceUpdateHTTP(c *gin.Context) {

	// 参数校验（基本校验）
	params := &dto.ServiceUpdateHTTPInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2009, err)
		return
	}

	// 参数校验（加强校验）
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2010, errors.New("IP列表与权重列表数量不匹配"))
		return
	}

	// 更新 HTTP 服务信息
	if err := s.useCase.UpdateHTTP(params, c); err != nil {
		middleware.ResponseError(c, 2011, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service_detail
// @Accept json
// @Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=do.ServiceDetail} "success"
// @Router /service/service_detail [GET]
func (s *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2012, err)
		return
	}
	serviceDetail, err := s.useCase.GetServiceDetail(params.ID, c)
	if err != nil {
		middleware.ResponseError(c, 2013, err)
		return
	}
	middleware.ResponseSuccess(c, serviceDetail)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/service_stat [GET]
func (s *ServiceController) ServiceStat(c *gin.Context) {
	params := &dto.ServiceStatInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2014, err)
		return
	}

	var todayList []int64
	for i := 0; i <= time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}

	var yesterdayList []int64
	for i := 0; i <= 23; i++ {
		yesterdayList = append(yesterdayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

package controller

import (
	v1 "github.com/captainlee1024/go-gateway/api/gateway/v1"
	"github.com/captainlee1024/go-gateway/internal/gateway/data"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"time"
)

func AppRegister(group *gin.RouterGroup) {
	repo := data.NewAppRepo()
	userCase := service.NewAppUseCase(repo)
	appController := NewAppController(userCase)

	group.GET("/app_list", appController.AppList)
	group.GET("/app_detail", appController.AppDetail)
	group.GET("/app_stat", appController.AppStat)
	group.GET("/app_delete", appController.AppDelete)
	group.POST("/app_add", appController.AppAdd)
	group.POST("/app_update", appController.AppUpdate)
}

type AppController struct {
	v1.UnimplementedAppServer
	useCase *service.AppUseCase
}

func NewAppController(useCase *service.AppUseCase) *AppController {
	return &AppController{useCase: useCase}
}

// AppList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.AppListOutput} "success"
// @Router /app/app_list [GET]
func (a *AppController) AppList(c *gin.Context) {

	// 参数校验
	params := &dto.AppListInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 查询
	outPut, err := a.useCase.AppList(params, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 结果返回前端
	middleware.ResponseSuccess(c, outPut)
}

// AppDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "ID"
// @Success 200 {object} middleware.Response{data=dto.AppDetailOutput} "success"
// @Router /app/app_detail [GET]
func (a *AppController) AppDetail(c *gin.Context) {
	// 参数校验
	params := &dto.AppDetailInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 查询详情
	outPut, err := a.useCase.AppDetail(params, c)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	// 结果返回前端
	middleware.ResponseSuccess(c, outPut)
}

// AppDetail godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [GET]
func (a *AppController) AppDelete(c *gin.Context) {
	// 参数校验
	params := &dto.AppDeleteInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	// 数据删除
	if err := a.useCase.AppDelete(params, c); err != nil {
		middleware.ResponseError(c, 2005, err)
		return
	}

	// 返回
	middleware.ResponseSuccess(c, "")
}

// AppAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.AppAddInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [POST]
func (a *AppController) AppAdd(c *gin.Context) {
	// 参数校验
	params := &dto.AppAddInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 校验用户是否输入了秘钥，没有就生成一个秘钥
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}

	// dto -> do 这里 do 复用的 dto

	// 进行数据库校验并入库
	if err := a.useCase.AppAdd(params, c); err != nil {
		middleware.ResponseError(c, 2007, err)
		return
	}

	// 返回数据
	middleware.ResponseSuccess(c, "")
}

// AppUpdate godoc
// @Summary 租户修改
// @Description 租户修改
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.AppUpdateInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [POST]
func (a *AppController) AppUpdate(c *gin.Context) {
	// 参数校验
	params := &dto.AppUpdateInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2008, err)
		return
	}

	// 秘钥在修改的时候会默认读取之前的秘钥填充到前段传过来，这里无需校验

	// 变更入库
	if err := a.useCase.AppUpdate(params, c); err != nil {
		middleware.ResponseError(c, 2009, err)
		return
	}

	// 返回
	middleware.ResponseSuccess(c, "")
}

// AppStat godoc
// @Summary 租户流量统计
// @Description 租户流量统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dto.AppStatOutput} "success"
// @Router /app/app_stat [GET]
func (a *AppController) AppStat(c *gin.Context) {
	params := &dto.AppStatInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2010, err)
		return
	}

	// 在实现代理的时候接入

	var todayList []int64
	for i := 0; i <= time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}

	var yesterdayList []int64
	for i := 0; i <= 23; i++ {
		yesterdayList = append(yesterdayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.AppStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

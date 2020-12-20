package controller

import (
	v1 "github.com/captainlee1024/go-gateway/api/gateway/v1"
	"github.com/captainlee1024/go-gateway/internal/gateway/data"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/gin-gonic/gin"
)

func ServiceRegister(group *gin.RouterGroup) {
	repo := data.NewServiceRepo()
	useCase := service.NewServiceUseCase(repo)
	serviceController := NewServiceController(useCase)

	group.GET("/service_list", serviceController.ServiceList)
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
		middleware.ResponseError(c, 2002, nil)
		return
	}

	middleware.ResponseSuccess(c, serviceListOutput)
}

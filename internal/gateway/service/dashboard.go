package service

import (
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/gin-gonic/gin"
)

type DashboardRepo interface {
	GetServiceNum(*gin.Context) (int64, error)
	GetAppNum(*gin.Context) (int64, error)
	GetServiceStat(*gin.Context) ([]*do.DashboardServiceStat, error)
}

type DashboardUseCase struct {
	repo DashboardRepo
}

func NewDashboardUseCase(repo DashboardRepo) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

// PanelGroupData 获取大盘的服务数和租户数指标
func (useCase *DashboardUseCase) PanelGroupData(c *gin.Context) (serviceNum, appNum int64, err error) {
	serviceNum, err = useCase.repo.GetServiceNum(c)
	if err != nil {
		return 0, 0, err
	}

	appNum, err = useCase.repo.GetAppNum(c)
	if err != nil {
		return 0, 0, err
	}

	return serviceNum, appNum, err
}

func (useCase *DashboardUseCase) GetServiceStat(c *gin.Context) (serviceStatList []*do.DashboardServiceStat, err error) {
	serviceStatList = make([]*do.DashboardServiceStat, 0, 2)
	serviceStatList, err = useCase.repo.GetServiceStat(c)
	if err != nil {
		return nil, err
	}
	return serviceStatList, err
}

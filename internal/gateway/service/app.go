package service

import (
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/data/flowcount"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/po"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/gin"
	"time"
)

type AppRepo interface {
	GetAppList(string, int, int, *gin.Context) ([]*po.App, int64, error)
	GetAppDetailByID(int64, *gin.Context) (*po.App, error)
	DeleteAppByID(int64, *gin.Context) error
	InsertApp(*po.App, *gin.Context) (int64, error)
	GetAppDetailByAppID(string, *gin.Context) (*po.App, error)
	UpdateApp(*po.App, *gin.Context) error
}

type AppUseCase struct {
	repo AppRepo
}

func NewAppUseCase(repo AppRepo) *AppUseCase {
	return &AppUseCase{repo: repo}
}

// AppList 获取租户列表
func (useCase *AppUseCase) AppList(appListInput *dto.AppListInput, c *gin.Context) (outPutList *dto.AppListOutput, err error) {

	outPutList = new(dto.AppListOutput)

	// 分页模糊查询
	appList, total, err := useCase.repo.GetAppList(appListInput.Info, appListInput.PageNo, appListInput.PageSize, c)
	if err != nil {
		return nil, err
	}

	// 组装数据返回
	for _, item := range appList {
		appCounter, err := flowcount.FlowCounterHandler.GetCounter(
			public.FlowAppPrefix + item.AppID)
		if err != nil {
			return nil, err
		}
		outItem := dto.AppListItemOutput{
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPs: item.WhiteIPs,
			Qps:      item.Qps,
			Qpd:      item.Qpd,
			ID:       item.ID,
			RealQps:  appCounter.QPS,
			RealQpd:  appCounter.TotalCount,
		}
		outPutList.List = append(outPutList.List, outItem)
	}
	outPutList.Total = total
	return outPutList, nil
}

// AppDetail 获取租户详细信息
func (useCase *AppUseCase) AppDetail(appDetailInput *dto.AppDetailInput, c *gin.Context) (outPut *dto.AppDetailOutput, err error) {

	// 查询
	appDetail, err := useCase.repo.GetAppDetailByID(appDetailInput.ID, c)
	if err != nil {
		return nil, err
	}

	// 组装返回的数据
	outPut = &dto.AppDetailOutput{
		ID:       appDetail.ID,
		AppID:    appDetail.AppID,
		Name:     appDetail.Name,
		Secret:   appDetail.Secret,
		Qps:      appDetail.Qps,
		Qpd:      appDetail.Qpd,
		WhiteIPs: appDetail.WhiteIPs,
		RealQpd:  0,
		RealQps:  0,
	}
	return outPut, nil
}

// AppDelete 删除租户
func (useCase *AppUseCase) AppDelete(appDeleteInput *dto.AppDeleteInput, c *gin.Context) (err error) {
	return useCase.repo.DeleteAppByID(appDeleteInput.ID, c)
}

// AppAdd 添加一个租户
func (useCase *AppUseCase) AppAdd(appAddInput *dto.AppAddInput, c *gin.Context) (err error) {
	// 验证 app_id 是否被占用
	if _, err = useCase.repo.GetAppDetailByAppID(appAddInput.AppID, c); err == nil {
		return errors.New("租户ID被占用，请重新输入")
	}

	currentTime := time.Now()
	appPo := &po.App{
		AppID:     appAddInput.AppID,
		Name:      appAddInput.Name,
		Secret:    appAddInput.Secret,
		WhiteIPs:  appAddInput.WhiteIPs,
		Qps:       appAddInput.Qps,
		Qpd:       appAddInput.Qpd,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		IsDelete:  0,
	}

	_, err = useCase.repo.InsertApp(appPo, c)
	return err
}

// AppUpdate 更改租户信息
func (useCase *AppUseCase) AppUpdate(appUpdateInput *dto.AppUpdateInput, c *gin.Context) (err error) {

	appPo := &po.App{
		ID:        appUpdateInput.ID,
		AppID:     appUpdateInput.AppID,
		Name:      appUpdateInput.Name,
		Secret:    appUpdateInput.Secret,
		Qps:       appUpdateInput.Qps,
		Qpd:       appUpdateInput.Qpd,
		WhiteIPs:  appUpdateInput.WhiteIPs,
		UpdatedAt: time.Now(),
	}

	return useCase.repo.UpdateApp(appPo, c)
}

// AppStat 获取租户流量统计信息
func (useCase *AppUseCase) AppStat(ID int64, c *gin.Context) (todayList, yesterdayList []int64, err error) {
	appInfo, err := useCase.repo.GetAppDetailByID(ID, c)
	if err != nil {
		return nil, nil, err
	}
	appCounter, err := flowcount.FlowCounterHandler.GetCounter(
		public.FlowAppPrefix + appInfo.AppID)
	if err != nil {
		return nil, nil, err
	}
	todayList = make([]int64, 0, 2)
	currentTime := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for i := 0; i <= currentTime.Hour(); i++ {
		dataTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := appCounter.GetHourData(dataTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList = make([]int64, 0, 2)
	yesterdayTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dataTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, loc)
		hourData, _ := appCounter.GetHourData(dataTime)
		yesterdayList = append(yesterdayList, hourData)
	}
	return todayList, yesterdayList, nil
}

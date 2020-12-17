package controller

import (
	"encoding/json"
	"fmt"
	v1 "github.com/captainlee1024/go-gateway/api/gateway/v1"
	"github.com/captainlee1024/go-gateway/internal/gateway/data"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminRegister(group *gin.RouterGroup) {
	repo := data.NewAdminRepo()
	useCase := service.NewAdminUseCase(repo)
	adminController := NewAdminController(useCase)

	group.GET("/admin_info", adminController.AdminInfo)
	group.POST("/change_pwd", adminController.ChangePwd)
}

type AdminController struct {
	v1.UnimplementedAdminServer
	useCase *service.AdminUseCase
}

func NewAdminController(useCase *service.AdminUseCase) *AdminController {
	return &AdminController{useCase: useCase}
}

// AdminInfo 管理员信息
// @Summary 管理员信息
// @Description 从登陆的 Session 中获取管理员信息
// @Tags 管理员接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [GET]
func (admin *AdminController) AdminInfo(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.KeyAdminSessionInfo)
	adminSessionInfo := &do.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	adminInfo := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Username:     adminSessionInfo.Username,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, adminInfo)
}

// ChangePwd 修改密码
// @Summary 修改密码
// @Description 修改已登录账户的密码
// @Tags 管理员接口
// @Accept application/json
// @Produce application/json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [POST]
func (admin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// dto -> do
	changePwdDo := &do.ChangePwd{
		Password: params.Password,
	}

	if err := admin.useCase.ChangePwd(changePwdDo, c); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}

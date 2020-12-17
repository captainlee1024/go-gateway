package controller

import (
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

func AdminLoginRegister(group *gin.RouterGroup) {
	repo := data.NewAdminLoginRepo()
	useCase := service.NewAdminLoginUseCase(repo)
	adminLogin := NewAdminLoginController(useCase)

	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLogOut)
}

type AdminLoginController struct {
	v1.UnimplementedAdminLoginServer
	useCase *service.AdminLoginUseCase
}

func NewAdminLoginController(useCase *service.AdminLoginUseCase) *AdminLoginController {
	return &AdminLoginController{useCase: useCase}
}

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @Accept application/json
// @Produce application/json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [POST]
func (a *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	adminLoginDo := &do.AdminLogin{
		Username: params.Username,
		Password: params.Password,
	}
	adminLoginDo, err := a.useCase.LoginCheck(adminLoginDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	outPut := &dto.AdminLoginOutput{Token: adminLoginDo.Username}
	middleware.ResponseSuccess(c, outPut)
}

// AdminLogOut 管理员退出登录
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [GET]
func (a *AdminLoginController) AdminLogOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.KeyAdminSessionInfo)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}

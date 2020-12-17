package service

import (
	"encoding/json"
	"fmt"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminRepo interface {
	GetAdmin(int, *gin.Context) (*do.AdminLogin, error)
	UpdatePassword(*do.AdminLogin, *gin.Context) error
}

type AdminUseCase struct {
	repo AdminRepo
}

func NewAdminUseCase(repo AdminRepo) *AdminUseCase {
	return &AdminUseCase{repo: repo}
}

func (a *AdminUseCase) ChangePwd(changePwdDo *do.ChangePwd, c *gin.Context) (err error) {
	sess := sessions.Default(c)
	sessionInfo := sess.Get(public.KeyAdminSessionInfo)
	adminSessionInfo := &do.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		return err
	}

	adminInfoDo, err := a.repo.GetAdmin(adminSessionInfo.ID, c)
	if err != nil {
		return err
	}

	newPassword := public.GenSaltPassword(adminInfoDo.Salt, changePwdDo.Password)
	adminInfoDo.Password = newPassword
	return a.repo.UpdatePassword(adminInfoDo, c)
}

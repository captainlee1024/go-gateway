package service

import (
	"encoding/json"
	"errors"
	"github.com/captainlee1024/go-gateway/internal/gateway/do"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

// AdminLoginRepo 存储接口
type AdminLoginRepo interface {
	GetAdmin(*do.AdminLogin, *gin.Context) (*do.AdminLogin, error)
}

// AdminLoginUseCase is .
type AdminLoginUseCase struct {
	repo AdminLoginRepo
}

func NewAdminLoginUseCase(repo AdminLoginRepo) *AdminLoginUseCase {
	return &AdminLoginUseCase{repo: repo}
}

// AdminLogin 进行登录验证
func (useCase *AdminLoginUseCase) LoginCheck(loginDo *do.AdminLogin, c *gin.Context) (admin *do.AdminLogin, err error) {
	admin = &do.AdminLogin{}

	admin, err = useCase.repo.GetAdmin(loginDo, c)
	if err != nil {
		return nil, err
	}
	if !loginDo.PasswordCheck(admin.Salt, admin.Password) {
		return nil, errors.New("密码错误！")
	}

	sessionInfo := &do.AdminSessionInfo{
		ID:        admin.ID,
		Username:  admin.Username,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessionInfo)
	if err != nil {
		return nil, err
	}
	sess := sessions.Default(c)
	sess.Set(public.KeyAdminSessionInfo, string(sessBts))
	sess.Save()

	return admin, nil
}

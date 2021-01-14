package controller

import (
	"encoding/base64"
	"errors"
	v1 "github.com/captainlee1024/go-gateway/api/gateway/v1"
	"github.com/captainlee1024/go-gateway/internal/gateway/dto"
	"github.com/captainlee1024/go-gateway/internal/gateway/middleware"
	"github.com/captainlee1024/go-gateway/internal/gateway/service"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"github.com/captainlee1024/go-gateway/internal/proxy_service/po"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func OauthRegister(group *gin.RouterGroup) {
	useCase := service.NewOauthUseCase()
	oauthController := NewOauthController(useCase)
	group.POST("/tokens", oauthController.Tokens)
}

type OauthController struct {
	v1.UnimplementedOauthServer
	useCase *service.OauthUseCase
}

func NewOauthController(useCase *service.OauthUseCase) *OauthController {
	return &OauthController{useCase: useCase}
}

// Tokens godoc
// @Summary 获取token
// @Description 获取token
// @Tags OAUTH
// @Accept application/json
// @Produce application/json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [POST]
func (o *OauthController) Tokens(c *gin.Context) {
	params := &dto.TokensInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 取出 app_id secret
	// 生成appList
	// 匹配 appID
	// 生成 JWT
	// 返回 Output
	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		middleware.ResponseError(c, 2002, errors.New("用户名或密码格式错误！"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	oauth := strings.Split(string(appSecret), ":")
	if len(oauth) != 2 {
		middleware.ResponseError(c, 2004, errors.New("用户名或密码格式错误！"))
	}

	appList := po.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		if oauth[0] == appInfo.AppID && oauth[1] == appInfo.Secret {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(public.JWTExpires * time.Second).In(loc).Unix(),
			}
			token, err := public.JWTEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 2005, err)
				return
			}

			outPut := &dto.TokensOutput{
				ExpiresIn:   public.JWTExpires,
				TokenType:   "Bearer",
				AccessToken: token,
				Scope:       "read_write",
			}
			middleware.ResponseSuccess(c, outPut)
			return
		}
	}
	middleware.ResponseError(c, 2006, errors.New("未找到APP信息！"))
}

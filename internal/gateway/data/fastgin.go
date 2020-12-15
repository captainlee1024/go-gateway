package data

import (
	"github.com/captainlee1024/go-gateway/internal/gateway/data/redis"
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	red "github.com/garyburd/redigo/redis"

	"github.com/gin-gonic/gin"
)

// SetAToken 设置 Token
func SetAToken(ID, token string, c *gin.Context) (err error) {
	trace := public.GetGinTraceContext(c)

	//_, err = redis.ConfDo(trace, "default", "SET", redis.GetRedisKey(redis.KeyFastGinJWTSetPrefix+ID), token, time.Hour*time.Duration(settings.GetIntConf("base.auth.jwt_expire")))
	_, err = redis.ConfDo(trace, "default", "SET", redis.GetRedisKey(redis.KeyFastGinJWTSetPrefix+ID), token)

	if err != nil {
		return err
	}
	return nil
}

// GetAToken 获取 Token
func GetAToken(ID string, c *gin.Context) (token string, err error) {
	trace := public.GetGinTraceContext(c)

	_token, err := redis.ConfDo(trace, "default", "GET", redis.GetRedisKey(redis.KeyFastGinJWTSetPrefix+ID))
	if err != nil {
		return "", err
	}
	token, err = red.String(_token, err)
	if err != nil {
		return "", nil
	}
	return
}

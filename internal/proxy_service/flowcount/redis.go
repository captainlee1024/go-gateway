package flowcount

import (
	"github.com/captainlee1024/go-gateway/internal/proxy_service/settings"
	"github.com/garyburd/redigo/redis"
)

func RedisConfPipline(pip ...func(c redis.Conn)) error {

	c, err := settings.ConnFactory("default")
	if err != nil {
		return err
	}
	defer c.Close()
	for _, f := range pip {
		f(c)
	}
	c.Flush()
	return nil
}

func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	c, err := settings.ConnFactory("default")
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.Do(commandName, args...)
}

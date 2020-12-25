package redis_server

import (
	"sec-dev-in-action-src/honeypot/server/logger"
	"sec-dev-in-action-src/honeypot/server/services/redis-server/redis"
)

func StartRedis(addr string, flag bool) error {
	logger.Log.Warningf("start redis service on %v", addr)
	err := redis.Run(addr, flag)
	return err
}

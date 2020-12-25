package models

import (
	"fmt"
	"time"

	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/vars"

	"github.com/go-redis/redis"
)

func init() {
	InitRedis()
}

func InitRedis() {
	var err error
	vars.RedisInstance, err = NewRedisClient(vars.RedisConf.Host, vars.RedisConf.Port,
		vars.RedisConf.Db, vars.RedisConf.Password)

	if err != nil {
		logger.Log.Errorf("connect redis failed, err: %v", err)
	}
}

func NewRedisClient(host string, port int, db int, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%v:%v", host, port),
		Password:    password,    // no password set
		DB:          db,          // use default DB
		ReadTimeout: time.Minute, // set timeout value = 60
	})

	_, err := client.Ping().Result()
	return client, err
}

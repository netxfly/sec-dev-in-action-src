package redis

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// redis server
var r = Default()

// redis client
var c = redis.NewClient(&redis.Options{
	Addr: ":6379",
})

func init() {
	go r.Run(":6379")
}

func TestPingCommand(t *testing.T) {
	s, err := c.Ping().Result()
	assert.Equal(t, "PONG", s)
	assert.NoError(t, err)

	pingCmd := redis.NewStringCmd("ping", "Hello,", "redis server!")
	c.Process(pingCmd)
	s, err = pingCmd.Result()
	assert.Equal(t, "Hello, redis server!", s)
	assert.NoError(t, err)
}

func TestSetCommand(t *testing.T) {
	s, err := c.Set("k", "v", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set("k2", nil, 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set("k3", "v", 1*time.Hour).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
}

func TestGetCommand(t *testing.T) {
	s, err := c.Get("k").Result()
	assert.Equal(t, "v", s)
	assert.NoError(t, err)
}

func TestDelCommand(t *testing.T) {
	i, err := c.Del("k", "k3").Result()
	assert.Equal(t, i, int64(2))
	assert.NoError(t, err)

	i, err = c.Del("abc").Result()
	assert.Zero(t, i)
	assert.NoError(t, err)
}

func TestTtlCommand(t *testing.T) {
	s, err := c.Set("aKey", "hey", 1*time.Minute).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
	s, err = c.Set("bKey", "hallo", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	ttl, err := c.TTL("aKey").Result()
	assert.True(t, ttl.Seconds() > 55 && ttl.Seconds() < 61, "ttl: %d", ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL("none").Result()
	assert.Equal(t, time.Duration(-2000000000), ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL("bKey").Result()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(-1000000000), ttl)
}

func TestExpiry(t *testing.T) {
	s, err := c.Set("x", "val", 10*time.Millisecond).Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", s)

	time.Sleep(10 * time.Millisecond)

	s, err = c.Get("x").Result()
	assert.Equal(t, "", s)
	assert.Error(t, err)
}

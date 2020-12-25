package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLPushCommand(t *testing.T) {
	i, err := c.LPush("lpushkey", "va").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), i)

	i, err = c.LPush("lpushkey", "vb").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	i, err = c.LPush("lpushkey", "vc", "vd").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(4), i)

	i, err = c.LPush("lpushkey2", "1", "2").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	i, err = c.LPush("lpush3key").Result()
	assert.Error(t, err)
}

func TestLPopCommand(t *testing.T) {
	s, err := c.LPop("lpop1").Result()
	assert.Zero(t, s)
	assert.Error(t, err)

	i, err := c.LPush("list", "a", "b").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	s, err = c.LPop("list").Result()
	assert.NoError(t, err)
	assert.Equal(t, "b", s)

	s, err = c.LPop("list").Result()
	assert.NoError(t, err)
	assert.Equal(t, "a", s)

	s, err = c.LPop("list").Result()
	assert.Error(t, err)
	assert.Zero(t, s)
}

func TestLRangeCommand(t *testing.T) {
	s, err := c.LRange("lrange", 0, 0).Result()
	assert.Error(t, err)
	assert.Zero(t, s)

	sl, err := c.Set("works", "esfkjsefj", 0).Result()
	assert.NoError(t, err)
	assert.NotZero(t, sl)
	assert.NotEmpty(t, sl)

	i, err := c.LPush("list2", "a", "b").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)
}

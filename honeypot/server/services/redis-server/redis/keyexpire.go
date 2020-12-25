package redis

import (
	"math"
	"time"
)

type KeyExpirer interface {
	Start(tick time.Duration, keyNum int, againPercentage int)
	Stop()
}

var _ KeyExpirer = (*Expirer)(nil)

type Expirer struct {
	redis *Redis

	done chan bool
}

func NewKeyExpirer(r *Redis) *Expirer {
	return &Expirer{
		redis: r,
		done:  make(chan bool, math.MaxInt32),
	}
}

// Start starts the Expirer.
//
// tick - How fast is the cleaner triggered.
//
// randomKeys - Amount of random expiring keys to get checked.
//
// againPercentage - If more than x% of keys were expired, start again in same tick.
func (e *Expirer) Start(tick time.Duration, randomKeys int, againPercentage int) {
	ticker := time.NewTicker(tick)
	for {
		select {
		case <-ticker.C:
			e.do(randomKeys, againPercentage)
		case <-e.done:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the
func (e *Expirer) Stop() {
	if e.done != nil {
		e.done <- true
		close(e.done)
	}
}

func (e *Expirer) do(randomKeys, againPercentage int) {
	var deletedKeys int

	dbs := make(map[*RedisDb]struct{})
	for _, db := range e.Redis().RedisDbs() {
		if !db.HasExpiringKeys() {
			continue
		}
		dbs[db] = struct{}{}
	}

	if len(dbs) == 0 {
		return
	}

	for c := 0; c < randomKeys; c++ {
		// get random db
		db := func() *RedisDb {
			for db := range dbs {
				return db
			}
			return nil // won't happen
		}()

		// get random key
		k := func() *string {
			for k := range db.ExpiringKeys() {
				return &k
			}
			return nil
		}()

		if k == nil {
			continue
		}

		// del if expired
		if db.DeleteExpired(k) != 0 {
			deletedKeys++
		}
	}

	// Start again in new goroutine so keys are deleted fast
	if againPercentage > 0 && deletedKeys/randomKeys*100 > againPercentage {
		go e.do(randomKeys, againPercentage)
	}
}

// Redis gets the redis instance.
func (e *Expirer) Redis() *Redis {
	return e.redis
}

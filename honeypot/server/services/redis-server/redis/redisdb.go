package redis

import (
	"sync"
	"time"
)

const (
	keysMapSize           = 32
	redisDbMapSizeDefault = 3
)

// A redis database.
// There can be more than one in a redis instance.
type RedisDb struct {
	// Database id
	id DatabaseId

	// All keys in this db.
	keys Keys

	// Keys with expire timestamp.
	expiringKeys ExpiringKeys

	// TODO long long avg_ttl;          /* Average TTL, just for stats */

	redis *Redis
}

// Redis databases map
type RedisDbs map[DatabaseId]*RedisDb

// Database id
type DatabaseId uint

// Key-Item map
type Keys map[string]Item

// Keys with expire timestamp.
type ExpiringKeys map[string]time.Time

// The item interface. An item is the value of a key.
type Item interface {
	// The pointer to the value.
	Value() interface{}

	// The id of the type of the Item.
	// This need to be constant for the type because it is
	// used when de-/serializing item from/to disk.
	Type() uint64
	// The type of the Item as readable string.
	TypeFancy() string

	// OnDelete is triggered before the key of the item is deleted.
	// db is the affected database.
	OnDelete(key *string, db *RedisDb)
}

// NewRedisDb creates a new db.
func NewRedisDb(id DatabaseId, r *Redis) *RedisDb {
	return &RedisDb{
		id:           id,
		redis:        r,
		keys:         make(Keys, keysMapSize),
		expiringKeys: make(ExpiringKeys, keysMapSize),
	}
}

// RedisDb gets the redis database by its id or creates and returns it if not exists.
func (r *Redis) RedisDb(dbId DatabaseId) *RedisDb {
	getDb := func() *RedisDb { // returns nil if db not exists
		if db, ok := r.redisDbs[dbId]; ok {
			return db
		}
		return nil
	}

	r.Mu().RLock()
	db := getDb()
	r.Mu().RUnlock()
	if db != nil {
		return db
	}

	// create db
	r.Mu().Lock()
	defer r.Mu().Unlock()
	// check if db does not exists again since
	// multiple "mutex readers" can come to this point
	db = getDb()
	if db != nil {
		return db
	}
	// now really create db of that id
	r.redisDbs[dbId] = NewRedisDb(dbId, r)
	return r.redisDbs[dbId]
}

// RedisDbs gets all redis databases.
func (r *Redis) RedisDbs() RedisDbs {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.redisDbs
}

// Redis gets the redis instance.
func (db *RedisDb) Redis() *Redis {
	return db.redis
}

// Mu gets the mutex.
func (db *RedisDb) Mu() *sync.RWMutex {
	return db.Redis().Mu()
}

// Id gets the db id.
func (db *RedisDb) Id() DatabaseId {
	return db.id
}

// Sets a key with an item which can have an expiration time.
func (db *RedisDb) Set(key *string, i Item, expires bool, expiry time.Time) {
	db.Mu().Lock()
	defer db.Mu().Unlock()
	db.keys[*key] = i
	if expires {
		db.expiringKeys[*key] = expiry
	}
}

// Returns the item by the key or nil if key does not exists.
func (db *RedisDb) Get(key *string) Item {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.get(key)
}

func (db *RedisDb) get(key *string) Item {
	i, _ := db.keys[*key]
	return i
}

// Deletes a key, returns number of deleted keys.
func (db *RedisDb) Delete(keys ...*string) int {
	db.Mu().Lock()
	defer db.Mu().Unlock()
	return db.delete(keys...)
}

// If checkExists is false, then return bool is reprehensible.
func (db *RedisDb) delete(keys ...*string) int {
	do := func(k *string) bool {
		if k == nil {
			return false
		}
		i := db.get(k)
		if i == nil {
			return false
		}
		i.OnDelete(k, db)
		delete(db.keys, *k)
		delete(db.expiringKeys, *k)
		return true
	}

	var c int
	for _, k := range keys {
		if do(k) {
			c++
		}
	}

	return c
}

func (db *RedisDb) DeleteExpired(keys ...*string) int {
	var c int
	for _, k := range keys {
		if k != nil && db.Expired(k) && db.Delete(k) > 0 {
			c++
		}
	}
	return c
}

// GetOrExpire gets the item or nil if expired or not exists. If 'deleteIfExpired' is true the key will be deleted.
func (db *RedisDb) GetOrExpire(key *string, deleteIfExpired bool) Item {
	// TODO mutex optimize this func so that a RLock is mainly first opened
	db.Mu().Lock()
	defer db.Mu().Unlock()
	i, ok := db.keys[*key]
	if !ok {
		return nil
	}
	if db.expired(key) {
		if deleteIfExpired {
			db.delete(key)
		}
		return nil
	}
	return i
}

// IsEmpty checks if db is empty.
func (db *RedisDb) IsEmpty() bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return len(db.keys) == 0
}

// HasExpiringKeys checks if db has any expiring keys.
func (db *RedisDb) HasExpiringKeys() bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return len(db.expiringKeys) != 0
}

// Check if key exists.
func (db *RedisDb) Exists(key *string) bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.exists(key)
}
func (db *RedisDb) exists(key *string) bool {
	_, ok := db.keys[*key]
	return ok
}

// Check if key has an expiry set.
func (db *RedisDb) Expires(key *string) bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expires(key)
}
func (db *RedisDb) expires(key *string) bool {
	_, ok := db.expiringKeys[*key]
	return ok
}

// Expired only check if a key can and is expired.
func (db *RedisDb) Expired(key *string) bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expired(key)
}
func (db *RedisDb) expired(key *string) bool {
	return db.expires(key) && TimeExpired(db.expiry(key))
}

// Expiry gets the expiry of the key has one.
func (db *RedisDb) Expiry(key *string) time.Time {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expiry(key)
}

func (db *RedisDb) expiry(key *string) time.Time {
	return db.expiringKeys[*key]
}

// Keys gets all keys in this db.
func (db *RedisDb) Keys() Keys {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.keys
}

// ExpiringKeys gets keys with an expiry set and their timeout.
func (db *RedisDb) ExpiringKeys() ExpiringKeys {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expiringKeys
}

// TimeExpired check if a timestamp is older than now.
func TimeExpired(expireAt time.Time) bool {
	return time.Now().After(expireAt)
}

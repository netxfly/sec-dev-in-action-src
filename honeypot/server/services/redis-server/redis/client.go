package redis

import (
	"github.com/redis-go/redcon"
	"sync"
)

// TODO Client flags
const (
// client is a master
// client is a slave
// ...
)

// A connected Client.
type Client struct {
	clientId ClientId
	// The client connection.
	conn redcon.Conn

	// Selected database (default 0)
	db DatabaseId

	redis *Redis
	flag  bool
}

// NewClient creates new client and adds it to the redis.
func (r *Redis) NewClient(conn redcon.Conn) *Client {
	c := &Client{
		conn:     conn,
		redis:    r,
		clientId: r.NextClientId(),
	}
	return c
}

// NextClientId atomically gets and increments a counter to return the next client id.
func (r *Redis) NextClientId() ClientId {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	id := r.nextClientId
	r.nextClientId++
	return id
}

// Clients gets the current connected clients.
func (r *Redis) Clients() Clients {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.clients
}

func (r *Redis) getClients() Clients {
	return r.clients
}

// Redis gets the redis instance.
func (c *Client) Redis() *Redis {
	return c.redis
}

// Mu the mutex.
func (c *Client) Mu() *sync.RWMutex {
	return c.Redis().Mu()
}

// ClientId get the client id.
func (c *Client) ClientId() ClientId {
	return c.clientId
}

// The client's connection.
func (c *Client) Conn() redcon.Conn {
	c.Mu().RLock()
	defer c.Mu().RUnlock()
	return c.conn
}

// SelectDb selects the clients database.
func (c *Client) SelectDb(db DatabaseId) {
	c.Mu().Lock()
	defer c.Mu().Unlock()
	c.db = db
}

// DbId gets the clients selected database id.
func (c *Client) DbId() DatabaseId {
	c.Mu().RLock()
	defer c.Mu().RUnlock()
	return c.db
}

// Db gets the clients selected database.
func (c *Client) Db() *RedisDb {
	return c.Redis().RedisDb(c.DbId())
}

// Disconnects and removes a Client.
func (c *Client) FreeClient() {
	c.Conn().Close() // TODO should we log on error?
	c.Mu().Lock()
	defer c.Mu().Unlock()
	delete(c.Redis().getClients(), c.ClientId())
}

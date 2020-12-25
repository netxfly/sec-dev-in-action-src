package redis

import (
	"strings"
	"sync"

	"github.com/redis-go/redcon"
)

const (
	SyntaxErr         = "ERR syntax error"
	InvalidIntErr     = "ERR value is not an integer or out of range"
	WrongTypeErr      = "WRONGTYPE Operation against a key holding the wrong kind of value"
	WrongNumOfArgsErr = "ERR wrong number of arguments for '%s' command"
)

// This is the redis server.
type Redis struct {
	// databases/keyspaces
	redisDbs RedisDbs

	// Locking is important, share this mutex around to provide state.
	mu *sync.RWMutex

	commands       Commands
	unknownCommand UnknownCommand

	handler Handler

	accept  Accept
	onClose OnClose

	// TODO version
	// TODO log writer
	// TODO modules
	// TODO redis options type

	keyExpirer KeyExpirer

	clients      Clients
	nextClientId ClientId
}

// A Handler is called when a request is received and after Accept
// (if Accept allowed the connection by returning true).
//
// For implementing an own handler see the default handler
// as a perfect example in the createDefault() function.
type Handler func(c *Client, cmd redcon.Command)

// Accept is called when a Client tries to connect and before everything else,
// the Client connection will be closed instantaneously if the function returns false.
type Accept func(c *Client) bool

// OnClose is called when a Client connection is closed.
type OnClose func(c *Client, err error)

// Client map
type Clients map[ClientId]*Client

// Client id
type ClientId uint64

// Gets the handler func.
func (r *Redis) HandlerFn() Handler {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.handler
}

// Sets the handler func.
// Live updates (while redis is running) works.
func (r *Redis) SetHandlerFn(new Handler) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.handler = new
}

// Gets the accept func.
func (r *Redis) AcceptFn() Accept {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.accept
}

// Sets the accept func.
// Live updates (while redis is running) works.
func (r *Redis) SetAcceptFn(new Accept) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.accept = new
}

// Gets the onclose func.
func (r *Redis) OnCloseFn() OnClose {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.onClose
}

// Sets the onclose func.
// Live updates (while redis is running) works.
func (r *Redis) SetOnCloseFn(new OnClose) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.onClose = new
}

// The mutex of the redis.
func (r *Redis) Mu() *sync.RWMutex {
	return r.mu
}

func (r *Redis) KeyExpirer() KeyExpirer {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.keyExpirer
}

func (r *Redis) SetKeyExpirer(ke KeyExpirer) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.keyExpirer = ke
}

var defaultRedis *Redis

// Default redis server.
// Initializes the default redis if not already.
// You can change the fields or value behind the pointer
// of the returned redis pointer to extend/change the default.
func Default() *Redis {
	if defaultRedis != nil {
		return defaultRedis
	}
	defaultRedis = createDefault()
	return defaultRedis
}

// createDefault creates a new default redis.
func createDefault() *Redis {
	// initialize default redis server
	r := &Redis{
		mu: new(sync.RWMutex),
		accept: func(c *Client) bool {
			return true
		},
		onClose: func(c *Client, err error) {
		},
		handler: func(c *Client, cmd redcon.Command) {
			cmdl := strings.ToLower(string(cmd.Args[0]))
			if c.Redis().CommandExists(cmdl) {
				c.Redis().CommandHandlerFn(cmdl)(c, cmd)
			} else {
				c.Redis().UnknownCommandFn()(c, cmd)
			}
		},
		unknownCommand: func(c *Client, cmd redcon.Command) {
			c.Conn().WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
		},
		commands: make(Commands, 0),
	}
	r.redisDbs = make(RedisDbs, redisDbMapSizeDefault)
	r.RedisDb(0) // initializes default db 0
	r.keyExpirer = KeyExpirer(NewKeyExpirer(r))

	r.RegisterCommands([]*Command{
		NewCommand("ping", PingCommand, CMD_STALE, CMD_FAST),
		NewCommand("set", SetCommand, CMD_WRITE, CMD_DENYOOM),
		NewCommand("get", GetCommand, CMD_READONLY, CMD_FAST),
		NewCommand("keys", KeysCommand, CMD_READONLY, CMD_FAST),
		NewCommand("del", DelCommand, CMD_WRITE),
		NewCommand("ttl", TtlCommand, CMD_READONLY, CMD_FAST),

		NewCommand("lpush", LPushCommand, CMD_WRITE, CMD_FAST, CMD_DENYOOM),
		NewCommand("rpush", RPushCommand, CMD_WRITE, CMD_FAST, CMD_DENYOOM),
		NewCommand("lpop", LPopCommand, CMD_WRITE, CMD_FAST),
		NewCommand("rpop", RPopCommand, CMD_WRITE, CMD_FAST),
		NewCommand("lrange", LRangeCommand, CMD_READONLY),
	})
	return r
}

package redis

import (
	"fmt"
	"github.com/redis-go/redcon"
	"time"
)

func TtlCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments (given %d, expected 1)", len(cmd.Args)-1))
		return
	}

	db := c.Db()
	key := string(cmd.Args[1])
	db.DeleteExpired(&key)
	if !db.Exists(&key) {
		c.Conn().WriteInt(-2)
		return
	}

	t := db.Expiry(&key)
	if t.IsZero() {
		c.Conn().WriteInt(-1)
		return
	}

	c.Conn().WriteInt64(int64(t.Sub(time.Now()).Seconds()))
}

package redis

import (
	"github.com/redis-go/redcon"
)

func DelCommand(c *Client, cmd redcon.Command) {
	db := c.Db()
	keys := make([]*string, 0, len(cmd.Args)-1)
	for i := 1; i < len(cmd.Args); i++ {
		k := string(cmd.Args[i])
		keys = append(keys, &k)
	}
	dels := db.Delete(keys...)
	c.Conn().WriteInt(dels)
}

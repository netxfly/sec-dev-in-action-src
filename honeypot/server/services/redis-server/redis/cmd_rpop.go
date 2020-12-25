package redis

import (
	"fmt"
	"github.com/redis-go/redcon"
)

func RPopCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "rpop"))
		return
	}
	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteNull()
		return
	} else if i.Type() != ListType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), ListTypeFancy))
		return
	}

	l := i.(*List)
	c.Redis().Mu().Lock()
	v, b := l.RPop()
	if b {
		db.Delete(&key)
	}
	c.Redis().Mu().Unlock()

	c.Conn().WriteBulkString(*v)
}

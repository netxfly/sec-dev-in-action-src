package redis

import (
	"fmt"
	"strconv"

	"github.com/redis-go/redcon"
)

func LRangeCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "lrange"))
		return
	}
	key := string(cmd.Args[1])
	start, err := strconv.Atoi(string(cmd.Args[2]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
		return
	}
	end, err := strconv.Atoi(string(cmd.Args[3]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
		return
	}

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
	c.Redis().Mu().RLock()
	values := l.LRange(start, end)
	c.Redis().Mu().RUnlock()

	c.Conn().WriteArray(len(values))
	for _, v := range values {
		c.Conn().WriteBulkString(v)
	}
}

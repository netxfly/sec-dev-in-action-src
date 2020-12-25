package redis

import (
	"fmt"
	"github.com/redis-go/redcon"
)

func GetCommand(c *Client, cmd redcon.Command) {
	key := string(cmd.Args[1])

	i := c.Db().GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteNull()
		return
	}

	if i.Type() != StringType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), StringTypeFancy))
		return
	}

	v := *i.Value().(*string)
	c.Conn().WriteBulkString(v)
}

func KeysCommand(c *Client, cmd redcon.Command) {
	key := string(cmd.Args[1])
	fmt.Println(key)
	i := c.Db().Keys()
	if i == nil {
		c.Conn().WriteNull()
		return
	}

	c.Conn().WriteArray(len(i))
	for k, _ := range i {
		c.Conn().WriteBulkString(k)
	}

}

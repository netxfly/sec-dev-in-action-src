package redis

import (
	"bytes"
	"github.com/redis-go/redcon"
)

func PingCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) > 1 {
		var buf bytes.Buffer
		for i := 1; i < len(cmd.Args); i++ {
			buf.Write(cmd.Args[i])
			buf.WriteString(" ")
		}
		s := buf.String()
		s = s[:len(s)-1]
		c.Conn().WriteString(s)
		return
	}
	c.Conn().WriteString("PONG")
}

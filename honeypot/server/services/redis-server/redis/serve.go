package redis

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/redis-go/redcon"

	"sec-dev-in-action-src/honeypot/server/logger"
	"sec-dev-in-action-src/honeypot/server/pusher"
	"sec-dev-in-action-src/honeypot/server/util"
)

var log = logging.MustGetLogger("redis-service")

// Run runs the default redis server.
// Initializes the default redis if not already.
func Run(addr string, flag bool) error {
	return Default().Run(addr, flag)
}

// Run runs the redis server.
func (r *Redis) Run(addr string, flag bool) error {
	go r.KeyExpirer().Start(100*time.Millisecond, 20, 25)
	return redcon.ListenAndServe(
		addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			if flag {
				rawIp, ProxyAddr, timeStamp := util.GetRawIpByConn(conn.NetConn())
				tmpCmd := make([]string, 0)
				for _, c := range cmd.Args {
					tmpCmd = append(tmpCmd, string(c))
				}

				var message pusher.HoneypotMessage
				message.Timestamp = timeStamp
				message.RawIp = rawIp
				message.ProxyAddr = ProxyAddr.String()
				message.ServiceType = "redis"

				data := make(map[string]interface{})
				data["cmd"] = strings.Join(tmpCmd, " ")
				message.Data = data

				strMessage, _ := message.Build()
				logger.Log.Info(strMessage)
				_ = message.Send()
			}

			r.HandlerFn()(r.NewClient(conn), cmd)
		},
		func(conn redcon.Conn) bool {
			return r.AcceptFn()(r.NewClient(conn))
		},
		func(conn redcon.Conn, err error) {
			r.OnCloseFn()(r.NewClient(conn), err)
		},
	)
}

// Run runs the redis server with tls.
func (r *Redis) RunTLS(addr string, tls *tls.Config) error {
	go r.KeyExpirer().Start(100*time.Millisecond, 20, 25)
	return redcon.ListenAndServeTLS(
		addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			r.HandlerFn()(r.NewClient(conn), cmd)
		},
		func(conn redcon.Conn) bool {
			return r.AcceptFn()(r.NewClient(conn))
		},
		func(conn redcon.Conn, err error) {
			r.OnCloseFn()(r.NewClient(conn), err)
		},
		tls,
	)
}

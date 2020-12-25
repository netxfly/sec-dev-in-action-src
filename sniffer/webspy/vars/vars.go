package vars

import "sync"

var (
	HttpHost = "127.0.0.1"
	HttpPort = 4000

	Data = sync.Pool{}
)

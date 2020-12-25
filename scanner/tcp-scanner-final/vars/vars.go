package vars

import "sync"

var (
	ThreadNum = 5000
	Result    *sync.Map

	Host    string
	Port    = "22,23,53,80-139"
	Mode    = "syn"
	Timeout = 2
)

func init() {
	Result = &sync.Map{}
}

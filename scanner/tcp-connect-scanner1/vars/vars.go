package vars

import "sync"

var (
	ThreadNum = 5000
	Result    *sync.Map
)

func init() {
	Result = &sync.Map{}
}

package models

import (
	"fmt"
	"sync"
)

type ProxyInfo struct {
	Addr     string
	Port     int
	Protocol string
}

var (
	Result sync.Map
)

func SaveProxies(isProxy bool, proxyInfo ProxyInfo, err error) error {
	if err == nil && isProxy {
		k := fmt.Sprintf("%v://%v:%v", proxyInfo.Protocol, proxyInfo.Addr, proxyInfo.Port)
		Result.Store(k, true)
	}

	return err
}

func PrintResult() {
	Result.Range(func(key, value interface{}) bool {
		fmt.Printf("%v\n", key)
		return true
	})
}

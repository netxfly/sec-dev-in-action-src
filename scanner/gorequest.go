package main

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

func main() {
	request := gorequest.New()
	resp, body, err := request.Proxy("http://sec.lu:8080").Get("http://mail.163.com").End()
	fmt.Printf("resp: %v, body: %v, err: %v\n", resp, body, err)
}

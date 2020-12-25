/*

Copyright (c) 2018 sec.lu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	Logger = log.New(os.Stderr, "", 0)
)

func usedTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		Logger.Printf("I am %v middleware", "usedTime")
		next.ServeHTTP(w, r)
		usedTime := time.Since(timeStart)
		Logger.Printf("used time: %v, I am %v middleware", usedTime, "usedTime")
	})
}
func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Printf("I am %v middleware", "Authentication")
		next.ServeHTTP(w, r)
	})
}

func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Printf("I am %v middleware", "Authorization")
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		Logger.Printf("I am %v middleware", "logging")
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func blog(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("blog"))
}

func admin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("admin"))
}

func user(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/blog/", blog)
	r.HandleFunc("/admin/", index)
	r.HandleFunc("/user/", index)

	r.Use(usedTimeMiddleware)
	r.Use(authenticationMiddleware)
	r.Use(authorizationMiddleware)

	err := http.ListenAndServe(":80", r)
	_ = err
}

package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"rlimit/rserver/config"
)

func newLimiter() *rate.Limiter {
	cfLoad := config.Load()
	return rate.NewLimiter(rate.Limit(cfLoad.Limit), cfLoad.Burst)
}

func limit(next http.Handler) http.Handler {
	limiter := newLimiter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/limiter", okHandler)    // 假设okHandler是一个处理函数
	http.ListenAndServe(":4000", limit(mux)) // 应用限流中间件
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

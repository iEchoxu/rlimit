package main

import (
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
	"log"
	"net/http"
)

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi there!"))
})

func main1() {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		log.Fatal(err)
	}

	// Maximum burst of 5 which refills at 20 tokens per minute.
	quota := throttled.RateQuota{MaxRate: throttled.PerMin(20), MaxBurst: 5}

	rateLimiter, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiterCtx{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	http.ListenAndServe(":8080", httpRateLimiter.RateLimit(myHandler))
}

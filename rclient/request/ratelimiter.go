package request

import (
	"context"
	"errors"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
	"log"
)

var globalRateLimiter = &rateLimiter{rateMap: make(map[string]*throttled.GCRARateLimiterCtx)}

type rateLimiter struct {
	rateMap map[string]*throttled.GCRARateLimiterCtx
}

func Register(limiter string, limitPerSec int) error {
	gcraRateLimiter, err := newRateLimiter(limitPerSec)
	if err != nil {
		return err
	}

	globalRateLimiter.rateMap[limiter] = gcraRateLimiter

	return nil
}

func newRateLimiter(limitPerSec int) (*throttled.GCRARateLimiterCtx, error) {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// MaxBurst 表示瞬时爆发的请求，如果设置得大于 0 表示最多可发送 MaxBurst + limitPerSec 个请求
	// MaxBurst 设置是关键，即使 limitPerSec 设置的比服务器的 limit 大，其每秒最多也只能发送服务器设置的 limit 个请求
	// 当 MaxBurst 设置为 0 时， limitPerSec 设置为 25 ，服务器设置的 limit 为 20 且允许 20 个爆发请求，此时每秒发送的请求个数为 20
	rateQuota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(limitPerSec), // 每秒的限额,可修改 throttled.PerHour() 等
		MaxBurst: 0,                             // 设置为 0 时永远会触发限速，同时限制每秒可发送 limitPerSec 个请求
	}

	gcraRateLimiter, gcraRateLimiterError := throttled.NewGCRARateLimiterCtx(store, rateQuota)
	if gcraRateLimiterError != nil {
		return nil, errors.New(gcraRateLimiterError.Error())
	}

	return gcraRateLimiter, nil
}

func GetRateLimiter(key string) (bool, *throttled.RateLimitResult, error) {
	rater, ok := globalRateLimiter.rateMap[key]
	if !ok {
		return false, nil, errors.New("CheckRateLimit-KeyError")
	}

	limited, result, err := rater.RateLimitCtx(context.Background(), key, 1) // quantity 为 0 时不能打印 result 等信息
	if err != nil {
		return false, nil, errors.New("CheckRateLimit-RateError")
	}

	return limited, &result, nil
}

func GetRateMapKey() {
	for k := range globalRateLimiter.rateMap {
		log.Println("已注册 ", k, " 限流器")
	}
}

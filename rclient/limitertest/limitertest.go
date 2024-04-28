package limitertest

import (
	"log"
	"rlimit/rclient/request"
)

const (
	LocalLimiterType RunMode = iota + 1
	ThrottledLimiterType
)

var limiterTypeMap = map[RunMode]func(LimiterParameter) Limiter{
	LocalLimiterType:     NewLocalLimiter,
	ThrottledLimiterType: NewThrottledLimiter,
}

// 根据配置项中的 httpClient 来决定是否传递 rateLimiter 参数来初始化 requestLimiter 接口实例
var requestLimiterOptionMap = make(map[request.HttpClientType]request.Request, 2)

type RunMode int

type Limiter interface {
	Run()
}

func New(runMode RunMode, limiterOption ...LimiterParaOption) Limiter {
	opts := new(LimiterParameter)
	for _, option := range limiterOption {
		option(opts)
	}
	return limiterTypeMap[runMode](*opts)
}

func registerThrottledRateLimiter(rateLimiter string, limiterSize int) {
	errRater := request.Register(rateLimiter, limiterSize)
	if errRater != nil {
		log.Println(errRater)
		return
	}

	request.GetRateMapKey()
}

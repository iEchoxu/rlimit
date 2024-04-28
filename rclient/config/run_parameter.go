package config

import (
	"fmt"
	"rlimit/rclient/limitertest"
	"rlimit/rclient/request"
)

// 注意: 此文件是 Main 函数运行时需要的参数与函数
// 注意: 此文件中的配置项为不经常改动的项

// 接口名称映射表
const (
	getRateLimiter     string = "limiter"
	getUserInfoLimiter string = "userInfo"
)

var interfaceName = getRateLimiter
var url = func() string {
	base := "http://127.0.0.1:4000/"
	return fmt.Sprintf("%s%s", base, interfaceName)
}()
var limiterSize = limitPerSecMap[interfaceName]
var totalSize = 55

var limiterOptionMap = map[request.HttpClientType][]limitertest.LimiterParaOption{
	request.DefaultHttpClient:   createCommonLimiterOptions(),
	request.RetryableHttpClient: append(createCommonLimiterOptions(), limitertest.WithRateLimiter(interfaceName)),
}

func createCommonLimiterOptions() []limitertest.LimiterParaOption {
	return []limitertest.LimiterParaOption{
		limitertest.WithUrlOption(url),
		limitertest.WithLimiterSize(limiterSize),
		limitertest.WithTotalSize(totalSize),
		limitertest.WithHttpClient(HttpClient),
		limitertest.WithLimiterType(limiterType),
	}
}

func NewRunner(runMode limitertest.RunMode, httpClient request.HttpClientType) limitertest.Limiter {
	return limitertest.New(runMode, limiterOptionMap[httpClient]...)
}

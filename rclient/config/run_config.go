package config

import (
	"rlimit/rclient/limitertest"
	"rlimit/rclient/request"
)

// 注意: 此文件包含的是需要经常修改的配置项

// RunMode  LocalLimiterType |  ThrottledLimiterType
var RunMode = limitertest.ThrottledLimiterType

var HttpClient = request.RetryableHttpClient

var limiterType = request.Limit

// limitPerSecMap 接口限流映射表，根据每个接口的限流情况不同修改为实际的值
var limitPerSecMap = map[string]int{
	getRateLimiter:     10,
	getUserInfoLimiter: 3,
}

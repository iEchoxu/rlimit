package request

import (
	"io"
)

const (
	Limit LimiterType = iota + 1
	LocalLimit
)

var requestFactories = map[LimiterType]func(option *RateLimiterRequest) Request{
	LocalLimit: LocalRateLimitRequestFactory,
	Limit:      RateLimitRequestFactory,
}

type LimiterType int

type Request interface {
	Get(url string) ([]byte, error)
	Post(url string, body io.Reader) ([]byte, error)
}

type LimiterOption func(*RateLimiterRequest)

type RateLimiterRequest struct {
	RateLimiterName string
}

func WithRateLimiterName(l string) LimiterOption {
	return func(r *RateLimiterRequest) {
		r.RateLimiterName = l
	}
}

// NewLimiter 构造接口类型
// 不要修改此代码，在 requestFactories 中添加新类型
func NewLimiter(limiterType LimiterType, limitOption ...LimiterOption) Request {
	opts := new(RateLimiterRequest)
	for _, option := range limitOption {
		option(opts)
	}

	// 如果没有找到对应的工厂函数，使用RateLimitRequest
	if _, exists := requestFactories[limiterType]; !exists {
		return requestFactories[Limit](opts)
	}

	return requestFactories[limiterType](opts)
}

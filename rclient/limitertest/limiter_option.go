package limitertest

import (
	"rlimit/rclient/request"
)

// LimiterParameter 用于往下传参数
type LimiterParameter struct {
	Url         string
	LimiterSize int
	TotalSize   int
	RateLimiter string
	HttpClient  request.HttpClientType
	LimiterType request.LimiterType
}

type LimiterParaOption func(option *LimiterParameter)

func WithUrlOption(url string) LimiterParaOption {
	return func(option *LimiterParameter) {
		option.Url = url
	}
}

func WithLimiterSize(limiterSize int) LimiterParaOption {
	return func(option *LimiterParameter) {
		option.LimiterSize = limiterSize
	}
}

func WithTotalSize(totalSize int) LimiterParaOption {
	return func(option *LimiterParameter) {
		option.TotalSize = totalSize
	}
}

func WithRateLimiter(rateLimiter string) LimiterParaOption {
	return func(option *LimiterParameter) {
		option.RateLimiter = rateLimiter
	}
}

func WithHttpClient(httpClient request.HttpClientType) LimiterParaOption {
	return func(option *LimiterParameter) {
		option.HttpClient = httpClient
	}
}

func WithLimiterType(limiterType request.LimiterType) LimiterParaOption {
	return func(option *LimiterParameter) {
		option.LimiterType = limiterType
	}
}

func NewLimiterOption(options ...LimiterParaOption) *LimiterParameter {
	opts := new(LimiterParameter)

	for _, option := range options {
		option(opts)
	}

	return opts
}

func (lo *LimiterParameter) buildURLS(url string, count int) []string {
	var urls []string
	for i := 0; i < count; i++ {
		urls = append(urls, url)
	}
	return urls
}

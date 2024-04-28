package request

import (
	"io"
)

type RateLimitRequest struct {
	*baseHttpRequest
}

// RateLimitRequestFactory 创建一个有速率限制的请求对象
func RateLimitRequestFactory(rateLimiterOption *RateLimiterRequest) Request {
	return &RateLimitRequest{
		baseHttpRequest: &baseHttpRequest{
			baseRateLimiterName: rateLimiterOption.RateLimiterName,
		},
	}
}

func (rl *RateLimitRequest) Get(url string) ([]byte, error) {
	return rl.makeRequest("GET", url, nil)
}

func (rl *RateLimitRequest) Post(url string, body io.Reader) ([]byte, error) {
	return rl.makeRequest("POST", url, body)
}

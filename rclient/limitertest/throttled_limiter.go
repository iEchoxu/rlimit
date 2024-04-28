package limitertest

import (
	"fmt"
	"log"
	"rlimit/rclient/request"
	"sync/atomic"
)

type ThrottledLimiter struct {
	LimiterParameter
}

func NewThrottledLimiter(limiterOption LimiterParameter) Limiter {
	return &ThrottledLimiter{
		LimiterParameter: limiterOption,
	}
}

func (t *ThrottledLimiter) Run() {
	urls := t.buildURLS(t.Url, t.TotalSize)
	successCount := int32(0)
	failCount := int32(0)

	// 只有使用 RetryableHttpClient 时才需要注册 GetRateLimiter
	if t.HttpClient == request.RetryableHttpClient {
		registerThrottledRateLimiter(t.RateLimiter, t.LimiterSize)
		requestLimiterOptionMap[t.HttpClient] = request.NewLimiter(
			request.Limit,
			request.WithRateLimiterName(t.RateLimiter),
		)
		fmt.Println("打印信息: ", t.Url, t.LimiterSize, t.TotalSize, t.RateLimiter, t.HttpClient, t.LimiterType)
	}

	if t.HttpClient == request.DefaultHttpClient {
		requestLimiterOptionMap[t.HttpClient] = request.NewLimiter(request.Limit)
		fmt.Println("打印信息: ", t.Url, t.LimiterSize, t.TotalSize, t.HttpClient, t.LimiterType)
	}

	request.SetClient(t.HttpClient)

	// 如果任务数量较多可以结合 sync.waitGroup 和 channel 来提高并发能力
	for _, url := range urls {
		get, err := requestLimiterOptionMap[t.HttpClient].Get(url)
		if err != nil {
			log.Println(err)
			atomic.AddInt32(&failCount, 1)
			// 可以将失败的请求传入一个 failChan 然后做失败后的任务处理
			continue
		}

		log.Println(string(get))
		atomic.AddInt32(&successCount, 1)
	}

	log.Printf("总任务数: %d\n", len(urls))
	log.Printf("成功的任务数: %d\n", atomic.LoadInt32(&successCount))
	log.Printf("失败的任务数: %d\n", atomic.LoadInt32(&failCount))
}

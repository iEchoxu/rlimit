package limitertest

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"rlimit/rclient/request"
	"sync"
	"sync/atomic"
	"time"
)

type LocalLimiter struct {
	LimiterParameter
}

func NewLocalLimiter(l LimiterParameter) Limiter {
	return &LocalLimiter{
		LimiterParameter: l,
	}
}

// Run 使用 rate.NewLimiter 结合协程并发发送多个请求
// 这段代码有可能会触发 429 错误，因为如果请求完成的非常快，会导致在非常短的时间内启动多个协程，造成瞬时请求大于 10
// 要确保每一秒内的请求均匀的分布
func (ll *LocalLimiter) Run() {
	//ll.batchSize = ll.batchSize-1 // 限制每秒请求数小于服务器设置的每秒请求时,不然会报错 429
	urls := ll.buildURLS(ll.Url, ll.TotalSize)

	limit := make(chan struct{}, ll.LimiterSize)
	var wg sync.WaitGroup

	successCount := int32(0)
	failCount := int32(0)

	rateLimiter := rate.NewLimiter(rate.Limit(ll.LimiterSize), ll.LimiterSize) // 每秒最大爆发请求最好设置得比服务器设置稍小些
	timeout := time.After(60 * time.Second)                                    // 用于控制 limit channel 等待时间

	// 仅当 httpClient 为 RetryableHttpClient 时才注册 ThrottledRateLimiter
	// 初始化执行函数并将 RateLimiter 传过去,因为 ThrottledRateLimiter 需要获取这个值
	if ll.HttpClient == request.RetryableHttpClient {
		registerThrottledRateLimiter(ll.RateLimiter, ll.LimiterSize)
		requestLimiterOptionMap[ll.HttpClient] = request.NewLimiter(
			request.Limit,
			request.WithRateLimiterName(ll.RateLimiter),
		)
		fmt.Println("打印信息: ", ll.Url, ll.LimiterSize, ll.TotalSize, ll.RateLimiter, ll.HttpClient, ll.LimiterType)
	}

	// 当使用 DefaultHttpClient 时不需要传递 RateLimiter 参数
	if ll.HttpClient == request.DefaultHttpClient {
		requestLimiterOptionMap[ll.HttpClient] = request.NewLimiter(request.Limit)
		fmt.Println("打印信息: ", ll.Url, ll.LimiterSize, ll.TotalSize, ll.HttpClient, ll.LimiterType)
	}

	// 注册 httpClient
	request.SetClient(ll.HttpClient)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer func() {
				<-time.After(time.Second)  // 等待 1s 后释放信号量
				<-limit
				wg.Done()
			}()

			select {
			case limit <- struct{}{}:
			case <-timeout:
				log.Printf("Timed out before launching goroutine for URL %s\n", url)
				atomic.AddInt32(&failCount, 1)
				return
			}

			// 等待获取一个令牌
			err := rateLimiter.Wait(context.Background())
			if err != nil {
				log.Printf("Failed to acquire token for request : %v\n", err)
				atomic.AddInt32(&failCount, 1)
				<-limit
				return
			}

			get, err := requestLimiterOptionMap[ll.HttpClient].Get(url)
			if err != nil {
				log.Println(err)
				atomic.AddInt32(&failCount, 1)
				// TODO 可以将失败的请求传入一个 failChan 然后做失败后的任务处理
				return
			}

			log.Println(string(get))
			atomic.AddInt32(&successCount, 1)
		}(url)

		// 因为给 httpClient 加了重试功能，所以下面的代码不需要了
		//time.Sleep(time.Duration((1/ll.batchSize)*1000) * time.Millisecond) // 控制协程之间的时间间隔
	}

	wg.Wait()

	log.Printf("总任务数: %d\n", len(urls))
	log.Printf("成功的任务数: %d\n", atomic.LoadInt32(&successCount))
	log.Printf("失败的任务数: %d\n", atomic.LoadInt32(&failCount))
}

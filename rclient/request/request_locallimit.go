package request

import (
	"io"
	"net"
	"net/http"
	"runtime"
	"time"
)

var defaultHttpClients = newDefaultHttpClient() // 设置成全局变量可让所有协程使用同一个 httpClient

type LocalRateLimitRequest struct {
	*baseHttpRequest
}

// LocalRateLimitRequestFactory 创建一个带速率限制的请求对象
func LocalRateLimitRequestFactory(rateLimiterOption *RateLimiterRequest) Request {
	return &LocalRateLimitRequest{
		baseHttpRequest: &baseHttpRequest{
			baseRateLimiterName: rateLimiterOption.RateLimiterName,
		},
	}
}

func (nr *LocalRateLimitRequest) Get(url string) ([]byte, error) {
	return nr.makeRequest("GET", url, nil)
}

func (nr *LocalRateLimitRequest) Post(url string, body io.Reader) ([]byte, error) {
	return nr.makeRequest("POST", url, body)
}

// newDefaultHttpClient 创建一个带有连接池和简单重试功能的 http client
func newDefaultHttpClient() *http.Client {
	return &http.Client{
		Transport: defaultPooledTransport(),
		Jar:       jar,
	}
}

func defaultPooledTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	dialContext := dialer.DialContext

	// 在服务器支持长连接时： MaxIdleConns  MaxIdleConnsPerHost MaxConnsPerHost 都设置为 1 时客户端只开启一个
	// 端口来发送请求，链接得到复用,这时服务器也只开一个端口来处理请求。
	// 当服务器开启了两个端口来处理请求时，要想客户端也只开启一个端口发送请求时需要设置 MaxIdleConnsPerHost MaxConnsPerHost 为 1
	// 且不能设置 MaxIdleConns=1 ，如果设置了这个，客户端还是会开启多个端口来发送请求。
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		MaxConnsPerHost:       25,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
	return transport
}

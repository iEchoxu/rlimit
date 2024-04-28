package request

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
)

func GetRetryHttpClient() *http.Client {
	return NewRetryHttpClient().StandardClient()
}

// NewRetryHttpClient 使用http连接池
// 自带重试功能
func NewRetryHttpClient() *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil
	return retryClient
}

var retryHttpClient = GetRetryHttpClient()

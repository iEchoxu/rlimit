# rLimit

rLimit, 简单的 http 请求限流案例。

## 如何使用

- `go mod tidy`
- 启动 server: `go run rlimit/rserver/cmd/main.go`
- 启动 client: `go run rlimit/rclient/cmd/main.go`

## 介绍

- config 目录下是 client 运行时所需要的参数，可在 run_config.go 中进行配置
- 运行模式有: LocalLimiterType、ThrottledLimiterType, 详细代码在 limitertest 目录中
- 请求方式: 使用 Throttled 库提供的 httpClient 以及使用 http 标准库实现的带简单重试功能的 httpClient
- LocalLimiterType 使用多协程 + channel + rate.NewLimiter(golang.org/x/time/rate) + 带重试功能的 httpClient 实现
- ThrottledLimiterType 使用  go-retryablehttp + throttled/v2 库实现限流与重试

## 功能
- 支持请求限流、重试
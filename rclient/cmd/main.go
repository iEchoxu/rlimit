package main

import (
	"fmt"
	"rlimit/rclient/config"
	"time"
)

func main() {
	start := time.Now()
	r := config.NewRunner(config.RunMode, config.HttpClient)
	r.Run()
	duration := time.Since(start)
	fmt.Printf("总耗时：%s\n\n", duration)
}

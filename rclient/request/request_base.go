package request

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// httpClient 类型, 值为: RetryableHttpClient 或 DefaultHttpClient
// RetryableHttpClient: 使用 go-retryablehttp 库提供的 httpClient
// DefaultHttpClient: 使用定制过的 http 标准库提供的 httpClient：有连接池以及简单的重试功能
const (
	RetryableHttpClient HttpClientType = iota + 1
	DefaultHttpClient
)

var httpClientMap = make(map[string]*http.Client)
var jar http.CookieJar // 全局变量，用于在所有请求间共享cookie

type HttpClientType int

type baseHttpRequest struct {
	baseRateLimiterName string
}

func init() {
	// 在程序初始化时创建cookie jar
	var err error
	jar, err = cookiejar.New(nil)
	if err != nil {
		log.Println(err)
	}
}

func (bh *baseHttpRequest) makeRequest(method, url string, body io.Reader) (bodyContent []byte, err error) {
	var httpClient *http.Client

	currentHttpclient := bh.getCurrentClient()
	if currentHttpclient == "defaultHttpClient" {
		httpClient = httpClientMap[currentHttpclient]
		bodyContent, err = bh.defaultHttpClientDo(method, url, body, httpClient)
		return
	}

	httpClient = httpClientMap[currentHttpclient]
	bodyContent, err = bh.retryableHttpClientDo(method, url, body, httpClient)
	return
}

func (bh *baseHttpRequest) retryableHttpClientDo(method, url string, body io.Reader, client *http.Client) (bodyContent []byte, err error) {
	limited, _, err := GetRateLimiter(bh.baseRateLimiterName)
	if err != nil {
		log.Println(err)
	}

	log.Println("当前限流器:", bh.baseRateLimiterName)

	if limited {
		//log.Println(result.Limit, result.Remaining, result.RetryAfter, result.ResetAfter)
		randSec := 0.1
		log.Printf("Too Many Requests,waiting %v second", randSec)
		time.Sleep(time.Second * time.Duration(randSec))
	}

	bodyContent, err = bh.httpClientDo(method, url, body, client)
	return
}

func (bh *baseHttpRequest) defaultHttpClientDo(method, url string, body io.Reader, client *http.Client) (bodyContent []byte, err error) {
	var respContent []byte
	var respError error

	for i := 0; i < 10; i++ {
		resp, respErr := bh.httpClientDo(method, url, body, client)
		if respErr != nil {
			log.Println("Too Many Requests: retry......")
			time.Sleep(1 * time.Second)
			continue
		}

		respContent = resp
		respError = respErr
		break
	}

	return respContent, respError
}

func (bh *baseHttpRequest) httpClientDo(method, url string, body io.Reader, httpClient *http.Client) (bodyContent []byte, err error) {
	// 创建请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error building %s request: %w", method, err)
	}

	// 设置请求头，指定 Content-Type 和 UserAgent
	req.Header.Set("User-Agent", bh.randomUserAgent())
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// 获取HTTP客户端
	client := httpClient

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending %s request: %w", method, err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: response status code %d", resp.StatusCode)
	}

	// 读取响应体
	bodyContent, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// 打印响应体
	//log.Println("Response Body:", string(bodyContent))

	// 打印所有的cookie
	//for _, cookie := range jar.Cookies(req.URL) {
	//	fmt.Printf("%s=%s\n", cookie.Name, cookie.Value)
	//}

	return bodyContent, nil
}

func (bh *baseHttpRequest) randomUserAgent() string {
	var userAgent = []string{
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	}

	rand.NewSource(time.Now().UnixNano())
	useAgentRand := userAgent[rand.Intn(len(userAgent))]

	return useAgentRand
}

func (bh *baseHttpRequest) getCurrentClient() (clientName string) {
	for key := range httpClientMap {
		clientName = key
	}

	return
}

// SetClient 设置使用哪个 httpClient
// RetryableHttpClient 使用  go-retryablehttp 库提供的 httpClient
// DefaultHttpClient 使用改造后的标准库 http 提供的 httpClient
func SetClient(clientType HttpClientType) {
	if clientType == RetryableHttpClient {
		httpClientMap["retryableHttpClient"] = retryHttpClient
		log.Println("current client: retryableHttpClient")
		return
	}

	httpClientMap["defaultHttpClient"] = defaultHttpClients
	log.Println("current client: defaultHttpClient")
}

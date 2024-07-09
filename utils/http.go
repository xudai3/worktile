package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HeaderOption struct {
	Name  string
	Value string
}

type QueryParameter struct {
	Key   string
	Value interface{}
}

type HttpClient struct {
	Client     *http.Client
	MaxRetries int
	RetryWait  time.Duration
}

const (
	ContentTypeJson = "application/json"
)

func NewHttpClient(client *http.Client, maxRetries int, retryWait int) *HttpClient {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	if retryWait <= 0 {
		retryWait = 2
	}
	return &HttpClient{
		Client:     client,
		MaxRetries: maxRetries,
		RetryWait:  time.Duration(retryWait) * time.Second,
	}
}

func (c *HttpClient) retryHttpRequest(reqFunc func() (*http.Request, error), respHandler func(*http.Response, error) ([]byte, error)) ([]byte, error) {
	retries := 0
	for {
		req, err := reqFunc()
		if err != nil {
			return nil, err
		}

		resp, err := c.Client.Do(req)
		if err != nil {
			// 处理网络请求错误，如需重试则继续循环
			if retries >= c.MaxRetries {
				return nil, fmt.Errorf("request failed after %d retries: %v", retries, err)
			}
			retries++
			time.Sleep(c.RetryWait)
			continue
		}
		defer resp.Body.Close()

		body, err := respHandler(resp, err)
		if err == nil {
			return body, nil
		}

		// 根据响应错误判断是否需要重试
		if shouldRetryOnResponseError(resp.StatusCode) {
			if retries >= c.MaxRetries {
				return nil, fmt.Errorf("request failed after %d retries with status code: %d", retries, resp.StatusCode)
			}
			retries++
			time.Sleep(c.RetryWait)
			continue
		}

		return nil, err
	}
}

// shouldRetryOnResponseError 根据响应状态码判断是否应该重试
func shouldRetryOnResponseError(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError:
		return true
	default:
		return false
	}
}

// func (c *HttpClient) GetRequest(url string, headerOptions ...HeaderOption) ([]byte, error) {
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, headerOption := range headerOptions {
// 		req.Header.Set(headerOption.Name, headerOption.Value)
// 	}
// 	resp, err := c.Client.Do(req)
// 	defer respDeferClose(resp)
// 	return respHandle(resp, err)
// }

func (c *HttpClient) GetRequest(url string, headerOptions ...HeaderOption) ([]byte, error) {
	reqFunc := func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		for _, option := range headerOptions {
			req.Header.Set(option.Name, option.Value)
		}
		return req, nil
	}
	return c.retryHttpRequest(reqFunc, respHandle)
}

func (c *HttpClient) Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) ([]byte, error) {
	fullUrl := url + ConvertToQueryParams(params)
	return c.GetRequest(fullUrl, headerOptions...)
}

func (c *HttpClient) GetWithParam(url string, param string, headerOptions ...HeaderOption) ([]byte, error) {
	fullUrl := url + "/" + param
	return c.GetRequest(fullUrl, headerOptions...)
}

// func (c *HttpClient) PostRequest(url string, body interface{}, headerOptions ...HeaderOption) ([]byte, error) {
// 	buf, err := json.Marshal(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", ContentTypeJson)
// 	for _, headerOption := range headerOptions {
// 		req.Header.Set(headerOption.Name, headerOption.Value)
// 	}
// 	resp, err := c.Client.Do(req)
// 	defer respDeferClose(resp)
// 	return respHandle(resp, err)
// }

func (c *HttpClient) PostRequest(url string, body interface{}, headerOptions ...HeaderOption) ([]byte, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reqFunc := func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", ContentTypeJson)
		for _, option := range headerOptions {
			req.Header.Set(option.Name, option.Value)
		}
		return req, nil
	}
	return c.retryHttpRequest(reqFunc, respHandle)
}

func (c *HttpClient) Post(url string, param string, body interface{}, headerOptions ...HeaderOption) ([]byte, error) {
	fullUrl := url + "/" + param
	return c.PostRequest(fullUrl, body, headerOptions...)
}

func BuildTokenHeaderOptions(accessToken string) HeaderOption {
	return HeaderOption{
		Name:  "authorization",
		Value: "Bearer " + accessToken,
	}
}

func respHandle(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ConvertToQueryParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		if v == nil {
			continue
		}
		buffer.WriteString(fmt.Sprintf("%s=%v&", k, v))
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}

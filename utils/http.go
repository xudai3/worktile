package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HeaderOption struct {
	Name string
	Value string
}

type QueryParameter struct {
	Key string
	Value interface{}
}

type HttpClient struct {
	Client *http.Client
}

const (
	ContentTypeJson = "application/json"
)

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{
		Client: client,
	}
}

func (c *HttpClient) GetRequest(url string, headerOptions ...HeaderOption) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}
	resp, err := c.Client.Do(req)
	defer respDeferClose(resp)
	return respHandle(resp, err)
}

func (c *HttpClient) Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) ([]byte, error) {
	fullUrl := url + ConvertToQueryParams(params)
	fmt.Printf("full url:%s", fullUrl)
	return c.GetRequest(fullUrl, headerOptions...)
}

func (c *HttpClient) PostRequest(url string, body interface{}, headerOptions ...HeaderOption) ([]byte, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ContentTypeJson)
	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}
	resp, err := c.Client.Do(req)
	defer respDeferClose(resp)
	return respHandle(resp, err)
}

func (c *HttpClient) Post(url string, param string, body interface{}, headerOptions ...HeaderOption) ([]byte, error) {
	fullUrl := url + "/" + param
	fmt.Printf("full url:%s", fullUrl)
	return c.PostRequest(fullUrl, body, headerOptions...)
}

func BuildTokenHeaderOptions(accessToken string) HeaderOption {
	return HeaderOption{
		Name: "authorization",
		Value: "Bearer " + accessToken,
	}
}

func respDeferClose(resp *http.Response) {
	if resp != nil {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}
}

func respHandle(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//respBody := string(b)
	return b, nil
}

func ConvertToQueryParams(params map[string]interface{}) string {
	//bs, err := json.Marshal(params)
	//if err != nil {
	//	return ""
	//}
	//params = map[string]interface{}{}
	//_ = json.Unmarshal(bs, &params)
	//
	if params == nil || len(params) == 0 {
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


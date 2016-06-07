package rhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

const (
	HTTP_RETURN_TYPE_MAP    = 1
	HTTP_RETURN_TYPE_SLICE  = 2
	HTTP_RETURN_TYPE_STRING = 3
)

// HttpClientGet http get from http client
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func HttpClientGet(c *http.Client, urlget string, retType int) (interface{}, error) {
	resp, err := c.Get(urlget)
	if err != nil {
		fmt.Printf("[ HttpClientGet() ]\r\n\tGet: %s\r\n", err)
		return nil, err
	}
	return GetResponseDecode(resp, retType)
}

// HttpGet 
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func HttpGet(urlget string, retType int) (interface{}, error) {
	resp, err := http.Get(urlget)
	if err != nil {
		fmt.Printf("[ HttpGet() ]\r\n\tGet: %s\r\n", err)
		return nil, err
	}
	return GetResponseDecode(resp, retType)
}

// HttpProxyGet 
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func HttpProxyGet(pgUrl, proxyIP, proxyPort string, retType int, auth ...*proxy.Auth) (interface{}, error) {

	client := GetProxyClient(proxyIP, proxyPort, auth...)

	resp, err := client.Get(pgUrl)
	if err != nil {
		fmt.Printf(" ** HTTPProxyGet() Error\r\n\tGet: %s\r\n", err)
		return nil, err
	}

	return GetResponseDecode(resp, retType)
}

// HttpPostForm http-post form
// retType = 1  => map[string]interface{};
// retType = 2  => []interface{};
// retType = 3  => string;
func HttpPostForm(address string, retType int, data string) (interface{}, error) {

	resp, err := http.Post(address, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		fmt.Printf(" ** HTTPPost() Error\r\n\tPost: %s\r\n", err)
		return nil, err
	}

	return GetResponseDecode(resp, retType)
}

// HttpPostJson http-post json
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func HttpPostJson(address string, retType int, data []byte) (interface{}, error) {

	resp, err := http.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Printf("ERROR [ HTTPPOST.json() ]\r\n\tPost: %s\r\n", err)
		return nil, err
	}

	return GetResponseDecode(resp, retType)
}

// HttpClientPostForm
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func HttpClientPostForm(c *http.Client, address string, retType int, data string) (interface{}, error) {
	resp, err := c.Post(address, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		fmt.Printf("[ HttpClientPostForm() ]\r\n\tPost: %s\r\n\tData: %s\r\n", err, data)
		return nil, err
	}
	return GetResponseDecode(resp, retType)
}

// HttpClientPostJson
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func HttpClientPostJson(c *http.Client, address string, retType int, data []byte) (interface{}, error) {
	resp, err := c.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Printf("[ HttpClientPostJson() ]\r\n\tPost: %s\r\n\tData: %s\r\n", err, string(data))
		return nil, err
	}
	return GetResponseDecode(resp, retType)
}

// GetResponseDecode
// retType = 1  => map[string]interface{}
// retType = 2  => []interface{}
// retType = 3  => string
func GetResponseDecode(resp *http.Response, retType int) (interface{}, error) {

	if resp == nil {
		return nil, errors.New("http.Response is nil!")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(" ** HTTP Response ERROR\r\n\tReadAll: %s\r\n", err)
		return nil, err
	}

	// fmt.Println("[ HTTPGet().body ]\r\n\t", string(body))

	switch retType {
	case HTTP_RETURN_TYPE_MAP:
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Printf(" ** HTTP Response ERROR\r\n\tUnmarshal: %s\r\n", err)
			fmt.Println("[ HTTPGet().body ]\r\n\t", string(body))
			return nil, err
		}
		return result, nil
	case HTTP_RETURN_TYPE_SLICE:
		result := make([]interface{}, 0)
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Printf(" ** HTTP Response ERROR\r\n\tUnmarshal: %s\r\n", err)
			fmt.Println("[ HTTPGet().body ]\r\n\t", string(body))
			return nil, err
		}
		return result, nil
	case HTTP_RETURN_TYPE_STRING:
		result := string(body)
		return result, nil
	}
	panic("Error retType")
}

// GetHttpClient create a http client
func GetHttpClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:  time.Second * 20,
		Deadline: time.Now().Add(30 * time.Second),
		// KeepAlive: time.Second * 30,
	}
	trans := &http.Transport{
		Dial: dialer.Dial,
		ResponseHeaderTimeout: 20 * time.Second,

		DialTLS:             dialer.Dial,
		TLSHandshakeTimeout: 20 * time.Second,
	}

	ret := &http.Client{
		Transport: trans,
		Timeout:   30 * time.Second,
	}
	return ret
}

// GetProxyClient create http client with proxy
func GetProxyClient(proxyIP, proxyPort string, auth ...*proxy.Auth) *http.Client {

	proxyurl := proxyIP + ":" + proxyPort
	var author *proxy.Auth = nil
	if auth != nil {
		author = auth[0]
	}
	dialer, err := proxy.SOCKS5("tcp", proxyurl, author,
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	)

	if err != nil {
		fmt.Printf(" ** Socket5ProxyClient() Error\r\n\tproxy.SOCKS5: %s\r\n", err)
		return nil
	}

	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{Transport: transport}
}

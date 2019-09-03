package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const Timeout = 20

type Options struct {
	Client    *http.Client
	Method    string
	URL       string
	Data      interface{}
	Json      bool
	Timeout   int
	Headers   map[string]string
	Retry     int
}

type Response struct {
	*http.Response
	Selector      *Selector
	Bytes         []byte
	History       []*http.Request
}

var DefaultHeaders = map[string]string{
	"Accept": "*/*",
	"Accept-Encoding": "",
	"Connection": "keep-alive",
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
}


func NodeNotFound(xpaths []string) error {
	return fmt.Errorf("HTML node not found. xpaths: %s", strings.Join(xpaths,","))
}


func InitCookie(client *http.Client) *http.Client {
	jar,_ := cookiejar.New(nil)
	client.Jar = jar
	return client
}

func InitClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,   // 防止保持长连接
		},
	}
	return InitCookie(client)
}

func Request(options Options) (*Response, error) {
	var bodyReader io.Reader
	method := options.Method
	method = strings.ToUpper(method)

	//set client
	client := options.Client
	if client == nil {
		client = InitClient()
	}

	//set history
	var history []*http.Request
	if client.CheckRedirect == nil {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			history = via
			if len(via) > 10 {
				return errors.New("stopped after 10 redirects")
			}
			history = append(history, req)
			return nil
		}
	}


	//set headers
	headers := http.Header{}

	for name, value := range DefaultHeaders {
		headers.Set(name, value)
	}

	if options.Method == "POST" {
		headers.Set("Content-Type","application/x-www-form-urlencoded")
	}

	if options.Headers != nil {
		for name, value := range options.Headers {
			headers.Set(name,value)
		}
	}

	// set body
	if options.Data != nil {
		switch h := options.Data.(type) {
		case string:
			bodyReader = strings.NewReader(h)
		case url.Values:
			bodyReader = strings.NewReader(h.Encode())
		case map[string]interface{}:
			if options.Json {
				if bodyBytes, err := json.Marshal(h); err != nil {
					panic(err)
				} else {
					bodyReader = strings.NewReader(string(bodyBytes))
				}
			}
		}

	}
	//fmt.Println(bodyReader)
	//new request
	request, err := http.NewRequest(method, options.URL, bodyReader)
	if err != nil {
		panic(err)
	}

	//set timeout
	timeout := Timeout
	if options.Timeout != 0 {
		timeout = options.Timeout
	}

	client.Timeout = time.Duration(timeout) * time.Second

	request.Header = headers

	retry := options.Retry
	if retry == 0 {
		retry = 1
	}

	var response *http.Response
	var response_body []byte

	for i := 0; i < retry; i++ {
		response, err = client.Do(request)
		if err == nil && response.StatusCode <  500 {
			response_body, err = _ReadResponseBody(response)
			if err == nil {
				break
			}
		}

		if err == nil && i+1 < retry {
			response.Body.Close()
		}
		fmt.Println(err)
		time.Sleep(time.Second * 1)
	}

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//defer response.Body.Close()
	//fmt.Println(response.StatusCode)

	//if response.StatusCode == 200 {
	//	body, err := ioutil.ReadAll(response.Body)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	fmt.Println(string(body))
	//
	//}
	//body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	return &Response{response,&Selector{body:response_body},response_body,history}, err
}

func _ReadResponseBody(r *http.Response) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 65536))
	_, err := io.Copy(buffer, r.Body)
	if err != nil {
		buffer.Truncate(0)
		return nil, err
	}
	temp := buffer.Bytes()
	length := len(temp)
	var body []byte
	//are we wasting more than 10% space?
	if cap(temp) > (length + length/10) {
		body = make([]byte, length)
		copy(body, temp)
	} else {
		body = temp
	}
	return body, nil
}

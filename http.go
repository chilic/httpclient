package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"reflect"

	gohttp "net/http"
	"github.com/gostuff/json"
)

type Response struct {
	gohttp.Response
}

//
// Check if the input value is a "primitive" that can be safely stringified
//
func canStringify(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

//
// Given a base URL and a bag of parameteters returns the URL with the encoded parameters
//
func URLWithPathParams(base string, path string, params map[string]interface{}) (u *url.URL) {

	u, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}

	if len(path) > 0 {
		u, err = u.Parse(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	q := u.Query()

	for k, v := range params {
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.Slice:
			if val.IsNil() { // TODO: add an option to ignore empty values
				q.Set(k, "")
				continue
			}
			fallthrough

		case reflect.Array:
			for i := 0; i < val.Len(); i++ {
				av := val.Index(i)

				if canStringify(av) {
					q.Add(k, fmt.Sprintf("%v", av))
				}
			}

		default:
			if canStringify(val) {
				q.Set(k, fmt.Sprintf("%v", v))
			} else {
				log.Fatal("Invalid type ", val)
			}
		}
	}

	u.RawQuery = q.Encode()
	return u
}

func URLWithParams(base string, params map[string]interface{}) (u *url.URL) {
	return URLWithPathParams(base, "", params)
}

//
// http.Get with params
//
func Get(url string, params map[string]interface{}) (*Response, error) {
	resp, err := gohttp.Get(URLWithParams(url, params).String())
	if err == nil {
		return &Response{*resp}, nil
	} else {
		return nil, err
	}
}

//
// http.Post with params
//
func Post(url string, params map[string]interface{}) (*Response, error) {
	resp, err := gohttp.PostForm(url, URLWithParams(url, params).Query())
	if err == nil {
		return &Response{*resp}, nil
	} else {
		return nil, err
	}
}

//
//  Read the body
//
func (resp *Response) Content() []byte {
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

//
//  Try to parse the response body as JSON
//
func (resp *Response) Json() *json.Jobj {
	return json.Loads(resp.Content())
}

////////////////////////////////////////////////////////////////////////
//
// http.Client with some defaults and stuff
//

type HttpClient struct {
	client *gohttp.Client

        BaseURL   *url.URL
	UserAgent string
	Headers   map[string]string
}

func NewHttpClient(base string) (httpClient *HttpClient) {
	httpClient = new(HttpClient)
	httpClient.client = &gohttp.Client{}

        if u, err := url.Parse(base); err != nil {
            log.Fatal(err)
        } else {
            httpClient.BaseURL = u
        }

	return
}

func (self *HttpClient) Request(method string, urlpath string, body io.Reader) (req *gohttp.Request) {
        if u, err := self.BaseURL.Parse(urlpath); err != nil {
            log.Fatal(err)
        } else {
            urlpath = u.String()
        }

	req, err := gohttp.NewRequest(method, urlpath, body)
	if err != nil {
		log.Fatal(err)
	}

	if len(self.UserAgent) > 0 {
		req.Header.Set("User-Agent", self.UserAgent)
	}

	for k, v := range self.Headers {
		req.Header.Set(k, v)
	}

	return
}

func (self *HttpClient) Do(req *gohttp.Request) (*Response, error) {
	resp, err := self.client.Do(req)
	if err == nil {
		return &Response{*resp}, nil
	} else {
		return nil, err
	}
}

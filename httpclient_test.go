package httpclient

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"testing"
)

const (
	BASE_URL     = "http://httpbin.org/"
	GET_URL      = BASE_URL + "get"
	POST_URL     = BASE_URL + "post"
	REDIRECT_URL = BASE_URL + "redirect-to?url=http://example.com"
)

var (
	params = map[string]interface{}{
		"string": "one",
		"int":    2,
		"number": 3.14,
		"bool":   true,
		"list":   []string{"one", "two", "three"},
		"empty":  []int{},
	}
)

func TestURLWithParams(test *testing.T) {
	test.Log(URLWithParams(GET_URL, params))
}

func TestURLWithPathParams(test *testing.T) {
	test.Log(URLWithPathParams(GET_URL, "another", nil))
	test.Log(URLWithPathParams(GET_URL+"/", "another", nil))
}

func TestGet(test *testing.T) {
	resp, err := Get(GET_URL, nil)
	if err != nil {
		test.Error(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		test.Log(string(body))
	}
}

func TestGetWithParams(test *testing.T) {
	resp, err := Get(GET_URL, params)
	if err != nil {
		test.Error(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		test.Log(string(body))
	}
}

func TestPostWithParams(test *testing.T) {

	resp, err := Post(POST_URL, params)
	if err != nil {
		test.Error(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		test.Log(string(body))
	}
}

func TestGetJSON(test *testing.T) {
	resp, err := Get(GET_URL, params)
	if err != nil {
		test.Error(err)
	} else {
		test.Log(resp.Json().Map())
	}
}

func TestClient(test *testing.T) {
	client := NewHttpClient(BASE_URL)
	client.UserAgent = "TestClient 0.1"
	client.Verbose = true

	req := client.Request("GET", "get", nil, nil)
	resp, err := client.Do(req)
	test.Log(err, string(resp.Content()))

	req = client.Request("POST", "post", bytes.NewBuffer([]byte("the body")), nil)
	resp, err = client.Do(req)
	test.Log(err, string(resp.Content()))
}

func TestClientGet(test *testing.T) {
	client := NewHttpClient(BASE_URL)
	client.UserAgent = "TestClient 0.1"
	client.Verbose = true

	resp, err := client.Get("get", nil, nil)
	test.Log(err, string(resp.Content()))
}

func TestClientPost(test *testing.T) {
	client := NewHttpClient(BASE_URL)
	client.UserAgent = "TestClient 0.1"
	client.Verbose = true

	data := bytes.NewBuffer([]byte("the body"))

	resp, err := client.Post("post", data, map[string]string{"Content-Type": "text/plain", "Content-Disposition": "attachment;filename=test.txt", "Content-Length": strconv.Itoa(data.Len())})
	test.Log(err, string(resp.Content()))
}

func TestClientGetRedirect(test *testing.T) {
	client := NewHttpClient(REDIRECT_URL)
	client.UserAgent = "TestClient 0.1"
	client.Verbose = true

	resp, err := client.Get("", nil, nil)
	test.Log(err, string(resp.Content()))
}

func TestClientHeadRedirect(test *testing.T) {
	client := NewHttpClient(REDIRECT_URL)
	client.UserAgent = "TestClient 0.1"
	client.Verbose = true

	resp, err := client.Head("", nil, nil)
	location := "<no-location>"
	if resp != nil {
		location = resp.Header.Get("Location")
	}
	test.Log(err, location)
}

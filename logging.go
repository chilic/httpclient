package httpclient

import (
	"fmt"
        "io"
        "os"
	"net/http"
	"net/http/httputil"
	"strconv"
)

// A transport that prints request and response

type LoggingTransport struct {
	t            *http.Transport
	requestBody  bool
	responseBody bool
}

func (lt *LoggingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	dreq, _ := httputil.DumpRequest(req, lt.requestBody)
	fmt.Println("REQUEST:", strconv.Quote(string(dreq)))
	fmt.Println("")

	resp, err = lt.t.RoundTrip(req)
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		dresp, _ := httputil.DumpResponse(resp, lt.responseBody)
		fmt.Println("RESPONSE:", string(dresp))
	}

	fmt.Println("")
	return
}

// Enable logging requests/response headers
//
// if requestBody == true, also log request body
// if responseBody == true, also log response body
func StartLogging(requestBody, responseBody bool) {
	http.DefaultTransport = &LoggingTransport{&http.Transport{}, requestBody, responseBody}
}

// Disable logging requests/responses
func StopLogging() {
	http.DefaultTransport = &http.Transport{}
}

// A Reader that "logs" progress

type ProgressReader struct {
    r   io.Reader
    c   [1]byte
    threshold int
    curr int
}

func NewProgressReader(r io.Reader, c byte, threshold int) *ProgressReader {
    if c == 0 {
        c = '.'
    }
    if threshold <= 0 {
        threshold = 10240
    }
    p := &ProgressReader{r: r, c: [1]byte{ c }, threshold: threshold, curr: 0}
    return p
}

func (p *ProgressReader) Read(b []byte) (int, error) {
    n, err := p.r.Read(b)

    p.curr += n

    if err == io.EOF {
        os.Stdout.Write([]byte{'\n'})
        os.Stdout.Sync()
    } else if p.curr >= p.threshold {
        p.curr -= p.threshold

        os.Stdout.Write(p.c[:])
        os.Stdout.Sync()
    }

    return n, err
}

func (p *ProgressReader) Close() error {
    if rc, ok := p.r.(io.ReadCloser); ok {
        return rc.Close()
    } else {
        return nil
    }
}

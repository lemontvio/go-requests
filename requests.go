package go_requests

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const UserAgent = `go-requests/0.0.1`

type Requests struct {
	proxy     string
	useragent string
	timeout   int
}

func New() *Requests {
	return &Requests{useragent: UserAgent}
}

func (requests *Requests) request(method string, raw string, options *Options) (response *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}

	if requests.proxy != "" {
		transport.Proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(requests.proxy)
		}
	}

	if requests.timeout > 0 {
		transport.IdleConnTimeout = time.Second * time.Duration(requests.timeout)
	}

	client := &http.Client{Transport: transport}

	if options != nil {
		options = options.Copy()
	} else {
		options = &Options{}
	}

	raw = strings.TrimSuffix(raw, "?")
	if options.Values != nil {
		raw += "?" + options.Values.Encode()
	}

	var req *http.Request
	if req, err = http.NewRequest(method, raw, options.Body); err != nil {
		return
	}

	if options.Header == nil {
		options.Header = http.Header{}
	}

	if v, ok := options.Header["user-agent"]; ok {
		delete(options.Header, "user-agent")
		options.Header.Set("User-Agent", v[0])
	}

	if _, ok := options.Header["User-Agent"]; !ok && options.UserAgent != "" {
		options.Header.Set("User-Agent", options.UserAgent)
	}

	if _, ok := options.Header["User-Agent"]; !ok {
		options.Header.Set("User-Agent", requests.useragent)
	}

	if v, ok := options.Header["referer"]; ok {
		delete(options.Header, "referer")
		options.Header.Set("Referer", v[0])
	}

	if _, ok := options.Header["Referer"]; !ok && options.Referer != "" {
		options.Header.Set("Referer", options.Referer)
	}

	if method != "GET" {
		if len(options.Header["Content-Type"]) < 1 || options.Header["Content-Type"][0] == "" {
			options.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	req.Header = options.Header
	if v, ok := options.Header["cookie"]; ok {
		delete(options.Header, "cookie")
		options.Header.Set("Cookie", v[0])
	}

	if _, ok := options.Header["Cookie"]; !ok && options.Cookie != nil {
		for _, cookie := range options.Cookie {
			req.AddCookie(cookie)
		}
	}

	response, err = client.Do(req)
	client.CloseIdleConnections()
	return
}

func (requests *Requests) Get(raw string, options *Options) (*http.Response, error) {
	return requests.request("GET", raw, options)
}

func (requests *Requests) Post(raw string, options *Options) (*http.Response, error) {
	return requests.request("POST", raw, options)
}

func (requests *Requests) Put(raw string, options *Options) (*http.Response, error) {
	return requests.request("PUT", raw, options)
}

func (requests *Requests) Delete(raw string, options *Options) (*http.Response, error) {
	return requests.request("DELETE", raw, options)
}

func (requests *Requests) SetProxy(raw string) {
	if raw != "" {
		requests.proxy = raw
	}
}

func (requests *Requests) SetUserAgent(raw string) {
	requests.useragent = raw
}

func (requests *Requests) SetTimeout(timeout int) {
	requests.timeout = timeout
}

//
// =====================================================================================================================
//

var std = New()

func Get(raw string, options *Options) (*http.Response, error) {
	return std.Get(raw, options)
}

func Post(raw string, options *Options) (*http.Response, error) {
	return std.Post(raw, options)
}

func Put(raw string, options *Options) (*http.Response, error) {
	return std.Put(raw, options)
}

func Delete(raw string, options *Options) (*http.Response, error) {
	return std.Delete(raw, options)
}

func SetProxy(raw string) {
	std.proxy = raw
}

func SetUserAgent(raw string) {
	std.useragent = raw
}

func SetTimeout(timeout int) {
	std.timeout = timeout
}

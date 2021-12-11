package requests

import (
	"io"
	"net/http"
	"net/url"
)

type Options struct {
	Values    *Values
	Body      io.Reader
	Header    http.Header
	UserAgent string
	Referer   string
	Cookie    []*http.Cookie
	Proxy     *url.URL
}

func (options *Options) Copy() *Options {
	opt := &Options{}
	opt.Values = options.Values
	opt.Body = options.Body
	opt.Header = options.Header.Clone()
	opt.UserAgent = options.UserAgent
	opt.Referer = options.Referer
	opt.Cookie = options.Cookie
	opt.Proxy = options.Proxy

	return opt
}

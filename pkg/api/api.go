package api

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/base"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost          = "POST"
	MethodPut           = "PUT"
	MethodDelete        = "DELETE"
)

type Api struct {
	req       *gorequest.SuperAgent
	baseURL   string
	AccessKey string
	secretKey string
}

func New(accessKey, secretKey string) *Api {
	baseURL := base.BaseUrl()
	req := gorequest.New()
	req.Type("json")
	req.Header.Set("x-onf-client", viper.GetString("app.name"))
	req.Header.Set("x-onf-version", viper.GetString("app.version"))
	return &Api{
		baseURL:   baseURL,
		req:       req,
		AccessKey: accessKey,
		secretKey: secretKey,
	}
}

type RequestOptions struct {
	Params map[string]string
	Body   interface{}
	Files  map[string]string
}

func (a *Api) Request(method Method, path string, opts *RequestOptions) *gorequest.SuperAgent {
	r := a.req.Clone()

	u, _ := url.Parse(fmt.Sprintf("%s%s", a.baseURL, path))
	m := string(method)
	r = r.CustomMethod(m, u.String())

	if opts != nil {
		if opts.Params != nil {
			for k, v := range opts.Params {
				u.Query().Add(k, v)
			}
		}
		if opts.Body != nil {
			r = r.Send(opts.Body)
			r.Header.Set("content-md5", contentChecksum(r.Data))
		}
	}

	// set date in GMT format
	r.Header.Set("date", time.Now().UTC().Format(http.TimeFormat))
	signature := a.GetSign(r.Method, u.RequestURI(), r.Header)
	r.Header.Set("authorization", fmt.Sprintf("ONF %s:%s", a.AccessKey, signature))

	return r
}

func (a *Api) Upload(path string, opts *RequestOptions) *gorequest.SuperAgent {
	r := a.Request(MethodPost, path, opts)
	if len(opts.Files) > 0 {
		r.Type("multipart")
		for n, f := range opts.Files {
			r.SendFile(f, n, "files")
		}
	}
	return r
}

func (a *Api) GetSign(method, uri string, header http.Header) string {
	data := strings.Join([]string{
		method,
		header.Get("content-type"),
		header.Get("content-md5"),
		header.Get("date"),
		canonicalize(header, "x-onf-"),
		uri,
	}, "\n")
	return sign(data, a.secretKey)
}

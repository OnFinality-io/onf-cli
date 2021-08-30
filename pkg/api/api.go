package api

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/utils/gorequest"
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
	version   int
	AccessKey string
	secretKey string
}

func New(accessKey, secretKey string, baseURL string) *Api {
	req := gorequest.New()
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-onf-client", viper.GetString("app.name"))
	req.Header.Set("x-onf-version", viper.GetString("app.version"))
	//req.Debug = true
	return &Api{
		baseURL:   baseURL,
		req:       req,
		AccessKey: accessKey,
		secretKey: secretKey,
		version:   1,
	}
}

func (a *Api) Clone() *Api {
	return &Api{
		baseURL:   a.baseURL,
		req:       a.req,
		AccessKey: a.AccessKey,
		secretKey: a.secretKey,
		version:   a.version,
	}
}

func (a *Api) Ver2() *Api {
	cloneApi := a.Clone()
	cloneApi.version = 2
	return cloneApi
}

type RequestOptions struct {
	Params map[string]string
	Body   interface{}
	Files  map[string]string
}

func (a *Api) Request(method Method, path string, opts *RequestOptions) *gorequest.SuperAgent {
	r := a.req.Clone()
	fullApi := fmt.Sprintf("%s/v%d", a.baseURL, a.version)
	u, _ := url.Parse(fmt.Sprintf("%s%s", fullApi, path))
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
	r.Header.Del("content-type")
	r.Type("multipart")
	for n, f := range opts.Files {
		r.SendFile(f, n, "file", true)
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

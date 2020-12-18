package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/OnFinality-io/onf-cli/pkg/base"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
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
	accessKey string
	secretKey string
}

func New(accessKey, secretKey string) *Api {
	baseURL := base.BaseUrl()
	req := gorequest.New()
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-onf-client", viper.GetString("app.name"))
	req.Header.Set("x-onf-version", viper.GetString("app.version"))
	return &Api{
		baseURL:   baseURL,
		req:       req,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

type RequestOptions struct {
	Params map[string]string
	Body   interface{}
}

func (a *Api) Request(method Method, path string, opts *RequestOptions) *gorequest.SuperAgent {
	r := a.req.Clone()

	u, _ := url.Parse(fmt.Sprintf("%s%s", a.baseURL, path))
	m := string(method)

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

	data := strings.Join([]string{
		m,
		r.Header.Get("content-type"),
		r.Header.Get("content-md5"),
		r.Header.Get("date"),
		canonicalize(r.Header, "x-onf-"),
		u.RequestURI(),
	}, "\n")
	signature := sign(data, a.secretKey)
	r.Header.Set("authorization", fmt.Sprintf("ONF %s:%s", a.accessKey, signature))

	return r.CustomMethod(m, u.String())
}

func (a *Api) Upload(path string, opts *RequestOptions) *gorequest.SuperAgent {
	// overwrite content-type
	contentTypeKey := "content-type"
	originalContentType := a.req.Header.Get(contentTypeKey)
	a.req.Header.Set(contentTypeKey, "multipart/form-data; boundary=WebAppBoundary")
	defer a.req.Header.Set(contentTypeKey, originalContentType)
	return a.Request(MethodPost, path, opts)
}

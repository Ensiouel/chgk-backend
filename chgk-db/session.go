package chgkdb

import (
	"net/http"
	"net/url"
	"time"
)

type Session struct {
	client *http.Client
}

func NewSession() *Session {
	return &Session{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (session *Session) Call(method string, params url.Values) *Request {
	if params == nil {
		params = url.Values{}
	}
	return &Request{
		session: session,
		params:  params,
		method:  method,
	}
}

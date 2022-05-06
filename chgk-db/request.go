package chgkdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Request struct {
	session *Session
	params  url.Values
	method  string
}

func (request *Request) Set(key string, value interface{}) *Request {
	request.params.Set(key, fmt.Sprintf("%v", value))
	return request
}

func (request *Request) Execute(v interface{}) error {
	resp, err := request.session.client.Get(fmt.Sprintf(APIAddress, request.method, request.params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &v)
}

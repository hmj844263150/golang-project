package rpc

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
)

type Response struct {
	Ctx     context.Context
	Nonce   int64
	Status  int
	Message string
	Body    map[string]interface{}
	Err     error

	W http.ResponseWriter
}

func NewResponse() *Response {
	r := &Response{Status: 200}
	return r
}

func (r *Response) Reset() {
	r.Ctx = nil
	r.Nonce = 0
	r.Status = 200
	r.Message = ""
	r.Body = nil
	r.Err = nil
}

func (r *Response) With(status int, message string) {
	r.Status, r.Message = status, message
}

func (r *Response) ValueWithErr(key string, value interface{}, err error) {
	if err != nil {
		r.Err = err
		return
	}
	r.Body[key] = value
}

var internalServerError = []byte(`{"status": 500, "result": "failed", "message": "Internal Server Error"}`)

func (r *Response) Json() []byte {
	if r.Err != nil {
		m := make(map[string]interface{})
		if r.Status == 200 {
			r.Status = 500
		}
		m["status"] = r.Status
		if r.Nonce != 0 {
			m["nonce"] = r.Nonce
		}
		m["message"] = r.Err.Error()
		bytes, _ := r.Marshal(m)
		return bytes
	}
	r.Body["status"] = r.Status
	if r.Nonce != 0 {
		r.Body["nonce"] = r.Nonce
	}
	if r.Message != "" {
		r.Body["message"] = r.Message
	}
	bytes, err := r.Marshal(r.Body)
	if err == nil {
		return bytes
	}
	return internalServerError
}

func (r *Response) Marshal(v interface{}) ([]byte, error) {
	bs, err := json.MarshalIndent(v, "", "")
	if err != nil {
		return nil, err
	}
	mark := 0
	count := 0
	for i, x := range bs {
		if x != '\n' {
			continue
		}
		if i > 0 && bs[i-1] == ',' {
			bs[i] = ' '
			continue
		}
		copy(bs[mark-count:i-count], bs[mark:i])
		count++
		mark = i + 1
	}
	if mark >= count {
		copy(bs[mark-count:], bs[mark:])
	}
	bs = bs[:len(bs)-count]
	return bs, nil
}

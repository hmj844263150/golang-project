package rpc

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func R2R(req *Request, r *http.Request) error {
	var err error

	req.Path = r.URL.Path
	req.SetMethod(r.Method)

	unmarshalGet(req, r)
	unmarshalMeta(req, r)
	unmarshalCookies(req, r)

	if req.Method != Get && req.Get[".ignore_body"] == nil {
		err = unmarshalBody(req, r)
		if err != nil {
			return err
		}
	}

	return nil
}

const maxPlus1 = 1024*1024 + 1

func unmarshalBody(req *Request, r *http.Request) error {
	length := maxPlus1
	if r.ContentLength != 0 && int(r.ContentLength) < length {
		length = int(r.ContentLength)
	}
	defer r.Body.Close()
	bytes := make([]byte, length)
	read, err := io.ReadFull(r.Body, bytes)
	if err != nil && err != io.EOF {
		return err
	}
	if read == 0 {
		return nil
	}
	if read == maxPlus1 {
		return errors.New("body length exceed limit")
	}
	body := make(map[string]interface{})
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Body = body
	req.Nonce, _ = req.GetInt("nonce", req.Body)
	return nil
}

func unmarshalGet(req *Request, r *http.Request) {
	get := make(map[string]interface{})
	values := r.URL.Query()
	for k, v := range values {
		get[k] = v[0]
		switch k {
		case "action":
			req.Action = v[0]
		case "method":
			req.SetMethod(v[0])
		}
	}
	req.Get = get
}

func unmarshalMeta(req *Request, r *http.Request) {
	meta := make(map[string]interface{})
	for k, vs := range r.Header {
		v := vs[0]
		switch k {
		case "token", "Token":
			req.Token = v
			meta[k] = v
		case "Time-Zone":
			loc, err := LoadLocation(v)
			if err == nil {
				req.Tzname = v
				req.Loc = loc
			}
			meta[k] = v
		case "Mode":
			req.SetMode(v)
			meta[k] = v
		default:
			continue
		}
	}
	req.Meta = meta
}

func unmarshalCookies(req *Request, r *http.Request) {
	if req.Token != "" {
		return
	}
	c, err := r.Cookie("Token")
	if err != nil {
		log.Println(err)
		return
	}
	req.Token = c.Value
}

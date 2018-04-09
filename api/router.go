package api

import (
	"encoding/json"
	"errors"
	"espressif.com/chip/factory/rpc"
	"fmt"
	"log"
	"strings"
)

var Router = &router{Root: &Seg{}}

type router struct {
	Root *Seg
}

func (r *router) Lookup(path string) (*Seg, map[string]interface{}) {
	if path != "/" && strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	ss := strings.Split(path, "/")[1:]
	segs := []*Seg{r.Root}
	last := len(ss) - 1
	var segMatch *Seg
	for i, part := range ss {
		segsLookup := []*Seg{}
		for _, seg := range segs {
			segsLookup = append(segsLookup, seg.Lookup(part)...)
		}
		if i == last {
			for _, seg := range segsLookup {
				if seg.Api != nil {
					segMatch = seg
					break
				}
			}
		}
		segs = segsLookup
	}
	if segMatch == nil {
		return nil, nil
	}
	args := make(map[string]interface{})
	i := last
	child := segMatch
	for child != nil {
		if child.IsPlaceholder {
			value, err := ptypeFunc[child.Ptype](ss[i])
			if err != nil {
				return nil, nil
			}
			args[child.Name] = value
		}
		i = i - 1
		child = child.Parent
	}
	return segMatch, args
}

func (r *router) Build() error {
	for _, a := range apiList {
		if strings.HasPrefix(a.Pattern, "/") {
			a.Pattern = a.Pattern[1:]
		}
		if strings.HasSuffix(a.Pattern, "/") {
			a.Pattern = a.Pattern[0 : len(a.Pattern)-1]
		}
		parent := r.Root
		ss := strings.Split(a.Pattern, "/")
		for _, s := range ss {
			subs := strings.Split(s, ":")
			if len(subs) == 1 {
				child := &Seg{Parent: parent, Name: s, IsPlaceholder: false}
				find := parent.find(child)
				if find != nil {
					parent = find
					continue
				}
				parent.Children = append(parent.Children, child)
				parent = child
				continue
			} else if len(subs) == 2 {
				ptype, ok := ptypeMap[subs[1]]
				if !ok {
					panic(fmt.Sprintf("unknow ptype: %s", subs[1]))
				}
				child := &Seg{Parent: parent, Name: subs[0], IsPlaceholder: true, Ptype: ptype}
				find := parent.find(child)
				if find != nil {
					parent = find
					continue
				}
				parent.Children = append(parent.Children, child)
				parent = child
			} else {
				panic("pattern not match")
			}
		}
		parent.Api = a
	}
	r.Root.Build()
	return nil
}

func (r *router) DispatchData(data []byte, req *rpc.Request, resp *rpc.Response) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	err := r.fillReq(data, req)
	if err != nil {
		return err
	}
	return r.Dispatch(req, resp)
}

func (r *router) Dispatch(req *rpc.Request, resp *rpc.Response) (err error) {
	seg, args := r.Lookup(req.Path)
	if seg == nil {
		return errors.New("path not exists, what are you looking for?")
	}
	req.Args = args
	resp.Body = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			switch rr := r.(type) {
			case string:
				err = errors.New(rr)
			case error:
				err = rr
			}
		}
	}()
	if seg.Api.Filter != nil {
		err := seg.Api.Filter(req, resp)
		if err != nil {
			return err
		}
	}
	seg.Api.Func(req, resp)
	return nil
}

func (r *router) fillReq(data []byte, req *rpc.Request) error {
	var err error
	req.Body = make(map[string]interface{})
	err = json.Unmarshal(data, &req.Body)
	if err != nil {
		return err
	}
	path, err := req.GetString("path", req.Body)
	if err != nil {
		return errors.New("need path")
	}
	req.Path = path
	method, _ := req.GetString("method", req.Body)
	req.SetMethod(method)

	get, err := req.GetMap("get", req.Body)
	if err == nil {
		req.Get = get
	}
	meta, err := req.GetMap("meta", req.Body)
	if err == nil {
		req.Meta = meta
	}
	authorization, err := req.GetString("Authorization", req.Meta)
	if err == nil {
		ss := strings.Split(authorization, " ")
		req.Token = ss[1]
	}
	timezone, err := req.GetString("Time-Zone", req.Meta)
	if err == nil {
		loc, err := rpc.LoadLocation(timezone)
		if err != nil {
			req.Loc = loc
			req.Tzname = timezone
		}
	}
	mode, err := req.GetString("mode", req.Meta)
	if err == nil {
		req.SetMode(mode)
	}
	action, err := req.GetString("action", req.Get)
	if err == nil {
		req.Action = action
	}
	return nil
}

func init() {
	err := Router.Build()
	if err != nil {
		panic(err)
	}
}

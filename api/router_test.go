package api

import (
	"errors"
	"espressif.com/chip/factory/rpc"
	"testing"
)

func keys(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("keys")
}

func keysId(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("keysId")
}

func products(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("products")
}

func productId(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("productId")
}

func productSerial(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("productSerial")
}

func datapointValue(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("datapointValue")
}

func datapointValueSet(req *rpc.Request, resp *rpc.Response) {
	resp.Err = errors.New("datapointValueSet")
}

func Test(t *testing.T) {
	apiList = []*Api{
		{`/keys`, keys},
		{`/keys/keyId:int`, keysId},
		{`/products`, products},
		{`/products/productId:int/`, productId},
		{`/product/productId:int/`, productId},
		{`/product/serial/productSerial:string/`, productSerial},
		{`/datapoint/value:float64`, datapointValue},
		{`/datapoint/value:float64/set`, datapointValueSet},
	}
	req, resp := &rpc.Request{}, &rpc.Response{}
	Router = &router{Root: &Seg{}}
	Router.Build()
	seg, args := Router.Lookup("/noexist")
	if seg != nil {
		t.Fatal("find /noexist?")
	}
	seg, args = Router.Lookup("/keys")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "keys" {
		t.Fatal("find /keys not match")
	}
	if len(args) > 0 {
		t.Fatal("len(args) == 0?")
	}
	seg, args = Router.Lookup("/keys/123456")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "keysId" {
		t.Fatal("func is not keysId")
		if args["keyId"].(int) != 123456 {
			t.Fatal("args not found 123456")
		}
	}
	seg, args = Router.Lookup("/keys/abcdef")
	if seg != nil {
		t.Fatal("find /keys/abcdef?")
	}
	seg, args = Router.Lookup("/keys/123456/abc")
	if seg != nil {
		t.Fatal("find /keys/123456/abc?")
	}
	seg, args = Router.Lookup("/products")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "products" {
		t.Fatal("func is not products")
	}
	seg, args = Router.Lookup("/products/12345")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "productId" {
		t.Fatal("func is not productId")
		if args["productId"].(int) != 12345 {
			t.Fatal("args not found 12345")
		}
	}
	seg, args = Router.Lookup("/product/12345")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "productId" {
		t.Fatal("func is not productId")
		if args["productId"].(int) != 12345 {
			t.Fatal("args not found 12345")
		}
	}
	seg, args = Router.Lookup("/product/serial/abcef")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "productSerial" {
		t.Fatal("func is not productSerial")
		if args["productSerial"].(string) != "abcef" {
			t.Fatal("args not found abcef")
		}
	}
	seg, args = Router.Lookup("/datapoint/123")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "datapointValue" {
		t.Fatal("func is not datapointValue")
		if args["value"].(float64) != 123.0 {
			t.Fatal("args not found 123.0")
		}
	}
	seg, args = Router.Lookup("/datapoint/12.3/set")
	seg.Api.Func(req, resp)
	if resp.Err.Error() != "datapointValueSet" {
		t.Fatal("func is not datapointValueSet")
		if args["value"].(float64) != 12.3 {
			t.Fatal("args not found 12.3")
		}
	}
}

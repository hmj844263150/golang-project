package mpn

import (
	"espressif.com/chip/factory/config"
	"espressif.com/chip/factory/rpc"
	"fmt"
	"io/ioutil"
	"os"
)

func Mpn(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		mpnConfigGet(req, resp)
	case rpc.Post:
		mpnConfigPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func mpnConfigGet(req *rpc.Request, resp *rpc.Response) {
	mpnSid, err := req.GetString("mpnSid", req.Get)
	if err != nil {
		resp.Err = rpc.BadRequest
		return
	}
	mpnConfigDir := config.Cfg.MpnConfigDir + "/" + mpnSid
	_, err = os.Stat(mpnConfigDir)

	if err != nil {
		resp.Err = rpc.NotFound
		return
	}
	f, err := os.Open(mpnConfigDir)
	if err != nil {
		resp.Err = rpc.NotFound
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		resp.Err = rpc.NotFound
		return
	}
	resp.Body["data"] = string(data)
}

func mpnConfigPost(req *rpc.Request, resp *rpc.Response) {
	mpnSid, err := req.GetString("mpnSid", req.Body)
	data, err := req.GetString("data", req.Body)
	if err != nil {
		fmt.Println("get data err")
		resp.Err = rpc.BadRequest
		return
	}
	mpnConfigDir := config.Cfg.MpnConfigDir + "/" + mpnSid

	fmt.Println(mpnConfigDir)
	f, err := os.Create(mpnConfigDir)

	if err != nil {
		resp.Err = rpc.BadRequest
		return
	}

	defer f.Close()
	_, err = f.Write([]byte(data))
	if err != nil {
		resp.Err = rpc.BadRequest
		return
	}
}

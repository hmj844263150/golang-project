package hengha

import (
	"espressif.com/chip/factory/rpc"
)

func Auth(req *rpc.Request, resp *rpc.Response) error {
	req.Auth(false)
	return nil
}

func AuthStaff(req *rpc.Request, resp *rpc.Response) error {
	req.Auth(true)
	return nil
}

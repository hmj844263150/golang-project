package hengha

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Createdb(req *rpc.Request, resp *rpc.Response) {
	dal.Createdb()
}

func Dropdb(req *rpc.Request, resp *rpc.Response) {
	dal.Dropdb()
}

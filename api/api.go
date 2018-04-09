package api

import (
	"espressif.com/chip/factory/rpc"
)

type Api struct {
	Pattern string
	Func    func(req *rpc.Request, resp *rpc.Response)
	Filter  func(req *rpc.Request, resp *rpc.Response) error
}

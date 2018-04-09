package factory

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Factory(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Post:
		factoryPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func factoryPost(req *rpc.Request, resp *rpc.Response) {
	factory, err := dal.UnmarshalFactory(req.Ctx, req.GetValue("factory", req.Body))
	if err != nil {
		resp.Err = err
		return
	}
	exist := dal.FindFactoryByName(req.Ctx, factory.Name)
	if exist != nil {
		resp.Err = rpc.Conflict
		return
	}
	err = factory.Generate()
	if err != nil {
		resp.Err = err
		return
	}
	err = factory.Save()
	if err != nil {
		resp.Err = err
		return
	}
	resp.Body["factory"] = factory
}

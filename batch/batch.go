package batch

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Batch(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Post:
		batchPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func batchPost(req *rpc.Request, resp *rpc.Response) {
	batch, err := dal.UnmarshalBatch(req.Ctx, req.GetValue("batch", req.Body))
	if err != nil {
		resp.Err = err
		return
	}
	if req.Factory != nil && !req.Factory.IsStaff {
		batch.FactorySid = req.Factory.Sid
	}
	find := dal.FindFactoryBySid(req.Ctx, batch.FactorySid)
	if find == nil {
		resp.Err = rpc.NotFound
		return
	}
	exist := dal.FindBatchBySid(req.Ctx, batch.Sid)
	if exist != nil {
		resp.Err = rpc.Conflict
		return
	}
	err = batch.Adjust()
	if err != nil {
		resp.Err = err
		return
	}
	err = batch.Save()
	if err != nil {
		resp.Err = err
		return
	}
	resp.Body["batch"] = batch
}

package batch

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Batchs(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		batchsGet(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func batchsGet(req *rpc.Request, resp *rpc.Response) {
	batchs := listBatch(req)
	resp.Body["batchs"] = batchs
}

func listBatch(req *rpc.Request) []*dal.Batch {
	req.FillRange(true, false)
	factorySid, err := req.GetString("factory_sid", req.Get)
	if err == nil {
		return dal.ListBatchByFactorySid(req.Ctx, factorySid, req.Offset, req.RowCount)
	}
	if req.Factory != nil {
		if req.Factory.IsStaff {
			return dal.ListBatchAll(req.Ctx, 0, 1000)
		}
		return dal.ListBatchByFactorySid(req.Ctx, req.Factory.Sid, req.Offset, req.RowCount)
	}
	return nil
}

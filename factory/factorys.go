package factory

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Factorys(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		factorysGet(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func factorysGet(req *rpc.Request, resp *rpc.Response) {
	req.FillRange(true, false)
	factorys := listFactory(req)
	withBatch, _ := req.GetBool("with_batch", req.Get)
	if withBatch {
		for _, factory := range factorys {
			batchs := dal.ListBatchByFactorySid(req.Ctx, factory.Sid, 0, 1000)
			factory.Padding("batchs", batchs)
		}
	}
	resp.Body["factorys"] = factorys
}

func listFactory(req *rpc.Request) []*dal.Factory {
	if req.Factory != nil {
		if req.Factory.IsStaff {
			return dal.ListFactoryAll(req.Ctx, req.Offset, req.RowCount)
		}
		return []*dal.Factory{req.Factory}
	}
	return nil
}

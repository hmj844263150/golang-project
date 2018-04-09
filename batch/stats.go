package batch

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
	"strings"
)

func Stats(req *rpc.Request, resp *rpc.Response) {
	str, err := req.GetString("batch_sids", req.Get)
	if err == nil {
		statsBatchs(req, resp, str)
		return
	}
	batchSid, err := req.GetString("batch_sid", req.Get)
	if err == nil {
		statsBatch(req, resp, batchSid)
		return
	}
	batchSidD, err := req.GetString("batch_detail", req.Get)
	if err == nil {
		statsBatchDetail(req, resp, batchSidD)
		return
	}
}

func statsBatchs(req *rpc.Request, resp *rpc.Response, str string) {
	batchSids := strings.Split(str, ",")
	if len(batchSids) == 0 {
		resp.Err = rpc.BadRequest
		return
	}
	batchs := make([]*dal.Batch, len(batchSids))
	for i, batchSid := range batchSids {
		batch := dal.FindBatchBySid(req.Ctx, batchSid)
		if batch == nil {
			resp.Err = rpc.BadRequest
			return
		}
		if req.Factory != nil && !req.Factory.IsStaff && req.Factory.Sid != batch.FactorySid {
			resp.Err = rpc.Forbidden
			return
		}
		batchs[i] = batch
		if batch.Success == batch.Cnt {
			continue
		}
		testdatas := dal.ListTestdataByBatch(req.Ctx, batch.Sid, 0, 1)
		if len(testdatas) == 0 {
			continue
		}
		if batch.Statsed == testdatas[0].Created {
			continue
		}
		batch.Success = dal.CountTestdataByBatchSuccess(req.Ctx, batch.Sid)
		batch.RightFirstTime = dal.CountTestdataByBatchRightFirstTime(req.Ctx, batch.Sid, batch.Sid)
		batch.Failed = dal.CountTestdataByBatchFailed(req.Ctx, batch.Sid)
		batch.Rejected = dal.CountTestdataByBatchRejected(req.Ctx, batch.Sid)
		batch.Statsed = testdatas[0].Created
		batch.Update(dal.BatchCol.Success, dal.BatchCol.RightFirstTime, dal.BatchCol.Failed, dal.BatchCol.Rejected, dal.BatchCol.Statsed)
	}
	resp.Body["batchs"] = batchs
}

func statsBatch(req *rpc.Request, resp *rpc.Response, batchSid string) {
	batch := dal.FindBatchBySid(req.Ctx, batchSid)
	if batch == nil {
		resp.Err = rpc.BadRequest
		return
	}
	resp.Body["days"], resp.Body["weeks"] = dal.ListBatchStatsByBatchIdRange(req.Ctx, batch.Id)
}

func statsBatchDetail(req *rpc.Request, resp *rpc.Response, batchSid string) {
	batch := dal.FindBatchBySid(req.Ctx, batchSid)
	if batch == nil {
		resp.Err = rpc.BadRequest
		return
	}
	//resp.Body["print_pass"] = dal.CountTestdataByBatchPrintPass(req.Ctx, batch.Sid)
	if batch.Success != batch.Cnt {
		testdatas := dal.ListTestdataByBatch(req.Ctx, batch.Sid, 0, 1)
		if len(testdatas) > 0 && batch.Statsed != testdatas[0].Created {
			batch.Success = dal.CountTestdataByBatchSuccess(req.Ctx, batch.Sid)
			//batch.RightFirstTime = dal.CountTestdataByBatchRightFirstTime(req.Ctx, batch.Sid, batch.Sid)
			batch.Failed = dal.CountTestdataByBatchFailed(req.Ctx, batch.Sid)
			batch.Rejected = dal.CountTestdataByBatchRejected(req.Ctx, batch.Sid)
			batch.Update(dal.BatchCol.Success, dal.BatchCol.RightFirstTime, dal.BatchCol.Failed, dal.BatchCol.Rejected, dal.BatchCol.Statsed)
		}
	}

	resp.Body["batch"] = batch
}

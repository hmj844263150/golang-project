package testdata

import (
	"errors"
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Stats(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		statsGet(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func statsGet(req *rpc.Request, resp *rpc.Response) {
	by, err := req.GetString("by", req.Get)
	if err != nil {
		resp.Err = err
		return
	}
	switch by {
	case "batch":
		batchStats(req, resp)
	case "datetime":
		datetimeStats(req, resp)
	case "datetime_batch":
		datetimeBatchStats(req, resp)
	case "mac":
		macStats(req, resp)
	default:
		resp.Err = errors.New("by parameter Not Found")
		return
	}
}

func batchStats(req *rpc.Request, resp *rpc.Response) {
	batch, err := req.GetString("batch_sid", req.Get)
	if err != nil {
		resp.Err = rpc.BadRequest
		return
	}
	successCnt := dal.CountTestdataByBatch(req.Ctx, batch, "success")
	failedCnt := dal.CountTestdataByBatch(req.Ctx, batch, "failed")
	resp.Body["batch"] = batch
	respStats(resp, successCnt, failedCnt)
}

func datetimeStats(req *rpc.Request, resp *rpc.Response) {
	req.FillRange(false, true)
	successCnt := dal.CountTestdataByDatetime(req.Ctx, req.Start, req.End, "success")
	failedCnt := dal.CountTestdataByDatetime(req.Ctx, req.Start, req.End, "failed")
	resp.Body["start"] = req.Start
	resp.Body["end"] = req.End
	respStats(resp, successCnt, failedCnt)
}

func datetimeBatchStats(req *rpc.Request, resp *rpc.Response) {
	req.FillRange(false, true)
	batch, err := req.GetString("batch_sid", req.Get)
	if err != nil {
		resp.Err = rpc.BadRequest
		return
	}
	successCnt := dal.CountTestdataByDatetimeBatch(req.Ctx, req.Start, req.End, batch, "success")
	failedCnt := dal.CountTestdataByDatetimeBatch(req.Ctx, req.Start, req.End, batch, "failed")
	resp.Body["start"] = req.Start
	resp.Body["end"] = req.End
	respStats(resp, successCnt, failedCnt)
}

func macStats(req *rpc.Request, resp *rpc.Response) {
	mac, err := req.GetString("esp_mac", req.Get)
	if err != nil {
		resp.Err = rpc.BadRequest
		return
	}
	successCnt := dal.CountTestdataByEspMac(req.Ctx, mac, "success")
	failedCnt := dal.CountTestdataByEspMac(req.Ctx, mac, "failed")
	resp.Body["mac"] = mac
	respStats(resp, successCnt, failedCnt)
}

func respStats(resp *rpc.Response, successCnt int, failedCnt int) {
	total := successCnt + failedCnt
	resp.Body["total"] = total
	resp.Body["success"] = successCnt
	resp.Body["failed"] = failedCnt
}

package testdata

import (
	"errors"
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
	"strings"
)

func Print(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Post:
		printPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func printPost(req *rpc.Request, resp *rpc.Response) {
	espMac, _ := req.GetString("esp_mac", req.Get)
	espMac = strings.ToLower(espMac)
	espMac = strings.Replace(espMac, ":", "", -1)
	testdata := dal.FindTestdataByEspMac(req.Ctx, espMac)

	if testdata == nil {
		resp.Err = rpc.NotFound
		return
	}
	if testdata.TestResult != "success" {
		resp.Err = errors.New("testdata.TestResult must success")
		return
	}
	batchSid := queryBatchSid(req)
	if testdata.BatchSid != batchSid {
		resp.Err = errors.New("batch_sid not match")
		return
	}
	testdata.QueryTimes++
	dryrun, _ := req.GetBool("dryrun", req.Get)
	if !dryrun {
		testdata.PrintTimes++
	}
	testdata.Update(dal.TestdataCol.QueryTimes, dal.TestdataCol.PrintTimes)
	resp.Body["test_result"] = testdata.TestResult
	resp.Body["esp_mac"] = testdata.EspMac
	resp.Body["cus_mac"] = testdata.CusMac
	resp.Body["print_times"] = testdata.PrintTimes
	resp.Body["testdata"] = testdata
}

func queryBatchSid(req *rpc.Request) string {
	batchSid, _ := req.GetString("batch_sid", req.Get)
	if batchSid != "" {
		return batchSid
	}
	batchSid, _ = req.GetString("batch_sid", req.Body)
	if batchSid != "" {
		return batchSid
	}
	v := req.GetValue("body.batch_sid", req.Body)
	if v != nil {
		batchSid, _ = v.(string)
	}
	return batchSid
}

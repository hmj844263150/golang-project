package testdata

import (
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
)

func Logs(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		logsGet(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func logsGet(req *rpc.Request, resp *rpc.Response) {
	req.FillRange(true, false)
	espMac, err := req.GetString("esp_mac", req.Get)
	if err != nil {
		resp.Err = err
		return
	}
	testdatas := dal.ListTestdataByEspMac(req.Ctx, espMac, req.Offset, req.RowCount)
	resp.Body["testdatas"] = testdatas
}

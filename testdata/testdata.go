package testdata

import (
	"errors"
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
	"strings"
	"sync"
	"time"
)

var mutexTestdata sync.Mutex

func Testdata(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Post:
		testdataPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func testdataPost(req *rpc.Request, resp *rpc.Response) {

	v := req.GetValue("testdata", req.Body)
	chkRepeatFlg, _ := req.GetBool("chk_repeat_flg", v.(map[string]interface{}))
	poType, _ := req.GetInt("po_type", v.(map[string]interface{}))
	testdata, err := dal.UnmarshalTestdata(req.Ctx, v)
	if err != nil {
		resp.Err = err
		return
	}

	var batch *dal.Batch
	if testdata.BatchSid != "" {
		batch = dal.FindBatchBySid(req.Ctx, testdata.BatchSid)
		if batch == nil {
			resp.Err = rpc.NotFound
			return
		}
	}
	if testdata.FactorySid != batch.FactorySid {
		resp.Err = errors.New("Factory Conflict")
		return
	}

	testdata.EspMac = strings.ToLower(testdata.EspMac)
	testdata.EspMac = strings.Replace(strings.Replace(testdata.EspMac, ":", "", -1), "-", "", -1)
	testdata.CusMac = strings.ToLower(testdata.CusMac)
	testdata.CusMac = strings.Replace(strings.Replace(testdata.CusMac, ":", "", -1), "-", "", -1)

	if len(testdata.EspMac) != 12 || (len(testdata.CusMac) > 0 && len(testdata.CusMac) != 12) {
		resp.Err = rpc.BadRequest
		return
	}

	// cus_mac equal esp_mac ?
	if len(testdata.CusMac) == 12 {
		if dal.FindTestdataByEspMac(req.Ctx, testdata.CusMac) != nil {
			resp.Err = errors.New("412 Cus Mac Conflict")
			return
		}
	}

	resp.Body["test_pass_record"] = "0"

	origin := dal.FindTestdataByEspMac(req.Ctx, testdata.EspMac)
	if len(testdata.CusMac) == 12 {
		origin = dal.FindTestdataByCusMac(req.Ctx, testdata.CusMac)
	}
	if origin != nil {
		// same cus_mac but diff esp_mac ?
		if len(testdata.CusMac) == 12 && testdata.EspMac != origin.EspMac {
			resp.Err = errors.New("411 Cus Mac Repeat")
			return
		}

		testPassRecord := dal.FindTestdataByEspMacSuccess(req.Ctx, testdata.EspMac)
		if len(testdata.CusMac) == 12 {
			testPassRecord = dal.FindTestdataByCusMacSuccess(req.Ctx, testdata.CusMac)
		}

		if testPassRecord != nil {
			resp.Body["test_pass_record"] = "1"
		}

		// check for repeat mac
		if chkRepeatFlg && testPassRecord != nil {
			if time.Now().Unix()-testPassRecord.Updated.Unix() > 120 {
				resp.Err = errors.New("410 Mac Repeat timeout")
				return
			}
		}

		// commit and already apply cus mac
		if testdata.IsCommit && origin.TestResult == "success" && origin.BatchIndex != 0 && origin.CusMac != "" {
			origin.IsCommit = true
			testdata = origin
		} else {
			if (poType == 0) && (origin.FactorySid != testdata.FactorySid || origin.BatchSid != testdata.BatchSid) {
				resp.Err = rpc.Conflict
				return
			}
			testdata.QueryTimes = origin.QueryTimes
			testdata.PrintTimes = origin.PrintTimes
			//testdata.CusMac = origin.CusMac

			oldBatchData := dal.FindTestdataByBatchSidEspMac(req.Ctx, testdata.BatchSid, testdata.EspMac)
			if len(testdata.CusMac) == 12 {
				oldBatchData = dal.FindTestdataByBatchSidCusMac(req.Ctx, testdata.BatchSid, testdata.CusMac)
			}
			if oldBatchData != nil {
				testdata.BatchIndex = oldBatchData.BatchIndex
			} else {
				testdata.BatchIndex = 0
			}

			origin.Latest = false
			origin.Update(dal.TestdataCol.Latest)
		}
	}

	mutexTestdata.Lock()
	defer mutexTestdata.Unlock()
	if testdata.BatchSid != "" {
		batch = dal.FindBatchBySid(req.Ctx, testdata.BatchSid)
		if batch == nil {
			resp.Err = rpc.NotFound
			return
		}
	}
	if batch != nil && testdata.TestResult == "success" && testdata.BatchIndex == 0 {
		//testdata.BatchIndex, testdata.CusMac, err = batch.NextRemainMac()
		testdata.BatchIndex, _, err = batch.NextRemainMac()
		if err != nil {
			resp.Err = err
			return
		}
	}

	testdata.Save()
	resp.Body["test_result"] = testdata.TestResult
	resp.Body["esp_mac"] = testdata.EspMac
	resp.Body["cus_mac"] = testdata.CusMac
	resp.Body["batch_index"] = testdata.BatchIndex
	resp.Body["testdata"] = testdata
	if batch != nil {
		resp.Body["batch_cnt"] = batch.Cnt
		resp.Body["batch_remain"] = batch.Remain
	}
}

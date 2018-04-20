package dal

import (
	"bytes"
	"fmt"
	"context"
	"espressif.com/chip/factory/db"
	"strconv"
	"strings"
	"time"
)

var TestdataTid = 1
var _ db.Doer = (*Testdata)(nil)
var testdatacols = []db.Col{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
var testdatafields = []string{"id", "created", "updated", "visibly", "module_id", "device_type", "fw_ver", "esp_mac", "cus_mac", "flash_id", "test_result", "test_msg", "factory_sid", "batch_sid", "efuse", "query_times", "print_times", "batch_index", "latest", "is_commit"}

var TestdataCol = struct {
	Id, Created, Updated, Visibly, ModuleId, DeviceType, FwVer, EspMac, CusMac, FlashId, TestResult, TestMsg, FactorySid, BatchSid, Efuse, QueryTimes, PrintTimes, BatchIndex, Latest, IsCommit, _ db.Col
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 0}

type Testdata struct {
	Id         int
	Created    time.Time
	Updated    time.Time
	Visibly    bool
	ModuleId   int
	DeviceType string
	FwVer      string
	EspMac     string
	CusMac     string
	FlashId    string
	TestResult string
	TestMsg    string
	FactorySid string
	BatchSid   string
	Efuse      string
	QueryTimes int
	PrintTimes int
	BatchIndex int
	Latest     bool
	IsCommit   bool

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewTestdata(ctx context.Context) *Testdata {
	now := time.Now()
	t := &Testdata{Created: now, Updated: now, Visibly: true}
	t.ext = GetExtFromContext(ctx)
	defaultTestdata(ctx, t)
	return t
}

func FindTestdata(ctx context.Context, id int) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func ListTestdata(ctx context.Context, ids ...int) []*Testdata {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(testdataSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testdatas := make([]*Testdata, len(dos))
	for ii, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		testdatas[ii] = testdata
	}
	return testdatas
}

func FindTestdataByEspMac(ctx context.Context, espMac string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[8], espMac)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func ListTestdataByEspMac(ctx context.Context, espMac string, offset int, rowCount int) []*Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[9], espMac, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testdatas := make([]*Testdata, len(dos))
	for ii, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		testdatas[ii] = testdata
	}
	return testdatas
}

func CountTestdataByEspMac(ctx context.Context, espMac string, testResult string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[10], espMac, testResult)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatch(ctx context.Context, batchSid string, testResult string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[11], batchSid, testResult)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByDatetime(ctx context.Context, start time.Time, end time.Time, testResult string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[12], start, end, testResult)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByDatetimeBatch(ctx context.Context, start time.Time, end time.Time, batchSid string, testResult string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[13], start, end, batchSid, testResult)
	if err != nil {
		return -1
	}
	return int(rows)
}

func ListTestdataByBatch(ctx context.Context, batchSid string, offset int, rowCount int) []*Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[14], batchSid, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testdatas := make([]*Testdata, len(dos))
	for ii, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		testdatas[ii] = testdata
	}
	return testdatas
}

func ListTestdataAll(ctx context.Context, offset int, rowCount int) []*Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[15], offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testdatas := make([]*Testdata, len(dos))
	for ii, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		testdatas[ii] = testdata
	}
	return testdatas
}

func ListTestdataByFactory(ctx context.Context, factorySid string, offset int, rowCount int) []*Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[16], factorySid, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testdatas := make([]*Testdata, len(dos))
	for ii, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		testdatas[ii] = testdata
	}
	return testdatas
}

func CountTestdataByBatchSuccess(ctx context.Context, batchSid string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[17], batchSid)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchRightFirstTime(ctx context.Context, batchSid1 string, batchSid2 string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[18], batchSid1, batchSid2)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchFailed(ctx context.Context, batchSid string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[19], batchSid)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchRejected(ctx context.Context, batchSid string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[20], batchSid)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchStartEnd(ctx context.Context, batchSid string, start time.Time, end time.Time) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[21], batchSid, start, end)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchSuccessStartEnd(ctx context.Context, batchSid string, start time.Time, end time.Time) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[22], batchSid, start, end)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchRightFirstTimeStartEnd(ctx context.Context, batchSid1 string, start1 time.Time, end1 time.Time, batchSid2 string, start2 time.Time, end2 time.Time) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[23], batchSid1, start1, end1, batchSid2, start2, end2)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchFailedStartEnd(ctx context.Context, batchSid string, start time.Time, end time.Time) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[24], batchSid, start, end)
	if err != nil {
		return -1
	}
	return int(rows)
}

func CountTestdataByBatchRejectedStartEnd(ctx context.Context, batchSid string, start time.Time, end time.Time) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[25], batchSid, start, end)
	if err != nil {
		return -1
	}
	return int(rows)
}

func FindTestdataByCusMac(ctx context.Context, cusMac string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[26], cusMac)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func FindTestdataByBatchSidEspMac(ctx context.Context, batchSid string, espMac string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[27], batchSid, espMac)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func FindTestdataByBatchSidCusMac(ctx context.Context, batchSid string, cusMac string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[28], batchSid, cusMac)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func FindTestdataByEspMacSuccess(ctx context.Context, espMac string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[29], espMac)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func FindTestdataByCusMacSuccess(ctx context.Context, cusMac string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[30], cusMac)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func CountTestdataByBatchPrintPass(ctx context.Context, batchSid string) int {
	rows, err := db.Open("Testdata").Count(testdataSqls[31], batchSid)
	if err != nil {
		return -1
	}
	return int(rows)
}

func ListTestdataByFactoryEspMac(ctx context.Context, factorySid string, espMac string, offset int, rowCount int) []*Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[32], factorySid, espMac, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testdatas := make([]*Testdata, len(dos))
	for ii, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		testdatas[ii] = testdata
	}
	return testdatas
}

func FindTestdataByBatchSidNewst(ctx context.Context, batchSid string) *Testdata {
	dos, err := db.Open("Testdata").Query(newTestdataDest, true, testdataSqls[33], batchSid)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testdata, _ := do.(*Testdata)
		if ext != nil {
			testdata.ext = ext
		}
		return testdata
	}
	return nil
}

func (t *Testdata) Save() error {
	now := time.Now()
	t.Created, t.Updated, t.Visibly = now, now, true
	var id int64
	var err error
	if t.Id == 0 {
		id, _, err = db.Open("Testdata").Exec(testdataSqls[0], t.ModuleId, t.DeviceType, t.FwVer, t.EspMac, t.CusMac, t.FlashId, t.TestResult, t.TestMsg, t.FactorySid, t.BatchSid, t.Efuse, t.QueryTimes, t.PrintTimes, t.BatchIndex, t.Latest, t.IsCommit)
	} else {
		id, _, err = db.Open("Testdata").Exec(testdataSqls[1], t.Id, t.ModuleId, t.DeviceType, t.FwVer, t.EspMac, t.CusMac, t.FlashId, t.TestResult, t.TestMsg, t.FactorySid, t.BatchSid, t.Efuse, t.QueryTimes, t.PrintTimes, t.BatchIndex, t.Latest, t.IsCommit)
	}
	if err != nil {
		return err
	}
	t.Id = int(id)
	return nil
}

func (t *Testdata) Update(cs ...db.Col) error {
	if t.Id == 0 {
		return logError("dal.Testdata Error: can not update row while id is zero")
	}
	t.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("Testdata").Exec(testdataSqls[2], t.Visibly, t.ModuleId, t.DeviceType, t.FwVer, t.EspMac, t.CusMac, t.FlashId, t.TestResult, t.TestMsg, t.FactorySid, t.BatchSid, t.Efuse, t.QueryTimes, t.PrintTimes, t.BatchIndex, t.Latest, t.IsCommit, t.Id)
		return err
	}
	cols, args, err := colsAndArgsTestdata(t, cs...)
	if err != nil {
		return err
	}
	args = append(args, t.Id)
	sqlstr := fmt.Sprintf(testdataSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("Testdata").Exec(sqlstr, args...)
	return err
}

func (t *Testdata) Invisibly() error {
	if t.Id == 0 {
		return logError("dal.Testdata Error: can not invisibly row while id is zero")
	}
	t.Updated = time.Now()
	t.Visibly = false
	_, _, err := db.Open("Testdata").Exec(testdataSqls[4], t.Id)
	return err
}

func (t *Testdata) Delete() error {
	if t.Id == 0 {
		return logError("dal.Testdata Error: can not delete row while id is zero")
	}
	t.Updated = time.Now()
	_, _, err := db.Open("Testdata").Exec(testdataSqls[5], t.Id)
	return err
}

func (t *Testdata) Valid() error {
	return t.valid()
}

func (t *Testdata) SetExt(ext *Ext) {
	t.ext = ext
}

func (t *Testdata) Padding(pkey string, pvalue interface{}) {
	if t.ext == nil {
		t.ext = &Ext{Loc: DefaultLoc}
	}
	if t.paddings == nil {
		t.paddings = make(map[string]interface{})
	}
	t.paddings[pkey] = pvalue
	t.ext.IsComplex = true
}

func (t *Testdata) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case TestdataCol.Id:
			if isColumnName {
				mm["id"] = t.Id
			} else {
				mm["Id"] = t.Id
			}
		case TestdataCol.Created:
			if isColumnName {
				mm["created"] = t.Created
			} else {
				mm["Created"] = t.Created
			}
		case TestdataCol.Updated:
			if isColumnName {
				mm["updated"] = t.Updated
			} else {
				mm["Updated"] = t.Updated
			}
		case TestdataCol.Visibly:
			if isColumnName {
				mm["visibly"] = t.Visibly
			} else {
				mm["Visibly"] = t.Visibly
			}
		case TestdataCol.ModuleId:
			if isColumnName {
				mm["module_id"] = t.ModuleId
			} else {
				mm["ModuleId"] = t.ModuleId
			}
		case TestdataCol.DeviceType:
			if isColumnName {
				mm["device_type"] = t.DeviceType
			} else {
				mm["DeviceType"] = t.DeviceType
			}
		case TestdataCol.FwVer:
			if isColumnName {
				mm["fw_ver"] = t.FwVer
			} else {
				mm["FwVer"] = t.FwVer
			}
		case TestdataCol.EspMac:
			if isColumnName {
				mm["esp_mac"] = t.EspMac
			} else {
				mm["EspMac"] = t.EspMac
			}
		case TestdataCol.CusMac:
			if isColumnName {
				mm["cus_mac"] = t.CusMac
			} else {
				mm["CusMac"] = t.CusMac
			}
		case TestdataCol.FlashId:
			if isColumnName {
				mm["flash_id"] = t.FlashId
			} else {
				mm["FlashId"] = t.FlashId
			}
		case TestdataCol.TestResult:
			if isColumnName {
				mm["test_result"] = t.TestResult
			} else {
				mm["TestResult"] = t.TestResult
			}
		case TestdataCol.TestMsg:
			if isColumnName {
				mm["test_msg"] = t.TestMsg
			} else {
				mm["TestMsg"] = t.TestMsg
			}
		case TestdataCol.FactorySid:
			if isColumnName {
				mm["factory_sid"] = t.FactorySid
			} else {
				mm["FactorySid"] = t.FactorySid
			}
		case TestdataCol.BatchSid:
			if isColumnName {
				mm["batch_sid"] = t.BatchSid
			} else {
				mm["BatchSid"] = t.BatchSid
			}
		case TestdataCol.Efuse:
			if isColumnName {
				mm["efuse"] = t.Efuse
			} else {
				mm["Efuse"] = t.Efuse
			}
		case TestdataCol.QueryTimes:
			if isColumnName {
				mm["query_times"] = t.QueryTimes
			} else {
				mm["QueryTimes"] = t.QueryTimes
			}
		case TestdataCol.PrintTimes:
			if isColumnName {
				mm["print_times"] = t.PrintTimes
			} else {
				mm["PrintTimes"] = t.PrintTimes
			}
		case TestdataCol.BatchIndex:
			if isColumnName {
				mm["batch_index"] = t.BatchIndex
			} else {
				mm["BatchIndex"] = t.BatchIndex
			}
		case TestdataCol.Latest:
			if isColumnName {
				mm["latest"] = t.Latest
			} else {
				mm["Latest"] = t.Latest
			}
		case TestdataCol.IsCommit:
			if isColumnName {
				mm["is_commit"] = t.IsCommit
			} else {
				mm["IsCommit"] = t.IsCommit
			}
		default:
			logError(fmt.Sprintf("dal.Testdata Error: unknow column num %d in talbe testdata", cc))
		}
	}
	return mm
}

func (t *Testdata) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	loc := DefaultLoc
	var numericEnum bool
	if t.ext != nil {
		if t.ext.IsComplex {
			return t.marshalJSONComplex()
		}
		loc = t.ext.Loc
		numericEnum = t.ext.NumericEnum
	}
	var buf bytes.Buffer
	buf.WriteString(`{"id":`)
	buf.WriteString(strconv.FormatInt(int64(t.Id), 10))
	buf.WriteString(`, "created":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(t.Created.Unix(), 10))
	} else {
		t.Created = t.Created.In(loc)
		buf.WriteString(`"` + t.Created.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "updated":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(t.Updated.Unix(), 10))
	} else {
		t.Updated = t.Updated.In(loc)
		buf.WriteString(`"` + t.Updated.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "visibly":`)
	if numericEnum {
		if t.Visibly {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(t.Visibly))
	}
	buf.WriteString(`, "module_id":`)
	buf.WriteString(strconv.FormatInt(int64(t.ModuleId), 10))
	buf.WriteString(`, "device_type":`)
	WriteJsonString(&buf, t.DeviceType)
	buf.WriteString(`, "fw_ver":`)
	WriteJsonString(&buf, t.FwVer)
	buf.WriteString(`, "esp_mac":`)
	WriteJsonString(&buf, t.EspMac)
	buf.WriteString(`, "cus_mac":`)
	WriteJsonString(&buf, t.CusMac)
	buf.WriteString(`, "flash_id":`)
	WriteJsonString(&buf, t.FlashId)
	buf.WriteString(`, "test_result":`)
	WriteJsonString(&buf, t.TestResult)
	buf.WriteString(`, "test_msg":`)
	WriteJsonString(&buf, t.TestMsg)
	buf.WriteString(`, "factory_sid":`)
	WriteJsonString(&buf, t.FactorySid)
	buf.WriteString(`, "batch_sid":`)
	WriteJsonString(&buf, t.BatchSid)
	buf.WriteString(`, "efuse":`)
	WriteJsonString(&buf, t.Efuse)
	buf.WriteString(`, "query_times":`)
	buf.WriteString(strconv.FormatInt(int64(t.QueryTimes), 10))
	buf.WriteString(`, "print_times":`)
	buf.WriteString(strconv.FormatInt(int64(t.PrintTimes), 10))
	buf.WriteString(`, "batch_index":`)
	buf.WriteString(strconv.FormatInt(int64(t.BatchIndex), 10))
	buf.WriteString(`, "latest":`)
	if numericEnum {
		if t.Latest {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(t.Latest))
	}
	buf.WriteString(`, "is_commit":`)
	if numericEnum {
		if t.IsCommit {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(t.IsCommit))
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (t *Testdata) marshalJSONComplex() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	if t.ext == nil {
		return nil, logError("dal.Testdata Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := t.ext.Loc
	numericEnum := t.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[TestdataTid]; ok {
		if vvv, ok := vv[t.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := t.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(TestdataCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(t.Id), 10))
	}
	isRender = isRenderField(TestdataCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(t.Created.Unix(), 10))
		} else {
			t.Created = t.Created.In(loc)
			buf.WriteString(`"` + t.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(TestdataCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(t.Updated.Unix(), 10))
		} else {
			t.Updated = t.Updated.In(loc)
			buf.WriteString(`"` + t.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(TestdataCol.Visibly, "visibly", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "visibly":`)
		if numericEnum {
			if t.Visibly {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(t.Visibly))
		}
	}
	isRender = isRenderField(TestdataCol.ModuleId, "module_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "module_id":`)
		buf.WriteString(strconv.FormatInt(int64(t.ModuleId), 10))
	}
	isRender = isRenderField(TestdataCol.DeviceType, "device_type", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "device_type":`)
		WriteJsonString(&buf, t.DeviceType)
	}
	isRender = isRenderField(TestdataCol.FwVer, "fw_ver", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "fw_ver":`)
		WriteJsonString(&buf, t.FwVer)
	}
	isRender = isRenderField(TestdataCol.EspMac, "esp_mac", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac":`)
		WriteJsonString(&buf, t.EspMac)
	}
	isRender = isRenderField(TestdataCol.CusMac, "cus_mac", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cus_mac":`)
		WriteJsonString(&buf, t.CusMac)
	}
	isRender = isRenderField(TestdataCol.FlashId, "flash_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "flash_id":`)
		WriteJsonString(&buf, t.FlashId)
	}
	isRender = isRenderField(TestdataCol.TestResult, "test_result", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "test_result":`)
		WriteJsonString(&buf, t.TestResult)
	}
	isRender = isRenderField(TestdataCol.TestMsg, "test_msg", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "test_msg":`)
		WriteJsonString(&buf, t.TestMsg)
	}
	isRender = isRenderField(TestdataCol.FactorySid, "factory_sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "factory_sid":`)
		WriteJsonString(&buf, t.FactorySid)
	}
	isRender = isRenderField(TestdataCol.BatchSid, "batch_sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "batch_sid":`)
		WriteJsonString(&buf, t.BatchSid)
	}
	isRender = isRenderField(TestdataCol.Efuse, "efuse", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "efuse":`)
		WriteJsonString(&buf, t.Efuse)
	}
	isRender = isRenderField(TestdataCol.QueryTimes, "query_times", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "query_times":`)
		buf.WriteString(strconv.FormatInt(int64(t.QueryTimes), 10))
	}
	isRender = isRenderField(TestdataCol.PrintTimes, "print_times", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "print_times":`)
		buf.WriteString(strconv.FormatInt(int64(t.PrintTimes), 10))
	}
	isRender = isRenderField(TestdataCol.BatchIndex, "batch_index", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "batch_index":`)
		buf.WriteString(strconv.FormatInt(int64(t.BatchIndex), 10))
	}
	isRender = isRenderField(TestdataCol.Latest, "latest", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "latest":`)
		if numericEnum {
			if t.Latest {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(t.Latest))
		}
	}
	isRender = isRenderField(TestdataCol.IsCommit, "is_commit", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "is_commit":`)
		if numericEnum {
			if t.IsCommit {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(t.IsCommit))
		}
	}
	if paddings != nil {
		var kk string
		var vv interface{}
		var str string
		var err error
		for kk, vv = range paddings {
			buf.WriteString(`, "` + kk + `":`)
			str, err = General(vv)
			if err != nil {
				return nil, err
			}
			buf.WriteString(str)
		}
	}
	buf.WriteString("}")
	bs := buf.Bytes()
	bs[0] = '{'
	return bs, nil
}

func (t *Testdata) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalTestdata with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalTestdata")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = testdatacols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case TestdataCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			t.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			t.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case TestdataCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			t.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case TestdataCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			t.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.ModuleId:
			vvv, ok := vv["module_id"]
			if !ok {
				continue
			}
			t.ModuleId, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.DeviceType:
			vvv, ok := vv["device_type"]
			if !ok {
				continue
			}
			t.DeviceType, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.FwVer:
			vvv, ok := vv["fw_ver"]
			if !ok {
				continue
			}
			t.FwVer, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.EspMac:
			vvv, ok := vv["esp_mac"]
			if !ok {
				continue
			}
			t.EspMac, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.CusMac:
			vvv, ok := vv["cus_mac"]
			if !ok {
				continue
			}
			t.CusMac, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.FlashId:
			vvv, ok := vv["flash_id"]
			if !ok {
				continue
			}
			t.FlashId, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.TestResult:
			vvv, ok := vv["test_result"]
			if !ok {
				continue
			}
			t.TestResult, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.TestMsg:
			vvv, ok := vv["test_msg"]
			if !ok {
				continue
			}
			t.TestMsg, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.FactorySid:
			vvv, ok := vv["factory_sid"]
			if !ok {
				continue
			}
			t.FactorySid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.BatchSid:
			vvv, ok := vv["batch_sid"]
			if !ok {
				continue
			}
			t.BatchSid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.Efuse:
			vvv, ok := vv["efuse"]
			if !ok {
				continue
			}
			t.Efuse, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.QueryTimes:
			vvv, ok := vv["query_times"]
			if !ok {
				continue
			}
			t.QueryTimes, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.PrintTimes:
			vvv, ok := vv["print_times"]
			if !ok {
				continue
			}
			t.PrintTimes, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.BatchIndex:
			vvv, ok := vv["batch_index"]
			if !ok {
				continue
			}
			t.BatchIndex, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.Latest:
			vvv, ok := vv["latest"]
			if !ok {
				continue
			}
			t.Latest, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case TestdataCol.IsCommit:
			vvv, ok := vv["is_commit"]
			if !ok {
				continue
			}
			t.IsCommit, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalTestdata(ctx context.Context, vi interface{}, cols ...db.Col) (*Testdata, error) {
	t := NewTestdata(ctx)
	_, err := t.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return t, err
}

func UnmarshalTestdatas(ctx context.Context, vi interface{}, cols ...db.Col) ([]*Testdata, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalTestdatas with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalTestdatas")
	}
	testdatas := make([]*Testdata, len(vv))
	for ii, vvv := range vv {
		var t *Testdata
		t, err = UnmarshalTestdata(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		testdatas[ii] = t
	}
	return testdatas, nil
}

func newTestdataDest(cols ...string) (db.Doer, []interface{}, error) {
	t := &Testdata{}
	if cols == nil || len(cols) == 0 {
		return t, []interface{}{&t.Id, &t.Created, &t.Updated, &t.Visibly, &t.ModuleId, &t.DeviceType, &t.FwVer, &t.EspMac, &t.CusMac, &t.FlashId, &t.TestResult, &t.TestMsg, &t.FactorySid, &t.BatchSid, &t.Efuse, &t.QueryTimes, &t.PrintTimes, &t.BatchIndex, &t.Latest, &t.IsCommit}, nil
	}
	dest := make([]interface{}, len(cols))
	for ii, col := range cols {
		switch col {
		case "id":
			dest[ii] = &t.Id
		case "created":
			dest[ii] = &t.Created
		case "updated":
			dest[ii] = &t.Updated
		case "visibly":
			dest[ii] = &t.Visibly
		case "module_id":
			dest[ii] = &t.ModuleId
		case "device_type":
			dest[ii] = &t.DeviceType
		case "fw_ver":
			dest[ii] = &t.FwVer
		case "esp_mac":
			dest[ii] = &t.EspMac
		case "cus_mac":
			dest[ii] = &t.CusMac
		case "flash_id":
			dest[ii] = &t.FlashId
		case "test_result":
			dest[ii] = &t.TestResult
		case "test_msg":
			dest[ii] = &t.TestMsg
		case "factory_sid":
			dest[ii] = &t.FactorySid
		case "batch_sid":
			dest[ii] = &t.BatchSid
		case "efuse":
			dest[ii] = &t.Efuse
		case "query_times":
			dest[ii] = &t.QueryTimes
		case "print_times":
			dest[ii] = &t.PrintTimes
		case "batch_index":
			dest[ii] = &t.BatchIndex
		case "latest":
			dest[ii] = &t.Latest
		case "is_commit":
			dest[ii] = &t.IsCommit
		default:
			return nil, nil, logError("dal.Testdata Error: unknow column " + col + " in talbe testdata")
		}
	}
	return t, dest, nil
}

func colsAndArgsTestdata(t *Testdata, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.Testdata Error: at least one column to colsAndArgsTestdata")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case TestdataCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = t.Id
		case TestdataCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = t.Created
		case TestdataCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = t.Updated
		case TestdataCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = t.Visibly
		case TestdataCol.ModuleId:
			cols[ii] = "`module_id` = ?"
			args[ii] = t.ModuleId
		case TestdataCol.DeviceType:
			cols[ii] = "`device_type` = ?"
			args[ii] = t.DeviceType
		case TestdataCol.FwVer:
			cols[ii] = "`fw_ver` = ?"
			args[ii] = t.FwVer
		case TestdataCol.EspMac:
			cols[ii] = "`esp_mac` = ?"
			args[ii] = t.EspMac
		case TestdataCol.CusMac:
			cols[ii] = "`cus_mac` = ?"
			args[ii] = t.CusMac
		case TestdataCol.FlashId:
			cols[ii] = "`flash_id` = ?"
			args[ii] = t.FlashId
		case TestdataCol.TestResult:
			cols[ii] = "`test_result` = ?"
			args[ii] = t.TestResult
		case TestdataCol.TestMsg:
			cols[ii] = "`test_msg` = ?"
			args[ii] = t.TestMsg
		case TestdataCol.FactorySid:
			cols[ii] = "`factory_sid` = ?"
			args[ii] = t.FactorySid
		case TestdataCol.BatchSid:
			cols[ii] = "`batch_sid` = ?"
			args[ii] = t.BatchSid
		case TestdataCol.Efuse:
			cols[ii] = "`efuse` = ?"
			args[ii] = t.Efuse
		case TestdataCol.QueryTimes:
			cols[ii] = "`query_times` = ?"
			args[ii] = t.QueryTimes
		case TestdataCol.PrintTimes:
			cols[ii] = "`print_times` = ?"
			args[ii] = t.PrintTimes
		case TestdataCol.BatchIndex:
			cols[ii] = "`batch_index` = ?"
			args[ii] = t.BatchIndex
		case TestdataCol.Latest:
			cols[ii] = "`latest` = ?"
			args[ii] = t.Latest
		case TestdataCol.IsCommit:
			cols[ii] = "`is_commit` = ?"
			args[ii] = t.IsCommit
		default:
			return nil, nil, logError(fmt.Sprintf("dal.Testdata Error: unknow column num %d in talbe testdata", cc))
		}
	}
	return cols, args, nil
}

var TestdataEnum = struct {
}{}

var testdataSqls = []string{
	/*
		CREATE TABLE `testdata` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `module_id` int(11) NOT NULL,
		  `device_type` varchar(64) NOT NULL,
		  `fw_ver` varchar(64) NOT NULL,
		  `esp_mac` varchar(64) NOT NULL,
		  `cus_mac` varchar(64) NOT NULL,
		  `flash_id` varchar(64) NOT NULL,
		  `test_result` varchar(64) NOT NULL,
		  `test_msg` varchar(64) NOT NULL,
		  `factory_sid` varchar(64) NOT NULL,
		  `batch_sid` varchar(64) NOT NULL,
		  `efuse` varchar(64) NOT NULL,
		  `query_times` int(11) NOT NULL,
		  `print_times` int(11) NOT NULL,
		  `batch_index` int(11) NOT NULL,
		  `latest` bool NOT NULL,
		  `is_commit` bool NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into testdata(`created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit`) values(now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*1*/ "insert into testdata(`id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit`) values(?, now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*2*/ "update testdata set updated = now(), `visibly` = ?, `module_id` = ?, `device_type` = ?, `fw_ver` = ?, `esp_mac` = ?, `cus_mac` = ?, `flash_id` = ?, `test_result` = ?, `test_msg` = ?, `factory_sid` = ?, `batch_sid` = ?, `efuse` = ?, `query_times` = ?, `print_times` = ?, `batch_index` = ?, `latest` = ?, `is_commit` = ? where id = ?",
	/*3*/ "update testdata set updated = now(), %s where id = ?",
	/*4*/ "update testdata set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from testdata where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where id in (%s) and visibly = 1",

	/*8*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and esp_mac = ? order by id desc limit 0, 1",
	/*9*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and esp_mac = ? order by id desc limit ?, ?",
	/*10*/ "select count(*) from testdata where visibly = 1 and esp_mac = ? and test_result = ?",
	/*11*/ "select count(*) from testdata where visibly = 1 and batch_sid = ? and test_result = ?",
	/*12*/ "select count(*) from testdata where visibly = 1 and created >= ? and created < ? and test_result = ?",
	/*13*/ "select count(*) from testdata where visibly = 1 and created >= ? and created < ? and batch_sid = ? and test_result = ?",
	/*14*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and batch_sid = ? order by id desc limit ?, ?",
	/*15*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 order by id desc limit ?, ?",
	/*16*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and factory_sid = ? order by id desc limit ?, ?",
	/*17*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result = 'success' and latest = 1",
	/*18*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result = 'success' and not exists (select * from testdata as testdata2 where visibly = 1 and batch_sid = ? and test_result != 'success' and testdata.esp_mac = testdata2.esp_mac)",
	/*19*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result != 'success'",
	/*20*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result != 'success' and latest = 1",
	/*21*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and created >= ? and created < ?",
	/*22*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result = 'success' and latest = 1 and created >= ? and created < ?",
	/*23*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result = 'success' and created >= ? and created < ? and not exists (select * from testdata as testdata2 where visibly = 1 and batch_sid = ? and test_result != 'success' and created >= ? and created < ? and testdata.esp_mac = testdata2.esp_mac)",
	/*24*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result != 'success' and created >= ? and created < ?",
	/*25*/ "select count(distinct(esp_mac)) from testdata where visibly = 1 and batch_sid = ? and test_result != 'success' and latest = 1 and created >= ? and created < ?",
	/*26*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and cus_mac = ? order by id desc limit 0, 1",
	/*27*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and batch_sid = ? and esp_mac = ? order by id desc limit 0, 1",
	/*28*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and batch_sid = ? and cus_mac = ? order by id desc limit 0, 1",
	/*29*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and esp_mac = ? and test_result = 'success' order by id desc limit 0, 1",
	/*30*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and cus_mac = ? and test_result = 'success' order by id desc limit 0, 1",
	/*31*/ "select count(distinct(esp_mac)) from testdata where visibly=1 and batch_sid = ? and test_result='success' and latest=1 and print_times=1",
	/*32*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and factory_sid = ? and esp_mac = ? order by id desc limit ?, ?",
	/*33*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`, `is_commit` from testdata where visibly = 1 and batch_sid = ? and id=(select max(id) from testdata where visibly=1 and batch_sid= ? )",
}

package dal

import (
	"bytes"
	"context"
	"espressif.com/chip/factory/db"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var TestlogTid = 1
var _ db.Doer = (*Testlog)(nil)
var testlogcols = []db.Col{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
var testlogfields = []string{"id", "created", "updated", "visibly", "module_id", "device_type", "fw_ver", "esp_mac", "cus_mac", "flash_id", "test_result", "test_msg", "factory_sid", "batch_sid", "efuse", "query_times", "print_times", "batch_index", "latest"}

var TestlogCol = struct {
	Id, Created, Updated, Visibly, ModuleId, DeviceType, FwVer, EspMac, CusMac, FlashId, TestResult, TestMsg, FactorySid, BatchSid, Efuse, QueryTimes, PrintTimes, BatchIndex, Latest, _ db.Col
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 0}

type Testlog struct {
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

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewTestlog(ctx context.Context) *Testlog {
	now := time.Now()
	t := &Testlog{Created: now, Updated: now, Visibly: true}
	t.ext = GetExtFromContext(ctx)
	defaultTestlog(ctx, t)
	return t
}

func FindTestlog(ctx context.Context, id int) *Testlog {
	dos, err := db.Open("Testlog").Query(newTestlogDest, true, testlogSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		testlog, _ := do.(*Testlog)
		if ext != nil {
			testlog.ext = ext
		}
		return testlog
	}
	return nil
}

func ListTestlog(ctx context.Context, ids ...int) []*Testlog {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(testlogSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("Testlog").Query(newTestlogDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	testlogs := make([]*Testlog, len(dos))
	for ii, do := range dos {
		testlog, _ := do.(*Testlog)
		if ext != nil {
			testlog.ext = ext
		}
		testlogs[ii] = testlog
	}
	return testlogs
}

func (t *Testlog) Save() error {
	now := time.Now()
	t.Created, t.Updated, t.Visibly = now, now, true
	var id int64
	var err error
	if t.Id == 0 {
		id, _, err = db.Open("Testlog").Exec(testlogSqls[0], t.ModuleId, t.DeviceType, t.FwVer, t.EspMac, t.CusMac, t.FlashId, t.TestResult, t.TestMsg, t.FactorySid, t.BatchSid, t.Efuse, t.QueryTimes, t.PrintTimes, t.BatchIndex, t.Latest)
	} else {
		id, _, err = db.Open("Testlog").Exec(testlogSqls[1], t.Id, t.ModuleId, t.DeviceType, t.FwVer, t.EspMac, t.CusMac, t.FlashId, t.TestResult, t.TestMsg, t.FactorySid, t.BatchSid, t.Efuse, t.QueryTimes, t.PrintTimes, t.BatchIndex, t.Latest)
	}
	if err != nil {
		return err
	}
	t.Id = int(id)
	return nil
}

func (t *Testlog) Update(cs ...db.Col) error {
	if t.Id == 0 {
		return logError("dal.Testlog Error: can not update row while id is zero")
	}
	t.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("Testlog").Exec(testlogSqls[2], t.Visibly, t.ModuleId, t.DeviceType, t.FwVer, t.EspMac, t.CusMac, t.FlashId, t.TestResult, t.TestMsg, t.FactorySid, t.BatchSid, t.Efuse, t.QueryTimes, t.PrintTimes, t.BatchIndex, t.Latest, t.Id)
		return err
	}
	cols, args, err := colsAndArgsTestlog(t, cs...)
	if err != nil {
		return err
	}
	args = append(args, t.Id)
	sqlstr := fmt.Sprintf(testlogSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("Testlog").Exec(sqlstr, args...)
	return err
}

func (t *Testlog) Invisibly() error {
	if t.Id == 0 {
		return logError("dal.Testlog Error: can not invisibly row while id is zero")
	}
	t.Updated = time.Now()
	t.Visibly = false
	_, _, err := db.Open("Testlog").Exec(testlogSqls[4], t.Id)
	return err
}

func (t *Testlog) Delete() error {
	if t.Id == 0 {
		return logError("dal.Testlog Error: can not delete row while id is zero")
	}
	t.Updated = time.Now()
	_, _, err := db.Open("Testlog").Exec(testlogSqls[5], t.Id)
	return err
}

func (t *Testlog) Valid() error {
	return t.valid()
}

func (t *Testlog) SetExt(ext *Ext) {
	t.ext = ext
}

func (t *Testlog) Padding(pkey string, pvalue interface{}) {
	if t.ext == nil {
		t.ext = &Ext{Loc: DefaultLoc}
	}
	if t.paddings == nil {
		t.paddings = make(map[string]interface{})
	}
	t.paddings[pkey] = pvalue
	t.ext.IsComplex = true
}

func (t *Testlog) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case TestlogCol.Id:
			if isColumnName {
				mm["id"] = t.Id
			} else {
				mm["Id"] = t.Id
			}
		case TestlogCol.Created:
			if isColumnName {
				mm["created"] = t.Created
			} else {
				mm["Created"] = t.Created
			}
		case TestlogCol.Updated:
			if isColumnName {
				mm["updated"] = t.Updated
			} else {
				mm["Updated"] = t.Updated
			}
		case TestlogCol.Visibly:
			if isColumnName {
				mm["visibly"] = t.Visibly
			} else {
				mm["Visibly"] = t.Visibly
			}
		case TestlogCol.ModuleId:
			if isColumnName {
				mm["module_id"] = t.ModuleId
			} else {
				mm["ModuleId"] = t.ModuleId
			}
		case TestlogCol.DeviceType:
			if isColumnName {
				mm["device_type"] = t.DeviceType
			} else {
				mm["DeviceType"] = t.DeviceType
			}
		case TestlogCol.FwVer:
			if isColumnName {
				mm["fw_ver"] = t.FwVer
			} else {
				mm["FwVer"] = t.FwVer
			}
		case TestlogCol.EspMac:
			if isColumnName {
				mm["esp_mac"] = t.EspMac
			} else {
				mm["EspMac"] = t.EspMac
			}
		case TestlogCol.CusMac:
			if isColumnName {
				mm["cus_mac"] = t.CusMac
			} else {
				mm["CusMac"] = t.CusMac
			}
		case TestlogCol.FlashId:
			if isColumnName {
				mm["flash_id"] = t.FlashId
			} else {
				mm["FlashId"] = t.FlashId
			}
		case TestlogCol.TestResult:
			if isColumnName {
				mm["test_result"] = t.TestResult
			} else {
				mm["TestResult"] = t.TestResult
			}
		case TestlogCol.TestMsg:
			if isColumnName {
				mm["test_msg"] = t.TestMsg
			} else {
				mm["TestMsg"] = t.TestMsg
			}
		case TestlogCol.FactorySid:
			if isColumnName {
				mm["factory_sid"] = t.FactorySid
			} else {
				mm["FactorySid"] = t.FactorySid
			}
		case TestlogCol.BatchSid:
			if isColumnName {
				mm["batch_sid"] = t.BatchSid
			} else {
				mm["BatchSid"] = t.BatchSid
			}
		case TestlogCol.Efuse:
			if isColumnName {
				mm["efuse"] = t.Efuse
			} else {
				mm["Efuse"] = t.Efuse
			}
		case TestlogCol.QueryTimes:
			if isColumnName {
				mm["query_times"] = t.QueryTimes
			} else {
				mm["QueryTimes"] = t.QueryTimes
			}
		case TestlogCol.PrintTimes:
			if isColumnName {
				mm["print_times"] = t.PrintTimes
			} else {
				mm["PrintTimes"] = t.PrintTimes
			}
		case TestlogCol.BatchIndex:
			if isColumnName {
				mm["batch_index"] = t.BatchIndex
			} else {
				mm["BatchIndex"] = t.BatchIndex
			}
		case TestlogCol.Latest:
			if isColumnName {
				mm["latest"] = t.Latest
			} else {
				mm["Latest"] = t.Latest
			}
		default:
			logError(fmt.Sprintf("dal.Testlog Error: unknow column num %d in talbe testlog", cc))
		}
	}
	return mm
}

func (t *Testlog) MarshalJSON() ([]byte, error) {
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
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (t *Testlog) marshalJSONComplex() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	if t.ext == nil {
		return nil, logError("dal.Testlog Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := t.ext.Loc
	numericEnum := t.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[TestlogTid]; ok {
		if vvv, ok := vv[t.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := t.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(TestlogCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(t.Id), 10))
	}
	isRender = isRenderField(TestlogCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(t.Created.Unix(), 10))
		} else {
			t.Created = t.Created.In(loc)
			buf.WriteString(`"` + t.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(TestlogCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(t.Updated.Unix(), 10))
		} else {
			t.Updated = t.Updated.In(loc)
			buf.WriteString(`"` + t.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(TestlogCol.Visibly, "visibly", includes, excludes, paddings)
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
	isRender = isRenderField(TestlogCol.ModuleId, "module_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "module_id":`)
		buf.WriteString(strconv.FormatInt(int64(t.ModuleId), 10))
	}
	isRender = isRenderField(TestlogCol.DeviceType, "device_type", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "device_type":`)
		WriteJsonString(&buf, t.DeviceType)
	}
	isRender = isRenderField(TestlogCol.FwVer, "fw_ver", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "fw_ver":`)
		WriteJsonString(&buf, t.FwVer)
	}
	isRender = isRenderField(TestlogCol.EspMac, "esp_mac", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac":`)
		WriteJsonString(&buf, t.EspMac)
	}
	isRender = isRenderField(TestlogCol.CusMac, "cus_mac", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cus_mac":`)
		WriteJsonString(&buf, t.CusMac)
	}
	isRender = isRenderField(TestlogCol.FlashId, "flash_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "flash_id":`)
		WriteJsonString(&buf, t.FlashId)
	}
	isRender = isRenderField(TestlogCol.TestResult, "test_result", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "test_result":`)
		WriteJsonString(&buf, t.TestResult)
	}
	isRender = isRenderField(TestlogCol.TestMsg, "test_msg", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "test_msg":`)
		WriteJsonString(&buf, t.TestMsg)
	}
	isRender = isRenderField(TestlogCol.FactorySid, "factory_sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "factory_sid":`)
		WriteJsonString(&buf, t.FactorySid)
	}
	isRender = isRenderField(TestlogCol.BatchSid, "batch_sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "batch_sid":`)
		WriteJsonString(&buf, t.BatchSid)
	}
	isRender = isRenderField(TestlogCol.Efuse, "efuse", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "efuse":`)
		WriteJsonString(&buf, t.Efuse)
	}
	isRender = isRenderField(TestlogCol.QueryTimes, "query_times", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "query_times":`)
		buf.WriteString(strconv.FormatInt(int64(t.QueryTimes), 10))
	}
	isRender = isRenderField(TestlogCol.PrintTimes, "print_times", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "print_times":`)
		buf.WriteString(strconv.FormatInt(int64(t.PrintTimes), 10))
	}
	isRender = isRenderField(TestlogCol.BatchIndex, "batch_index", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "batch_index":`)
		buf.WriteString(strconv.FormatInt(int64(t.BatchIndex), 10))
	}
	isRender = isRenderField(TestlogCol.Latest, "latest", includes, excludes, paddings)
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

func (t *Testlog) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalTestlog with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalTestlog")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = testlogcols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case TestlogCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			t.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			t.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case TestlogCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			t.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case TestlogCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			t.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.ModuleId:
			vvv, ok := vv["module_id"]
			if !ok {
				continue
			}
			t.ModuleId, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.DeviceType:
			vvv, ok := vv["device_type"]
			if !ok {
				continue
			}
			t.DeviceType, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.FwVer:
			vvv, ok := vv["fw_ver"]
			if !ok {
				continue
			}
			t.FwVer, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.EspMac:
			vvv, ok := vv["esp_mac"]
			if !ok {
				continue
			}
			t.EspMac, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.CusMac:
			vvv, ok := vv["cus_mac"]
			if !ok {
				continue
			}
			t.CusMac, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.FlashId:
			vvv, ok := vv["flash_id"]
			if !ok {
				continue
			}
			t.FlashId, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.TestResult:
			vvv, ok := vv["test_result"]
			if !ok {
				continue
			}
			t.TestResult, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.TestMsg:
			vvv, ok := vv["test_msg"]
			if !ok {
				continue
			}
			t.TestMsg, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.FactorySid:
			vvv, ok := vv["factory_sid"]
			if !ok {
				continue
			}
			t.FactorySid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.BatchSid:
			vvv, ok := vv["batch_sid"]
			if !ok {
				continue
			}
			t.BatchSid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.Efuse:
			vvv, ok := vv["efuse"]
			if !ok {
				continue
			}
			t.Efuse, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.QueryTimes:
			vvv, ok := vv["query_times"]
			if !ok {
				continue
			}
			t.QueryTimes, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.PrintTimes:
			vvv, ok := vv["print_times"]
			if !ok {
				continue
			}
			t.PrintTimes, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.BatchIndex:
			vvv, ok := vv["batch_index"]
			if !ok {
				continue
			}
			t.BatchIndex, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case TestlogCol.Latest:
			vvv, ok := vv["latest"]
			if !ok {
				continue
			}
			t.Latest, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalTestlog(ctx context.Context, vi interface{}, cols ...db.Col) (*Testlog, error) {
	t := NewTestlog(ctx)
	_, err := t.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return t, err
}

func UnmarshalTestlogs(ctx context.Context, vi interface{}, cols ...db.Col) ([]*Testlog, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalTestlogs with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalTestlogs")
	}
	testlogs := make([]*Testlog, len(vv))
	for ii, vvv := range vv {
		var t *Testlog
		t, err = UnmarshalTestlog(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		testlogs[ii] = t
	}
	return testlogs, nil
}

func newTestlogDest(cols ...string) (db.Doer, []interface{}, error) {
	t := &Testlog{}
	if cols == nil || len(cols) == 0 {
		return t, []interface{}{&t.Id, &t.Created, &t.Updated, &t.Visibly, &t.ModuleId, &t.DeviceType, &t.FwVer, &t.EspMac, &t.CusMac, &t.FlashId, &t.TestResult, &t.TestMsg, &t.FactorySid, &t.BatchSid, &t.Efuse, &t.QueryTimes, &t.PrintTimes, &t.BatchIndex, &t.Latest}, nil
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
		default:
			return nil, nil, logError("dal.Testlog Error: unknow column " + col + " in talbe testlog")
		}
	}
	return t, dest, nil
}

func colsAndArgsTestlog(t *Testlog, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.Testlog Error: at least one column to colsAndArgsTestlog")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case TestlogCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = t.Id
		case TestlogCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = t.Created
		case TestlogCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = t.Updated
		case TestlogCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = t.Visibly
		case TestlogCol.ModuleId:
			cols[ii] = "`module_id` = ?"
			args[ii] = t.ModuleId
		case TestlogCol.DeviceType:
			cols[ii] = "`device_type` = ?"
			args[ii] = t.DeviceType
		case TestlogCol.FwVer:
			cols[ii] = "`fw_ver` = ?"
			args[ii] = t.FwVer
		case TestlogCol.EspMac:
			cols[ii] = "`esp_mac` = ?"
			args[ii] = t.EspMac
		case TestlogCol.CusMac:
			cols[ii] = "`cus_mac` = ?"
			args[ii] = t.CusMac
		case TestlogCol.FlashId:
			cols[ii] = "`flash_id` = ?"
			args[ii] = t.FlashId
		case TestlogCol.TestResult:
			cols[ii] = "`test_result` = ?"
			args[ii] = t.TestResult
		case TestlogCol.TestMsg:
			cols[ii] = "`test_msg` = ?"
			args[ii] = t.TestMsg
		case TestlogCol.FactorySid:
			cols[ii] = "`factory_sid` = ?"
			args[ii] = t.FactorySid
		case TestlogCol.BatchSid:
			cols[ii] = "`batch_sid` = ?"
			args[ii] = t.BatchSid
		case TestlogCol.Efuse:
			cols[ii] = "`efuse` = ?"
			args[ii] = t.Efuse
		case TestlogCol.QueryTimes:
			cols[ii] = "`query_times` = ?"
			args[ii] = t.QueryTimes
		case TestlogCol.PrintTimes:
			cols[ii] = "`print_times` = ?"
			args[ii] = t.PrintTimes
		case TestlogCol.BatchIndex:
			cols[ii] = "`batch_index` = ?"
			args[ii] = t.BatchIndex
		case TestlogCol.Latest:
			cols[ii] = "`latest` = ?"
			args[ii] = t.Latest
		default:
			return nil, nil, logError(fmt.Sprintf("dal.Testlog Error: unknow column num %d in talbe testlog", cc))
		}
	}
	return cols, args, nil
}

var TestlogEnum = struct {
}{}

var testlogSqls = []string{
	/*
		CREATE TABLE `testlog` (
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
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into testlog(`created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`) values(now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*1*/ "insert into testlog(`id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest`) values(?, now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*2*/ "update testlog set updated = now(), `visibly` = ?, `module_id` = ?, `device_type` = ?, `fw_ver` = ?, `esp_mac` = ?, `cus_mac` = ?, `flash_id` = ?, `test_result` = ?, `test_msg` = ?, `factory_sid` = ?, `batch_sid` = ?, `efuse` = ?, `query_times` = ?, `print_times` = ?, `batch_index` = ?, `latest` = ? where id = ?",
	/*3*/ "update testlog set updated = now(), %s where id = ?",
	/*4*/ "update testlog set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from testlog where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest` from testlog where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `module_id`, `device_type`, `fw_ver`, `esp_mac`, `cus_mac`, `flash_id`, `test_result`, `test_msg`, `factory_sid`, `batch_sid`, `efuse`, `query_times`, `print_times`, `batch_index`, `latest` from testlog where id in (%s) and visibly = 1",
}

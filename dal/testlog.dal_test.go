package dal

import (
	"encoding/json"
	"espressif.com/chip/factory/db"
	"log"
	"strings"
	"testing"
	"time"
)

func TestTestlogSave(tt *testing.T) {
	t := NewTestlog(nil)
	randSetTestlog(tt, t)
	now := time.Now()
	t.Id = 0
	t.Save()
	if t.Id <= 0 {
		log.Println("save error?")
		tt.Fatal("after Save(), t.Id must be great than zero")
	}
	if !t.Visibly {
		tt.Fatal("after Save(), t.Visibly must be true")
	}
	if t.Created != t.Updated {
		tt.Fatal("after Save(), t.Created must equals t.Updated")
	}
	if !almostSameTime(t.Created, now, 1) {
		tt.Fatal("after Save(), t.Created must be now")
	}
	if !almostSameTime(t.Updated, now, 1) {
		tt.Fatal("after Save(), t.Updated must be now")
	}
	compareTestlog(tt, t)
}

func randSetTestlog(tt *testing.T, t *Testlog) {
	t.Id = randInt32()
	t.Created = randTime()
	t.Updated = randTime()
	t.Visibly = randBool()
	t.ModuleId = randInt32()
	t.DeviceType = randStr(64)
	t.FwVer = randStr(64)
	t.EspMac = randStr(64)
	t.CusMac = randStr(64)
	t.FlashId = randStr(64)
	t.TestResult = randStr(64)
	t.TestMsg = randStr(64)
	t.FactorySid = randStr(64)
	t.BatchSid = randStr(64)
	t.Efuse = randStr(64)
	t.QueryTimes = randInt32()
	t.PrintTimes = randInt32()
	t.BatchIndex = randInt32()
	t.Latest = randBool()
}

func compareTestlog(tt *testing.T, t *Testlog) {
	t2 := FindTestlog(nil, t.Id)
	if t.Id != t2.Id {
		tt.Fatal("insert and find compare failed, field: Id")
	}
	if !almostSameTime(t.Created, t2.Created, 1) {
		tt.Fatal("insert and find compare failed, field: Created")
	}
	if !almostSameTime(t.Updated, t2.Updated, 1) {
		tt.Fatal("insert and find compare failed, field: Updated")
	}
	if t.Visibly != t2.Visibly {
		tt.Fatal("insert and find compare failed, field: Visibly")
	}
	if t.ModuleId != t2.ModuleId {
		tt.Fatal("insert and find compare failed, field: ModuleId")
	}
	if t.DeviceType != t2.DeviceType {
		tt.Fatal("insert and find compare failed, field: DeviceType")
	}
	if t.FwVer != t2.FwVer {
		tt.Fatal("insert and find compare failed, field: FwVer")
	}
	if t.EspMac != t2.EspMac {
		tt.Fatal("insert and find compare failed, field: EspMac")
	}
	if t.CusMac != t2.CusMac {
		tt.Fatal("insert and find compare failed, field: CusMac")
	}
	if t.FlashId != t2.FlashId {
		tt.Fatal("insert and find compare failed, field: FlashId")
	}
	if t.TestResult != t2.TestResult {
		tt.Fatal("insert and find compare failed, field: TestResult")
	}
	if t.TestMsg != t2.TestMsg {
		tt.Fatal("insert and find compare failed, field: TestMsg")
	}
	if t.FactorySid != t2.FactorySid {
		tt.Fatal("insert and find compare failed, field: FactorySid")
	}
	if t.BatchSid != t2.BatchSid {
		tt.Fatal("insert and find compare failed, field: BatchSid")
	}
	if t.Efuse != t2.Efuse {
		tt.Fatal("insert and find compare failed, field: Efuse")
	}
	if t.QueryTimes != t2.QueryTimes {
		tt.Fatal("insert and find compare failed, field: QueryTimes")
	}
	if t.PrintTimes != t2.PrintTimes {
		tt.Fatal("insert and find compare failed, field: PrintTimes")
	}
	if t.BatchIndex != t2.BatchIndex {
		tt.Fatal("insert and find compare failed, field: BatchIndex")
	}
	if t.Latest != t2.Latest {
		tt.Fatal("insert and find compare failed, field: Latest")
	}
}

func TestTestlogUpdate(tt *testing.T) {
	t := NewTestlog(nil)
	randSetTestlog(tt, t)
	t.Id = 0
	t.Save()
	id := t.Id
	created := t.Created
	randSetTestlog(tt, t)
	t.Id = id
	t.Visibly = true
	t.Created = created
	t.Update()
	if t.Created.After(t.Updated) || t.Created.Equal(t.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareTestlog(tt, t)
}

func TestTestlogInvisibly(tt *testing.T) {
	t := NewTestlog(nil)
	randSetTestlog(tt, t)
	t.Id = 0
	t.Save()
	t.Invisibly()
	t2 := FindTestlog(nil, t.Id)
	if t2 != nil {
		tt.Fatal("after Invisibly, FindTestlog() must return nil")
	}
}

func TestTestlogDelete(tt *testing.T) {
	t := NewTestlog(nil)
	randSetTestlog(tt, t)
	t.Id = 0
	t.Save()
	t.Delete()
	t2 := FindTestlog(nil, t.Id)
	if t2 != nil {
		tt.Fatal("after Invisibly, FindTestlog() must return nil")
	}
}

func TestTestlogUnmarshalMap(tt *testing.T) {
	t := NewTestlog(nil)
	mm := make(map[string]interface{})
	id := randInt()
	created := randTime()
	updated := randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err := t.UnmarshalMap(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	if !compareTestlogValue(mm["id"], t.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareTestlogValue(mm["created"], t.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareTestlogValue(mm["updated"], t.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = t.UnmarshalMap(nil, mm, TestlogCol.Id)
	if !compareTestlogValue(mm["id"], t.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareTestlogValue(mm["created"], t.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareTestlogValue(mm["updated"], t.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareTestlogJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareTestlogValue(jsonValue, fieldValue)
}

func compareTestlogValue(jsonValue interface{}, fieldValue interface{}) bool {
	switch v1 := jsonValue.(type) {
	case int:
		switch v2 := fieldValue.(type) {
		case int:
			return v1 == v2
		}
	case string:
		switch v2 := fieldValue.(type) {
		case string:
			return v1 == v2
		case time.Time:
			v3, err := Time(v1, DefaultLoc)
			if err != nil {
				return false
			}
			return almostSameTime(v3, v2, 1)
		}
	case float64:
		switch v2 := fieldValue.(type) {
		case int:
			return almostSameFloat(v1, float64(v2), 0.01)
		case int64:
			return almostSameFloat(v1, float64(v2), 0.01)
		case float64:
			return almostSameFloat(v1, v2, 0.01)
		case float32:
			return almostSameFloat(v1, float64(v2), 0.01)
		}
	case bool:
		switch v2 := fieldValue.(type) {
		case bool:
			return v1 == v2
		}
	case time.Time:
		switch v2 := fieldValue.(type) {
		case time.Time:
			return almostSameTime(v1, v2, 1)
		}
	}
	return false
}

func marshalAndUnmarshalTestlog(tt *testing.T, t *Testlog) map[string]interface{} {
	bs, err := json.Marshal(t)
	if err != nil {
		tt.Fatal(err)
	}
	mm := make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	return mm
}

func TestTestlogExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	t := NewTestlog(nil)
	t.SetExt(ext)
	randSetTestlog(tt, t)
	mm := marshalAndUnmarshalTestlog(tt, t)
	mmLen := len(mm)
	vv, ok := mm["created"].(string)
	if !ok {
		tt.Fatal("time type must be string in json")
	}
	if !strings.HasSuffix(vv, "+05:00") {
		// tt.Fatal("ext.loc has not affect")
	}
	t.ext.Verbose = "v"
	t.ext.IsComplex = true
	if _, ok := dalVerboses[TestlogTid]; !ok {
		dalVerboses[TestlogTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[TestlogTid][t.ext.Verbose]
	dalVerboses[TestlogTid][t.ext.Verbose] = []map[db.Col]interface{}{
		{TestlogCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalTestlog(tt, t)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareTestlogValue(id, t.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[TestlogTid][t.ext.Verbose] = []map[db.Col]interface{}{
		{}, {TestlogCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalTestlog(tt, t)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[TestlogTid][t.ext.Verbose] = origin
}

func TestTestlogPadding(tt *testing.T) {
	t := NewTestlog(nil)
	key1 := randStr(16)
	key2 := randStr(16)
	key3 := randStr(16)
	key4 := randStr(16)
	key5 := randStr(16)
	val1 := randStr(16)
	val2 := randFloat64()
	val3 := randInt()
	val4 := randBool()
	val5 := randTime()
	t.Padding(key1, val1)
	t.Padding(key2, val2)
	t.Padding(key3, val3)
	t.Padding(key4, val4)
	t.Padding(key5, val5)
	mm := marshalAndUnmarshalTestlog(tt, t)
	if !compareTestlogValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareTestlogValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareTestlogValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareTestlogValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareTestlogValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestTestlogMarshalJSON(tt *testing.T) {
	t := NewTestlog(nil)
	randSetTestlog(tt, t)
	mm := marshalAndUnmarshalTestlog(tt, t)
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareTestlogValue(jsonValue, t.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareTestlogValue(jsonValue, t.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareTestlogValue(jsonValue, t.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareTestlogValue(jsonValue, t.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["module_id"]
	if !compareTestlogValue(jsonValue, t.ModuleId) {
		tt.Fatal("json Marshal and Unmarshal compare field (ModuleId) failed")
	}
	jsonValue = mm["device_type"]
	if !compareTestlogValue(jsonValue, t.DeviceType) {
		tt.Fatal("json Marshal and Unmarshal compare field (DeviceType) failed")
	}
	jsonValue = mm["fw_ver"]
	if !compareTestlogValue(jsonValue, t.FwVer) {
		tt.Fatal("json Marshal and Unmarshal compare field (FwVer) failed")
	}
	jsonValue = mm["esp_mac"]
	if !compareTestlogValue(jsonValue, t.EspMac) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMac) failed")
	}
	jsonValue = mm["cus_mac"]
	if !compareTestlogValue(jsonValue, t.CusMac) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMac) failed")
	}
	jsonValue = mm["flash_id"]
	if !compareTestlogValue(jsonValue, t.FlashId) {
		tt.Fatal("json Marshal and Unmarshal compare field (FlashId) failed")
	}
	jsonValue = mm["test_result"]
	if !compareTestlogValue(jsonValue, t.TestResult) {
		tt.Fatal("json Marshal and Unmarshal compare field (TestResult) failed")
	}
	jsonValue = mm["test_msg"]
	if !compareTestlogValue(jsonValue, t.TestMsg) {
		tt.Fatal("json Marshal and Unmarshal compare field (TestMsg) failed")
	}
	jsonValue = mm["factory_sid"]
	if !compareTestlogValue(jsonValue, t.FactorySid) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactorySid) failed")
	}
	jsonValue = mm["batch_sid"]
	if !compareTestlogValue(jsonValue, t.BatchSid) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchSid) failed")
	}
	jsonValue = mm["efuse"]
	if !compareTestlogValue(jsonValue, t.Efuse) {
		tt.Fatal("json Marshal and Unmarshal compare field (Efuse) failed")
	}
	jsonValue = mm["query_times"]
	if !compareTestlogValue(jsonValue, t.QueryTimes) {
		tt.Fatal("json Marshal and Unmarshal compare field (QueryTimes) failed")
	}
	jsonValue = mm["print_times"]
	if !compareTestlogValue(jsonValue, t.PrintTimes) {
		tt.Fatal("json Marshal and Unmarshal compare field (PrintTimes) failed")
	}
	jsonValue = mm["batch_index"]
	if !compareTestlogValue(jsonValue, t.BatchIndex) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchIndex) failed")
	}
	jsonValue = mm["latest"]
	if !compareTestlogValue(jsonValue, t.Latest) {
		tt.Fatal("json Marshal and Unmarshal compare field (Latest) failed")
	}
}

func TestTestlogMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	t := NewTestlog(nil)
	t.SetExt(ext)
	randSetTestlog(tt, t)
	origin := dalVerboses[TestlogTid][t.ext.Verbose]
	dalVerboses[TestlogTid][t.ext.Verbose] = []map[db.Col]interface{}{
		{}, {TestlogCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalTestlog(tt, t)
	vv, ok := mm["created"].(string)
	if !ok {
		tt.Fatal("time type must be string in json")
	}
	if !strings.HasSuffix(vv, "+08:00") {
		// tt.Fatal("ext.loc has not affect")
	}
	_, ok = mm["updated"]
	if ok {
		tt.Fatal("MarshalJSONComplex must exclude Update")
	}
	rint := randInt()
	rstr := randStr(16)
	t.Padding("id", rint)
	t.Padding("updated", rstr)
	mm = marshalAndUnmarshalTestlog(tt, t)
	if !compareTestlogValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareTestlogValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[TestlogTid][t.ext.Verbose] = origin
}

func TestTestlogUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["module_id"] = randInt32()
	mm["device_type"] = randStr(64)
	mm["fw_ver"] = randStr(64)
	mm["esp_mac"] = randStr(64)
	mm["cus_mac"] = randStr(64)
	mm["flash_id"] = randStr(64)
	mm["test_result"] = randStr(64)
	mm["test_msg"] = randStr(64)
	mm["factory_sid"] = randStr(64)
	mm["batch_sid"] = randStr(64)
	mm["efuse"] = randStr(64)
	mm["query_times"] = randInt32()
	mm["print_times"] = randInt32()
	mm["batch_index"] = randInt32()
	mm["latest"] = randBool()
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	t, err := UnmarshalTestlog(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareTestlogValue(jsonValue, t.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareTestlogValue(jsonValue, t.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareTestlogValue(jsonValue, t.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareTestlogValue(jsonValue, t.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["module_id"]
	if !compareTestlogValue(jsonValue, t.ModuleId) {
		tt.Fatal("json Marshal and Unmarshal compare field (ModuleId) failed")
	}
	jsonValue = mm["device_type"]
	if !compareTestlogValue(jsonValue, t.DeviceType) {
		tt.Fatal("json Marshal and Unmarshal compare field (DeviceType) failed")
	}
	jsonValue = mm["fw_ver"]
	if !compareTestlogValue(jsonValue, t.FwVer) {
		tt.Fatal("json Marshal and Unmarshal compare field (FwVer) failed")
	}
	jsonValue = mm["esp_mac"]
	if !compareTestlogValue(jsonValue, t.EspMac) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMac) failed")
	}
	jsonValue = mm["cus_mac"]
	if !compareTestlogValue(jsonValue, t.CusMac) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMac) failed")
	}
	jsonValue = mm["flash_id"]
	if !compareTestlogValue(jsonValue, t.FlashId) {
		tt.Fatal("json Marshal and Unmarshal compare field (FlashId) failed")
	}
	jsonValue = mm["test_result"]
	if !compareTestlogValue(jsonValue, t.TestResult) {
		tt.Fatal("json Marshal and Unmarshal compare field (TestResult) failed")
	}
	jsonValue = mm["test_msg"]
	if !compareTestlogValue(jsonValue, t.TestMsg) {
		tt.Fatal("json Marshal and Unmarshal compare field (TestMsg) failed")
	}
	jsonValue = mm["factory_sid"]
	if !compareTestlogValue(jsonValue, t.FactorySid) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactorySid) failed")
	}
	jsonValue = mm["batch_sid"]
	if !compareTestlogValue(jsonValue, t.BatchSid) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchSid) failed")
	}
	jsonValue = mm["efuse"]
	if !compareTestlogValue(jsonValue, t.Efuse) {
		tt.Fatal("json Marshal and Unmarshal compare field (Efuse) failed")
	}
	jsonValue = mm["query_times"]
	if !compareTestlogValue(jsonValue, t.QueryTimes) {
		tt.Fatal("json Marshal and Unmarshal compare field (QueryTimes) failed")
	}
	jsonValue = mm["print_times"]
	if !compareTestlogValue(jsonValue, t.PrintTimes) {
		tt.Fatal("json Marshal and Unmarshal compare field (PrintTimes) failed")
	}
	jsonValue = mm["batch_index"]
	if !compareTestlogValue(jsonValue, t.BatchIndex) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchIndex) failed")
	}
	jsonValue = mm["latest"]
	if !compareTestlogValue(jsonValue, t.Latest) {
		tt.Fatal("json Marshal and Unmarshal compare field (Latest) failed")
	}
}

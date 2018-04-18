package dal

import (
	"encoding/json"
	"testing"
	"time"
	"log"
	"strings"
	"espressif.com/chip/factory/db"
)

func TestBatchSave(tt *testing.T) {
	b := NewBatch(nil)
	randSetBatch(tt, b)
	now := time.Now()
	b.Id = 0
	b.Save()
	if b.Id <= 0 {
		log.Println("save error?")
		tt.Fatal("after Save(), b.Id must be great than zero")
	}
	if !b.Visibly {
		tt.Fatal("after Save(), b.Visibly must be true")
	}
	if b.Created != b.Updated {
		tt.Fatal("after Save(), b.Created must equals b.Updated")
	}
	if !almostSameTime(b.Created, now, 1) {
		tt.Fatal("after Save(), b.Created must be now")
	}
	if !almostSameTime(b.Updated, now, 1) {
		tt.Fatal("after Save(), b.Updated must be now")
	}
	compareBatch(tt, b)
}

func randSetBatch(tt *testing.T, b *Batch) {
	b.Id = randInt32()
	b.Created = randTime()
	b.Updated = randTime()
	b.Visibly = randBool()
	b.Sid = randStr(64)
	b.FactorySid = randStr(64)
	b.Name = randStr(64)
	b.Desc = randStr(128)
	b.Cnt = randInt32()
	b.Remain = randInt32()
	b.EspMacFrom = randStr(64)
	b.EspMacTo = randStr(64)
	b.CusMacFrom = randStr(64)
	b.CusMacTo = randStr(64)
	b.EspMacNumFrom = randInt32()
	b.EspMacNumTo = randInt32()
	b.CusMacNumFrom = randInt32()
	b.CusMacNumTo = randInt32()
	b.IsCus = randBool()
	b.Success = randInt32()
	b.RightFirstTime = randInt32()
	b.Failed = randInt32()
	b.Rejected = randInt32()
	b.Statsed = randTime()
	b.PrintNum = randInt32()
}

func compareBatch(tt *testing.T, b *Batch) {
	b2 := FindBatch(nil, b.Id)
	if b.Id != b2.Id {
		tt.Fatal("insert and find compare failed, field: Id")
	}
	if !almostSameTime(b.Created, b2.Created, 1) {
		tt.Fatal("insert and find compare failed, field: Created")
	}
	if !almostSameTime(b.Updated, b2.Updated, 1) {
		tt.Fatal("insert and find compare failed, field: Updated")
	}
	if b.Visibly != b2.Visibly {
		tt.Fatal("insert and find compare failed, field: Visibly")
	}
	if b.Sid != b2.Sid {
		tt.Fatal("insert and find compare failed, field: Sid")
	}
	if b.FactorySid != b2.FactorySid {
		tt.Fatal("insert and find compare failed, field: FactorySid")
	}
	if b.Name != b2.Name {
		tt.Fatal("insert and find compare failed, field: Name")
	}
	if b.Desc != b2.Desc {
		tt.Fatal("insert and find compare failed, field: Desc")
	}
	if b.Cnt != b2.Cnt {
		tt.Fatal("insert and find compare failed, field: Cnt")
	}
	if b.Remain != b2.Remain {
		tt.Fatal("insert and find compare failed, field: Remain")
	}
	if b.EspMacFrom != b2.EspMacFrom {
		tt.Fatal("insert and find compare failed, field: EspMacFrom")
	}
	if b.EspMacTo != b2.EspMacTo {
		tt.Fatal("insert and find compare failed, field: EspMacTo")
	}
	if b.CusMacFrom != b2.CusMacFrom {
		tt.Fatal("insert and find compare failed, field: CusMacFrom")
	}
	if b.CusMacTo != b2.CusMacTo {
		tt.Fatal("insert and find compare failed, field: CusMacTo")
	}
	if b.EspMacNumFrom != b2.EspMacNumFrom {
		tt.Fatal("insert and find compare failed, field: EspMacNumFrom")
	}
	if b.EspMacNumTo != b2.EspMacNumTo {
		tt.Fatal("insert and find compare failed, field: EspMacNumTo")
	}
	if b.CusMacNumFrom != b2.CusMacNumFrom {
		tt.Fatal("insert and find compare failed, field: CusMacNumFrom")
	}
	if b.CusMacNumTo != b2.CusMacNumTo {
		tt.Fatal("insert and find compare failed, field: CusMacNumTo")
	}
	if b.IsCus != b2.IsCus {
		tt.Fatal("insert and find compare failed, field: IsCus")
	}
	if b.Success != b2.Success {
		tt.Fatal("insert and find compare failed, field: Success")
	}
	if b.RightFirstTime != b2.RightFirstTime {
		tt.Fatal("insert and find compare failed, field: RightFirstTime")
	}
	if b.Failed != b2.Failed {
		tt.Fatal("insert and find compare failed, field: Failed")
	}
	if b.Rejected != b2.Rejected {
		tt.Fatal("insert and find compare failed, field: Rejected")
	}
	if !almostSameTime(b.Statsed, b2.Statsed, 1) {
		tt.Fatal("insert and find compare failed, field: Statsed")
	}
	if b.PrintNum != b2.PrintNum {
		tt.Fatal("insert and find compare failed, field: PrintNum")
	}
}

func TestBatchUpdate(tt *testing.T) {
	b := NewBatch(nil)
	randSetBatch(tt, b)
	b.Id = 0
	b.Save()
	id := b.Id
	created := b.Created
	randSetBatch(tt, b)
	b.Id = id
	b.Visibly = true
	b.Created = created
	b.Update()
	if b.Created.After(b.Updated) || b.Created.Equal(b.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareBatch(tt, b)
}

func TestBatchInvisibly(tt *testing.T) {
	b := NewBatch(nil)
	randSetBatch(tt, b)
	b.Id = 0
	b.Save()
	b.Invisibly()
	b2 := FindBatch(nil, b.Id)
	if b2 != nil {
		tt.Fatal("after Invisibly, FindBatch() must return nil")
	}
}

func TestBatchDelete(tt *testing.T) {
	b := NewBatch(nil)
	randSetBatch(tt, b)
	b.Id = 0
	b.Save()
	b.Delete()
	b2 := FindBatch(nil, b.Id)
	if b2 != nil {
		tt.Fatal("after Invisibly, FindBatch() must return nil")
	}
}

func TestBatchUnmarshalMap(tt *testing.T) {
	b := NewBatch(nil)
	mm := make(map[string]interface{})
	id := randInt()
	created := randTime()
	updated := randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err := b.UnmarshalMap(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	if !compareBatchValue(mm["id"], b.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareBatchValue(mm["created"], b.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareBatchValue(mm["updated"], b.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = b.UnmarshalMap(nil, mm, BatchCol.Id)
	if !compareBatchValue(mm["id"], b.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareBatchValue(mm["created"], b.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareBatchValue(mm["updated"], b.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareBatchJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareBatchValue(jsonValue, fieldValue)
}

func compareBatchValue(jsonValue interface{}, fieldValue interface{}) bool {
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

func marshalAndUnmarshalBatch(tt *testing.T, b *Batch) map[string]interface{} {
	bs, err := json.Marshal(b)
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

func TestBatchExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	b := NewBatch(nil)
	b.SetExt(ext)
	randSetBatch(tt, b)
	mm := marshalAndUnmarshalBatch(tt, b)
	mmLen := len(mm)
	vv, ok := mm["created"].(string)
	if !ok {
		tt.Fatal("time type must be string in json")
	}
	if !strings.HasSuffix(vv, "+05:00") {
		// tt.Fatal("ext.loc has not affect")
	}
	b.ext.Verbose = "v"
	b.ext.IsComplex = true
	if _, ok := dalVerboses[BatchTid]; !ok {
		dalVerboses[BatchTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[BatchTid][b.ext.Verbose]
	dalVerboses[BatchTid][b.ext.Verbose] = []map[db.Col]interface{}{
		{BatchCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalBatch(tt, b)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareBatchValue(id, b.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[BatchTid][b.ext.Verbose] = []map[db.Col]interface{}{
		{}, {BatchCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalBatch(tt, b)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[BatchTid][b.ext.Verbose] = origin
}

func TestBatchPadding(tt *testing.T) {
	b := NewBatch(nil)
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
	b.Padding(key1, val1)
	b.Padding(key2, val2)
	b.Padding(key3, val3)
	b.Padding(key4, val4)
	b.Padding(key5, val5)
	mm := marshalAndUnmarshalBatch(tt, b)
	if !compareBatchValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareBatchValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareBatchValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareBatchValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareBatchValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestBatchMarshalJSON(tt *testing.T) {
	b := NewBatch(nil)
	randSetBatch(tt, b)
	mm := marshalAndUnmarshalBatch(tt, b)
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareBatchValue(jsonValue, b.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareBatchValue(jsonValue, b.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareBatchValue(jsonValue, b.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareBatchValue(jsonValue, b.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["sid"]
	if !compareBatchValue(jsonValue, b.Sid) {
		tt.Fatal("json Marshal and Unmarshal compare field (Sid) failed")
	}
	jsonValue = mm["factory_sid"]
	if !compareBatchValue(jsonValue, b.FactorySid) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactorySid) failed")
	}
	jsonValue = mm["name"]
	if !compareBatchValue(jsonValue, b.Name) {
		tt.Fatal("json Marshal and Unmarshal compare field (Name) failed")
	}
	jsonValue = mm["desc"]
	if !compareBatchValue(jsonValue, b.Desc) {
		tt.Fatal("json Marshal and Unmarshal compare field (Desc) failed")
	}
	jsonValue = mm["cnt"]
	if !compareBatchValue(jsonValue, b.Cnt) {
		tt.Fatal("json Marshal and Unmarshal compare field (Cnt) failed")
	}
	jsonValue = mm["remain"]
	if !compareBatchValue(jsonValue, b.Remain) {
		tt.Fatal("json Marshal and Unmarshal compare field (Remain) failed")
	}
	jsonValue = mm["esp_mac_from"]
	if !compareBatchValue(jsonValue, b.EspMacFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacFrom) failed")
	}
	jsonValue = mm["esp_mac_to"]
	if !compareBatchValue(jsonValue, b.EspMacTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacTo) failed")
	}
	jsonValue = mm["cus_mac_from"]
	if !compareBatchValue(jsonValue, b.CusMacFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacFrom) failed")
	}
	jsonValue = mm["cus_mac_to"]
	if !compareBatchValue(jsonValue, b.CusMacTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacTo) failed")
	}
	jsonValue = mm["esp_mac_num_from"]
	if !compareBatchValue(jsonValue, b.EspMacNumFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacNumFrom) failed")
	}
	jsonValue = mm["esp_mac_num_to"]
	if !compareBatchValue(jsonValue, b.EspMacNumTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacNumTo) failed")
	}
	jsonValue = mm["cus_mac_num_from"]
	if !compareBatchValue(jsonValue, b.CusMacNumFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacNumFrom) failed")
	}
	jsonValue = mm["cus_mac_num_to"]
	if !compareBatchValue(jsonValue, b.CusMacNumTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacNumTo) failed")
	}
	jsonValue = mm["is_cus"]
	if !compareBatchValue(jsonValue, b.IsCus) {
		tt.Fatal("json Marshal and Unmarshal compare field (IsCus) failed")
	}
	jsonValue = mm["success"]
	if !compareBatchValue(jsonValue, b.Success) {
		tt.Fatal("json Marshal and Unmarshal compare field (Success) failed")
	}
	jsonValue = mm["right_first_time"]
	if !compareBatchValue(jsonValue, b.RightFirstTime) {
		tt.Fatal("json Marshal and Unmarshal compare field (RightFirstTime) failed")
	}
	jsonValue = mm["failed"]
	if !compareBatchValue(jsonValue, b.Failed) {
		tt.Fatal("json Marshal and Unmarshal compare field (Failed) failed")
	}
	jsonValue = mm["rejected"]
	if !compareBatchValue(jsonValue, b.Rejected) {
		tt.Fatal("json Marshal and Unmarshal compare field (Rejected) failed")
	}
	jsonValue = mm["statsed"]
	if !compareBatchValue(jsonValue, b.Statsed) {
		tt.Fatal("json Marshal and Unmarshal compare field (Statsed) failed")
	}
	jsonValue = mm["print_num"]
	if !compareBatchValue(jsonValue, b.PrintNum) {
		tt.Fatal("json Marshal and Unmarshal compare field (PrintNum) failed")
	}
}

func TestBatchMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	b := NewBatch(nil)
	b.SetExt(ext)
	randSetBatch(tt, b)
	origin := dalVerboses[BatchTid][b.ext.Verbose]
	dalVerboses[BatchTid][b.ext.Verbose] = []map[db.Col]interface{}{
		{}, {BatchCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalBatch(tt, b)
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
	b.Padding("id", rint)
	b.Padding("updated", rstr)
	mm = marshalAndUnmarshalBatch(tt, b)
	if !compareBatchValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareBatchValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[BatchTid][b.ext.Verbose] = origin
}

func TestBatchUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["sid"] = randStr(64)
	mm["factory_sid"] = randStr(64)
	mm["name"] = randStr(64)
	mm["desc"] = randStr(128)
	mm["cnt"] = randInt32()
	mm["remain"] = randInt32()
	mm["esp_mac_from"] = randStr(64)
	mm["esp_mac_to"] = randStr(64)
	mm["cus_mac_from"] = randStr(64)
	mm["cus_mac_to"] = randStr(64)
	mm["esp_mac_num_from"] = randInt32()
	mm["esp_mac_num_to"] = randInt32()
	mm["cus_mac_num_from"] = randInt32()
	mm["cus_mac_num_to"] = randInt32()
	mm["is_cus"] = randBool()
	mm["success"] = randInt32()
	mm["right_first_time"] = randInt32()
	mm["failed"] = randInt32()
	mm["rejected"] = randInt32()
	mm["statsed"] = randTime()
	mm["print_num"] = randInt32()
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	b, err := UnmarshalBatch(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareBatchValue(jsonValue, b.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareBatchValue(jsonValue, b.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareBatchValue(jsonValue, b.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareBatchValue(jsonValue, b.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["sid"]
	if !compareBatchValue(jsonValue, b.Sid) {
		tt.Fatal("json Marshal and Unmarshal compare field (Sid) failed")
	}
	jsonValue = mm["factory_sid"]
	if !compareBatchValue(jsonValue, b.FactorySid) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactorySid) failed")
	}
	jsonValue = mm["name"]
	if !compareBatchValue(jsonValue, b.Name) {
		tt.Fatal("json Marshal and Unmarshal compare field (Name) failed")
	}
	jsonValue = mm["desc"]
	if !compareBatchValue(jsonValue, b.Desc) {
		tt.Fatal("json Marshal and Unmarshal compare field (Desc) failed")
	}
	jsonValue = mm["cnt"]
	if !compareBatchValue(jsonValue, b.Cnt) {
		tt.Fatal("json Marshal and Unmarshal compare field (Cnt) failed")
	}
	jsonValue = mm["remain"]
	if !compareBatchValue(jsonValue, b.Remain) {
		tt.Fatal("json Marshal and Unmarshal compare field (Remain) failed")
	}
	jsonValue = mm["esp_mac_from"]
	if !compareBatchValue(jsonValue, b.EspMacFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacFrom) failed")
	}
	jsonValue = mm["esp_mac_to"]
	if !compareBatchValue(jsonValue, b.EspMacTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacTo) failed")
	}
	jsonValue = mm["cus_mac_from"]
	if !compareBatchValue(jsonValue, b.CusMacFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacFrom) failed")
	}
	jsonValue = mm["cus_mac_to"]
	if !compareBatchValue(jsonValue, b.CusMacTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacTo) failed")
	}
	jsonValue = mm["esp_mac_num_from"]
	if !compareBatchValue(jsonValue, b.EspMacNumFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacNumFrom) failed")
	}
	jsonValue = mm["esp_mac_num_to"]
	if !compareBatchValue(jsonValue, b.EspMacNumTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMacNumTo) failed")
	}
	jsonValue = mm["cus_mac_num_from"]
	if !compareBatchValue(jsonValue, b.CusMacNumFrom) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacNumFrom) failed")
	}
	jsonValue = mm["cus_mac_num_to"]
	if !compareBatchValue(jsonValue, b.CusMacNumTo) {
		tt.Fatal("json Marshal and Unmarshal compare field (CusMacNumTo) failed")
	}
	jsonValue = mm["is_cus"]
	if !compareBatchValue(jsonValue, b.IsCus) {
		tt.Fatal("json Marshal and Unmarshal compare field (IsCus) failed")
	}
	jsonValue = mm["success"]
	if !compareBatchValue(jsonValue, b.Success) {
		tt.Fatal("json Marshal and Unmarshal compare field (Success) failed")
	}
	jsonValue = mm["right_first_time"]
	if !compareBatchValue(jsonValue, b.RightFirstTime) {
		tt.Fatal("json Marshal and Unmarshal compare field (RightFirstTime) failed")
	}
	jsonValue = mm["failed"]
	if !compareBatchValue(jsonValue, b.Failed) {
		tt.Fatal("json Marshal and Unmarshal compare field (Failed) failed")
	}
	jsonValue = mm["rejected"]
	if !compareBatchValue(jsonValue, b.Rejected) {
		tt.Fatal("json Marshal and Unmarshal compare field (Rejected) failed")
	}
	jsonValue = mm["statsed"]
	if !compareBatchValue(jsonValue, b.Statsed) {
		tt.Fatal("json Marshal and Unmarshal compare field (Statsed) failed")
	}
	jsonValue = mm["print_num"]
	if !compareBatchValue(jsonValue, b.PrintNum) {
		tt.Fatal("json Marshal and Unmarshal compare field (PrintNum) failed")
	}
}

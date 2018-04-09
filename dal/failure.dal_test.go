package dal

import (
	"encoding/json"
	"espressif.com/chip/factory/db"
	"log"
	"strings"
	"testing"
	"time"
)

func TestFailureSave(tt *testing.T) {
	f := NewFailure(nil)
	randSetFailure(tt, f)
	now := time.Now()
	f.Id = 0
	f.Save()
	if f.Id <= 0 {
		log.Println("save error?")
		tt.Fatal("after Save(), f.Id must be great than zero")
	}
	if !f.Visibly {
		tt.Fatal("after Save(), f.Visibly must be true")
	}
	if f.Created != f.Updated {
		tt.Fatal("after Save(), f.Created must equals f.Updated")
	}
	if !almostSameTime(f.Created, now, 1) {
		tt.Fatal("after Save(), f.Created must be now")
	}
	if !almostSameTime(f.Updated, now, 1) {
		tt.Fatal("after Save(), f.Updated must be now")
	}
	compareFailure(tt, f)
}

func randSetFailure(tt *testing.T, f *Failure) {
	f.Id = randInt32()
	f.Created = randTime()
	f.Updated = randTime()
	f.Visibly = randBool()
	f.FactoryId = randInt32()
	f.BatchId = randInt32()
	f.MacInt = randInt32()
	f.Mode = randEnumInt(FailureEnum.Mode.RMAP)
	f.Mac = randStr(0)
	f.Latest = randBool()
	f.raw = randStr(4096)
	f.IsFailed = randBool()
	f.FailureMode = randEnumInt(FailureEnum.FailureMode.RMAP)
	f.FailureMode1 = randEnumInt(FailureEnum.FailureMode1.RMAP)
	f.FailureMode2 = randEnumInt(FailureEnum.FailureMode2.RMAP)
	f.FailureMode3 = randEnumInt(FailureEnum.FailureMode3.RMAP)
}

func compareFailure(tt *testing.T, f *Failure) {
	f2 := FindFailure(nil, f.Id)
	if f.Id != f2.Id {
		tt.Fatal("insert and find compare failed, field: Id")
	}
	if !almostSameTime(f.Created, f2.Created, 1) {
		tt.Fatal("insert and find compare failed, field: Created")
	}
	if !almostSameTime(f.Updated, f2.Updated, 1) {
		tt.Fatal("insert and find compare failed, field: Updated")
	}
	if f.Visibly != f2.Visibly {
		tt.Fatal("insert and find compare failed, field: Visibly")
	}
	if f.FactoryId != f2.FactoryId {
		tt.Fatal("insert and find compare failed, field: FactoryId")
	}
	if f.BatchId != f2.BatchId {
		tt.Fatal("insert and find compare failed, field: BatchId")
	}
	if f.MacInt != f2.MacInt {
		tt.Fatal("insert and find compare failed, field: MacInt")
	}
	if f.Mode != f2.Mode {
		tt.Fatal("insert and find compare failed, field: Mode")
	}
	if f.Mac != f2.Mac {
		tt.Fatal("insert and find compare failed, field: Mac")
	}
	if f.Latest != f2.Latest {
		tt.Fatal("insert and find compare failed, field: Latest")
	}
	if f.raw != f2.raw {
		tt.Fatal("insert and find compare failed, field: raw")
	}
	if f.IsFailed != f2.IsFailed {
		tt.Fatal("insert and find compare failed, field: IsFailed")
	}
	if f.FailureMode != f2.FailureMode {
		tt.Fatal("insert and find compare failed, field: FailureMode")
	}
	if f.FailureMode1 != f2.FailureMode1 {
		tt.Fatal("insert and find compare failed, field: FailureMode1")
	}
	if f.FailureMode2 != f2.FailureMode2 {
		tt.Fatal("insert and find compare failed, field: FailureMode2")
	}
	if f.FailureMode3 != f2.FailureMode3 {
		tt.Fatal("insert and find compare failed, field: FailureMode3")
	}
}

func TestFailureUpdate(tt *testing.T) {
	f := NewFailure(nil)
	randSetFailure(tt, f)
	f.Id = 0
	f.Save()
	id := f.Id
	created := f.Created
	randSetFailure(tt, f)
	f.Id = id
	f.Visibly = true
	f.Created = created
	f.Update()
	if f.Created.After(f.Updated) || f.Created.Equal(f.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareFailure(tt, f)
}

func TestFailureInvisibly(tt *testing.T) {
	f := NewFailure(nil)
	randSetFailure(tt, f)
	f.Id = 0
	f.Save()
	f.Invisibly()
	f2 := FindFailure(nil, f.Id)
	if f2 != nil {
		tt.Fatal("after Invisibly, FindFailure() must return nil")
	}
}

func TestFailureDelete(tt *testing.T) {
	f := NewFailure(nil)
	randSetFailure(tt, f)
	f.Id = 0
	f.Save()
	f.Delete()
	f2 := FindFailure(nil, f.Id)
	if f2 != nil {
		tt.Fatal("after Invisibly, FindFailure() must return nil")
	}
}

func TestFailureValid(tt *testing.T) {
	f := NewFailure(nil)
	f.Mode = randInt()
	f.FailureMode = randInt()
	f.FailureMode1 = randInt()
	f.FailureMode2 = randInt()
	f.FailureMode3 = randInt()
	if f.Valid() == nil {
		tt.Fatal("enum using randInt must not be valid")
	}
}

func TestFailureUnmarshalMap(tt *testing.T) {
	f := NewFailure(nil)
	mm := make(map[string]interface{})
	id := randInt()
	created := randTime()
	updated := randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err := f.UnmarshalMap(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	if !compareFailureValue(mm["id"], f.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareFailureValue(mm["created"], f.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareFailureValue(mm["updated"], f.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = f.UnmarshalMap(nil, mm, FailureCol.Id)
	if !compareFailureValue(mm["id"], f.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareFailureValue(mm["created"], f.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareFailureValue(mm["updated"], f.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareFailureJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareFailureValue(jsonValue, fieldValue)
}

func compareFailureValue(jsonValue interface{}, fieldValue interface{}) bool {
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

func marshalAndUnmarshalFailure(tt *testing.T, f *Failure) map[string]interface{} {
	bs, err := json.Marshal(f)
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

func TestFailureExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	f := NewFailure(nil)
	f.SetExt(ext)
	randSetFailure(tt, f)
	mm := marshalAndUnmarshalFailure(tt, f)
	mmLen := len(mm)
	vv, ok := mm["created"].(string)
	if !ok {
		tt.Fatal("time type must be string in json")
	}
	if !strings.HasSuffix(vv, "+05:00") {
		// tt.Fatal("ext.loc has not affect")
	}
	f.ext.Verbose = "v"
	f.ext.IsComplex = true
	if _, ok := dalVerboses[FailureTid]; !ok {
		dalVerboses[FailureTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[FailureTid][f.ext.Verbose]
	dalVerboses[FailureTid][f.ext.Verbose] = []map[db.Col]interface{}{
		{FailureCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalFailure(tt, f)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareFailureValue(id, f.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[FailureTid][f.ext.Verbose] = []map[db.Col]interface{}{
		{}, {FailureCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalFailure(tt, f)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[FailureTid][f.ext.Verbose] = origin
}

func TestFailurePadding(tt *testing.T) {
	f := NewFailure(nil)
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
	f.Padding(key1, val1)
	f.Padding(key2, val2)
	f.Padding(key3, val3)
	f.Padding(key4, val4)
	f.Padding(key5, val5)
	mm := marshalAndUnmarshalFailure(tt, f)
	if !compareFailureValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareFailureValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareFailureValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareFailureValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareFailureValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestFailureMarshalJSON(tt *testing.T) {
	f := NewFailure(nil)
	randSetFailure(tt, f)
	mm := marshalAndUnmarshalFailure(tt, f)
	var jsonValue interface{}
	var jsonStr string
	var ok bool
	jsonValue = mm["id"]
	if !compareFailureValue(jsonValue, f.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareFailureValue(jsonValue, f.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareFailureValue(jsonValue, f.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareFailureValue(jsonValue, f.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["factory_id"]
	if !compareFailureValue(jsonValue, f.FactoryId) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactoryId) failed")
	}
	jsonValue = mm["batch_id"]
	if !compareFailureValue(jsonValue, f.BatchId) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchId) failed")
	}
	jsonValue = mm["mac_int"]
	if !compareFailureValue(jsonValue, f.MacInt) {
		tt.Fatal("json Marshal and Unmarshal compare field (MacInt) failed")
	}
	jsonValue = mm["mode"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.Mode.MAP[jsonStr], f.Mode) {
		tt.Fatal("json Marshal and Unmarshal compare field (Mode) failed")
	}
	jsonValue = mm["mac"]
	if !compareFailureValue(jsonValue, f.Mac) {
		tt.Fatal("json Marshal and Unmarshal compare field (Mac) failed")
	}
	jsonValue = mm["latest"]
	if !compareFailureValue(jsonValue, f.Latest) {
		tt.Fatal("json Marshal and Unmarshal compare field (Latest) failed")
	}
	jsonValue = mm["raw"]
	if !compareFailureValue(jsonValue, f.raw) {
		tt.Fatal("json Marshal and Unmarshal compare field (raw) failed")
	}
	jsonValue = mm["is_failed"]
	if !compareFailureValue(jsonValue, f.IsFailed) {
		tt.Fatal("json Marshal and Unmarshal compare field (IsFailed) failed")
	}
	jsonValue = mm["failure_mode"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode.MAP[jsonStr], f.FailureMode) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode) failed")
	}
	jsonValue = mm["failure_mode1"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode1.MAP[jsonStr], f.FailureMode1) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode1) failed")
	}
	jsonValue = mm["failure_mode2"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode2.MAP[jsonStr], f.FailureMode2) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode2) failed")
	}
	jsonValue = mm["failure_mode3"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode3.MAP[jsonStr], f.FailureMode3) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode3) failed")
	}
}

func TestFailureMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	f := NewFailure(nil)
	f.SetExt(ext)
	randSetFailure(tt, f)
	origin := dalVerboses[FailureTid][f.ext.Verbose]
	dalVerboses[FailureTid][f.ext.Verbose] = []map[db.Col]interface{}{
		{}, {FailureCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalFailure(tt, f)
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
	f.Padding("id", rint)
	f.Padding("updated", rstr)
	mm = marshalAndUnmarshalFailure(tt, f)
	if !compareFailureValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareFailureValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[FailureTid][f.ext.Verbose] = origin
}

func TestFailureUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["factory_id"] = randInt32()
	mm["batch_id"] = randInt32()
	mm["mac_int"] = randInt32()
	mm["mode"] = randEnumStr(FailureEnum.Mode.RMAP)
	mm["mac"] = randStr(0)
	mm["latest"] = randBool()
	mm["raw"] = randStr(4096)
	mm["is_failed"] = randBool()
	mm["failure_mode"] = randEnumStr(FailureEnum.FailureMode.RMAP)
	mm["failure_mode1"] = randEnumStr(FailureEnum.FailureMode1.RMAP)
	mm["failure_mode2"] = randEnumStr(FailureEnum.FailureMode2.RMAP)
	mm["failure_mode3"] = randEnumStr(FailureEnum.FailureMode3.RMAP)
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	f, err := UnmarshalFailure(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	var jsonStr string
	var ok bool
	jsonValue = mm["id"]
	if !compareFailureValue(jsonValue, f.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareFailureValue(jsonValue, f.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareFailureValue(jsonValue, f.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareFailureValue(jsonValue, f.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["factory_id"]
	if !compareFailureValue(jsonValue, f.FactoryId) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactoryId) failed")
	}
	jsonValue = mm["batch_id"]
	if !compareFailureValue(jsonValue, f.BatchId) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchId) failed")
	}
	jsonValue = mm["mac_int"]
	if !compareFailureValue(jsonValue, f.MacInt) {
		tt.Fatal("json Marshal and Unmarshal compare field (MacInt) failed")
	}
	jsonValue = mm["mode"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.Mode.MAP[jsonStr], f.Mode) {
		tt.Fatal("json Marshal and Unmarshal compare field (Mode) failed")
	}
	jsonValue = mm["mac"]
	if !compareFailureValue(jsonValue, f.Mac) {
		tt.Fatal("json Marshal and Unmarshal compare field (Mac) failed")
	}
	jsonValue = mm["latest"]
	if !compareFailureValue(jsonValue, f.Latest) {
		tt.Fatal("json Marshal and Unmarshal compare field (Latest) failed")
	}
	jsonValue = mm["raw"]
	if !compareFailureValue(jsonValue, f.raw) {
		tt.Fatal("json Marshal and Unmarshal compare field (raw) failed")
	}
	jsonValue = mm["is_failed"]
	if !compareFailureValue(jsonValue, f.IsFailed) {
		tt.Fatal("json Marshal and Unmarshal compare field (IsFailed) failed")
	}
	jsonValue = mm["failure_mode"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode.MAP[jsonStr], f.FailureMode) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode) failed")
	}
	jsonValue = mm["failure_mode1"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode1.MAP[jsonStr], f.FailureMode1) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode1) failed")
	}
	jsonValue = mm["failure_mode2"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode2.MAP[jsonStr], f.FailureMode2) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode2) failed")
	}
	jsonValue = mm["failure_mode3"]
	if jsonStr, ok = jsonValue.(string); !ok {
		tt.Fatal("enum must be string type")
	}
	if !compareFailureValue(FailureEnum.FailureMode3.MAP[jsonStr], f.FailureMode3) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailureMode3) failed")
	}
}

package dal

import (
	"encoding/json"
	"espressif.com/chip/factory/db"
	"log"
	"strings"
	"testing"
	"time"
)

func TestFactorySave(tt *testing.T) {
	f := NewFactory(nil)
	randSetFactory(tt, f)
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
	compareFactory(tt, f)
}

func randSetFactory(tt *testing.T, f *Factory) {
	f.Id = randInt32()
	f.Created = randTime()
	f.Updated = randTime()
	f.Visibly = randBool()
	f.Sid = randStr(64)
	f.Name = randStr(64)
	f.Location = randStr(64)
	f.Token = randStr(64)
	f.IsStaff = randBool()
}

func compareFactory(tt *testing.T, f *Factory) {
	f2 := FindFactory(nil, f.Id)
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
	if f.Sid != f2.Sid {
		tt.Fatal("insert and find compare failed, field: Sid")
	}
	if f.Name != f2.Name {
		tt.Fatal("insert and find compare failed, field: Name")
	}
	if f.Location != f2.Location {
		tt.Fatal("insert and find compare failed, field: Location")
	}
	if f.Token != f2.Token {
		tt.Fatal("insert and find compare failed, field: Token")
	}
	if f.IsStaff != f2.IsStaff {
		tt.Fatal("insert and find compare failed, field: IsStaff")
	}
}

func TestFactoryUpdate(tt *testing.T) {
	f := NewFactory(nil)
	randSetFactory(tt, f)
	f.Id = 0
	f.Save()
	id := f.Id
	created := f.Created
	randSetFactory(tt, f)
	f.Id = id
	f.Visibly = true
	f.Created = created
	f.Update()
	if f.Created.After(f.Updated) || f.Created.Equal(f.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareFactory(tt, f)
}

func TestFactoryInvisibly(tt *testing.T) {
	f := NewFactory(nil)
	randSetFactory(tt, f)
	f.Id = 0
	f.Save()
	f.Invisibly()
	f2 := FindFactory(nil, f.Id)
	if f2 != nil {
		tt.Fatal("after Invisibly, FindFactory() must return nil")
	}
}

func TestFactoryDelete(tt *testing.T) {
	f := NewFactory(nil)
	randSetFactory(tt, f)
	f.Id = 0
	f.Save()
	f.Delete()
	f2 := FindFactory(nil, f.Id)
	if f2 != nil {
		tt.Fatal("after Invisibly, FindFactory() must return nil")
	}
}

func TestFactoryUnmarshalMap(tt *testing.T) {
	f := NewFactory(nil)
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
	if !compareFactoryValue(mm["id"], f.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareFactoryValue(mm["created"], f.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareFactoryValue(mm["updated"], f.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = f.UnmarshalMap(nil, mm, FactoryCol.Id)
	if !compareFactoryValue(mm["id"], f.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareFactoryValue(mm["created"], f.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareFactoryValue(mm["updated"], f.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareFactoryJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareFactoryValue(jsonValue, fieldValue)
}

func compareFactoryValue(jsonValue interface{}, fieldValue interface{}) bool {
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

func marshalAndUnmarshalFactory(tt *testing.T, f *Factory) map[string]interface{} {
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

func TestFactoryExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	f := NewFactory(nil)
	f.SetExt(ext)
	randSetFactory(tt, f)
	mm := marshalAndUnmarshalFactory(tt, f)
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
	if _, ok := dalVerboses[FactoryTid]; !ok {
		dalVerboses[FactoryTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[FactoryTid][f.ext.Verbose]
	dalVerboses[FactoryTid][f.ext.Verbose] = []map[db.Col]interface{}{
		{FactoryCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalFactory(tt, f)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareFactoryValue(id, f.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[FactoryTid][f.ext.Verbose] = []map[db.Col]interface{}{
		{}, {FactoryCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalFactory(tt, f)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[FactoryTid][f.ext.Verbose] = origin
}

func TestFactoryPadding(tt *testing.T) {
	f := NewFactory(nil)
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
	mm := marshalAndUnmarshalFactory(tt, f)
	if !compareFactoryValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareFactoryValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareFactoryValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareFactoryValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareFactoryValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestFactoryMarshalJSON(tt *testing.T) {
	f := NewFactory(nil)
	randSetFactory(tt, f)
	mm := marshalAndUnmarshalFactory(tt, f)
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareFactoryValue(jsonValue, f.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareFactoryValue(jsonValue, f.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareFactoryValue(jsonValue, f.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareFactoryValue(jsonValue, f.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["sid"]
	if !compareFactoryValue(jsonValue, f.Sid) {
		tt.Fatal("json Marshal and Unmarshal compare field (Sid) failed")
	}
	jsonValue = mm["name"]
	if !compareFactoryValue(jsonValue, f.Name) {
		tt.Fatal("json Marshal and Unmarshal compare field (Name) failed")
	}
	jsonValue = mm["location"]
	if !compareFactoryValue(jsonValue, f.Location) {
		tt.Fatal("json Marshal and Unmarshal compare field (Location) failed")
	}
	jsonValue = mm["token"]
	if !compareFactoryValue(jsonValue, f.Token) {
		tt.Fatal("json Marshal and Unmarshal compare field (Token) failed")
	}
	jsonValue = mm["is_staff"]
	if !compareFactoryValue(jsonValue, f.IsStaff) {
		tt.Fatal("json Marshal and Unmarshal compare field (IsStaff) failed")
	}
}

func TestFactoryMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	f := NewFactory(nil)
	f.SetExt(ext)
	randSetFactory(tt, f)
	origin := dalVerboses[FactoryTid][f.ext.Verbose]
	dalVerboses[FactoryTid][f.ext.Verbose] = []map[db.Col]interface{}{
		{}, {FactoryCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalFactory(tt, f)
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
	mm = marshalAndUnmarshalFactory(tt, f)
	if !compareFactoryValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareFactoryValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[FactoryTid][f.ext.Verbose] = origin
}

func TestFactoryUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["sid"] = randStr(64)
	mm["name"] = randStr(64)
	mm["location"] = randStr(64)
	mm["token"] = randStr(64)
	mm["is_staff"] = randBool()
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	f, err := UnmarshalFactory(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareFactoryValue(jsonValue, f.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareFactoryValue(jsonValue, f.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareFactoryValue(jsonValue, f.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareFactoryValue(jsonValue, f.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["sid"]
	if !compareFactoryValue(jsonValue, f.Sid) {
		tt.Fatal("json Marshal and Unmarshal compare field (Sid) failed")
	}
	jsonValue = mm["name"]
	if !compareFactoryValue(jsonValue, f.Name) {
		tt.Fatal("json Marshal and Unmarshal compare field (Name) failed")
	}
	jsonValue = mm["location"]
	if !compareFactoryValue(jsonValue, f.Location) {
		tt.Fatal("json Marshal and Unmarshal compare field (Location) failed")
	}
	jsonValue = mm["token"]
	if !compareFactoryValue(jsonValue, f.Token) {
		tt.Fatal("json Marshal and Unmarshal compare field (Token) failed")
	}
	jsonValue = mm["is_staff"]
	if !compareFactoryValue(jsonValue, f.IsStaff) {
		tt.Fatal("json Marshal and Unmarshal compare field (IsStaff) failed")
	}
}

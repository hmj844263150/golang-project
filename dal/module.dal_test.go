package dal

import (
	"encoding/json"
	"espressif.com/chip/factory/db"
	"log"
	"strings"
	"testing"
	"time"
)

func TestModuleSave(tt *testing.T) {
	m := NewModule(nil)
	randSetModule(tt, m)
	now := time.Now()
	m.Id = 0
	m.Save()
	if m.Id <= 0 {
		log.Println("save error?")
		tt.Fatal("after Save(), m.Id must be great than zero")
	}
	if !m.Visibly {
		tt.Fatal("after Save(), m.Visibly must be true")
	}
	if m.Created != m.Updated {
		tt.Fatal("after Save(), m.Created must equals m.Updated")
	}
	if !almostSameTime(m.Created, now, 1) {
		tt.Fatal("after Save(), m.Created must be now")
	}
	if !almostSameTime(m.Updated, now, 1) {
		tt.Fatal("after Save(), m.Updated must be now")
	}
	compareModule(tt, m)
}

func randSetModule(tt *testing.T, m *Module) {
	m.Id = randInt32()
	m.Created = randTime()
	m.Updated = randTime()
	m.Visibly = randBool()
	m.EspMac = randStr(64)
}

func compareModule(tt *testing.T, m *Module) {
	m2 := FindModule(nil, m.Id)
	if m.Id != m2.Id {
		tt.Fatal("insert and find compare failed, field: Id")
	}
	if !almostSameTime(m.Created, m2.Created, 1) {
		tt.Fatal("insert and find compare failed, field: Created")
	}
	if !almostSameTime(m.Updated, m2.Updated, 1) {
		tt.Fatal("insert and find compare failed, field: Updated")
	}
	if m.Visibly != m2.Visibly {
		tt.Fatal("insert and find compare failed, field: Visibly")
	}
	if m.EspMac != m2.EspMac {
		tt.Fatal("insert and find compare failed, field: EspMac")
	}
}

func TestModuleUpdate(tt *testing.T) {
	m := NewModule(nil)
	randSetModule(tt, m)
	m.Id = 0
	m.Save()
	id := m.Id
	created := m.Created
	randSetModule(tt, m)
	m.Id = id
	m.Visibly = true
	m.Created = created
	m.Update()
	if m.Created.After(m.Updated) || m.Created.Equal(m.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareModule(tt, m)
}

func TestModuleInvisibly(tt *testing.T) {
	m := NewModule(nil)
	randSetModule(tt, m)
	m.Id = 0
	m.Save()
	m.Invisibly()
	m2 := FindModule(nil, m.Id)
	if m2 != nil {
		tt.Fatal("after Invisibly, FindModule() must return nil")
	}
}

func TestModuleDelete(tt *testing.T) {
	m := NewModule(nil)
	randSetModule(tt, m)
	m.Id = 0
	m.Save()
	m.Delete()
	m2 := FindModule(nil, m.Id)
	if m2 != nil {
		tt.Fatal("after Invisibly, FindModule() must return nil")
	}
}

func TestModuleUnmarshalMap(tt *testing.T) {
	m := NewModule(nil)
	mm := make(map[string]interface{})
	id := randInt()
	created := randTime()
	updated := randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err := m.UnmarshalMap(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	if !compareModuleValue(mm["id"], m.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareModuleValue(mm["created"], m.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareModuleValue(mm["updated"], m.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = m.UnmarshalMap(nil, mm, ModuleCol.Id)
	if !compareModuleValue(mm["id"], m.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareModuleValue(mm["created"], m.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareModuleValue(mm["updated"], m.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareModuleJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareModuleValue(jsonValue, fieldValue)
}

func compareModuleValue(jsonValue interface{}, fieldValue interface{}) bool {
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

func marshalAndUnmarshalModule(tt *testing.T, m *Module) map[string]interface{} {
	bs, err := json.Marshal(m)
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

func TestModuleExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	m := NewModule(nil)
	m.SetExt(ext)
	randSetModule(tt, m)
	mm := marshalAndUnmarshalModule(tt, m)
	mmLen := len(mm)
	vv, ok := mm["created"].(string)
	if !ok {
		tt.Fatal("time type must be string in json")
	}
	if !strings.HasSuffix(vv, "+05:00") {
		// tt.Fatal("ext.loc has not affect")
	}
	m.ext.Verbose = "v"
	m.ext.IsComplex = true
	if _, ok := dalVerboses[ModuleTid]; !ok {
		dalVerboses[ModuleTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[ModuleTid][m.ext.Verbose]
	dalVerboses[ModuleTid][m.ext.Verbose] = []map[db.Col]interface{}{
		{ModuleCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalModule(tt, m)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareModuleValue(id, m.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[ModuleTid][m.ext.Verbose] = []map[db.Col]interface{}{
		{}, {ModuleCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalModule(tt, m)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[ModuleTid][m.ext.Verbose] = origin
}

func TestModulePadding(tt *testing.T) {
	m := NewModule(nil)
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
	m.Padding(key1, val1)
	m.Padding(key2, val2)
	m.Padding(key3, val3)
	m.Padding(key4, val4)
	m.Padding(key5, val5)
	mm := marshalAndUnmarshalModule(tt, m)
	if !compareModuleValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareModuleValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareModuleValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareModuleValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareModuleValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestModuleMarshalJSON(tt *testing.T) {
	m := NewModule(nil)
	randSetModule(tt, m)
	mm := marshalAndUnmarshalModule(tt, m)
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareModuleValue(jsonValue, m.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareModuleValue(jsonValue, m.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareModuleValue(jsonValue, m.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareModuleValue(jsonValue, m.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["esp_mac"]
	if !compareModuleValue(jsonValue, m.EspMac) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMac) failed")
	}
}

func TestModuleMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	m := NewModule(nil)
	m.SetExt(ext)
	randSetModule(tt, m)
	origin := dalVerboses[ModuleTid][m.ext.Verbose]
	dalVerboses[ModuleTid][m.ext.Verbose] = []map[db.Col]interface{}{
		{}, {ModuleCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalModule(tt, m)
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
	m.Padding("id", rint)
	m.Padding("updated", rstr)
	mm = marshalAndUnmarshalModule(tt, m)
	if !compareModuleValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareModuleValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[ModuleTid][m.ext.Verbose] = origin
}

func TestModuleUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["esp_mac"] = randStr(64)
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	m, err := UnmarshalModule(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareModuleValue(jsonValue, m.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareModuleValue(jsonValue, m.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareModuleValue(jsonValue, m.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareModuleValue(jsonValue, m.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["esp_mac"]
	if !compareModuleValue(jsonValue, m.EspMac) {
		tt.Fatal("json Marshal and Unmarshal compare field (EspMac) failed")
	}
}

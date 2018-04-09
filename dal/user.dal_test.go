package dal

import (
	"encoding/json"
	"testing"
	"time"
	"log"
	"strings"
	"espressif.com/chip/factory/db"
)

func TestUserSave(tt *testing.T) {
	u := NewUser(nil)
	randSetUser(tt, u)
	now := time.Now()
	u.Id = 0
	u.Save()
	if u.Id <= 0 {
		log.Println("save error?")
		tt.Fatal("after Save(), u.Id must be great than zero")
	}
	if !u.Visibly {
		tt.Fatal("after Save(), u.Visibly must be true")
	}
	if u.Created != u.Updated {
		tt.Fatal("after Save(), u.Created must equals u.Updated")
	}
	if !almostSameTime(u.Created, now, 1) {
		tt.Fatal("after Save(), u.Created must be now")
	}
	if !almostSameTime(u.Updated, now, 1) {
		tt.Fatal("after Save(), u.Updated must be now")
	}
	compareUser(tt, u)
}

func randSetUser(tt *testing.T, u *User) {
	u.Id = randInt32()
	u.Created = randTime()
	u.Updated = randTime()
	u.Visibly = randBool()
	u.Account = randStr(64)
	u.Password = randStr(64)
	u.Name = randStr(64)
	u.FactorySid = randStr(64)
	u.GroupId = randInt32()
	u.Email = randStr(128)
	u.Description = randStr(256)
}

func compareUser(tt *testing.T, u *User) {
	u2 := FindUser(nil, u.Id)
	if u.Id != u2.Id {
		tt.Fatal("insert and find compare failed, field: Id")
	}
	if !almostSameTime(u.Created, u2.Created, 1) {
		tt.Fatal("insert and find compare failed, field: Created")
	}
	if !almostSameTime(u.Updated, u2.Updated, 1) {
		tt.Fatal("insert and find compare failed, field: Updated")
	}
	if u.Visibly != u2.Visibly {
		tt.Fatal("insert and find compare failed, field: Visibly")
	}
	if u.Account != u2.Account {
		tt.Fatal("insert and find compare failed, field: Account")
	}
	if u.Password != u2.Password {
		tt.Fatal("insert and find compare failed, field: Password")
	}
	if u.Name != u2.Name {
		tt.Fatal("insert and find compare failed, field: Name")
	}
	if u.FactorySid != u2.FactorySid {
		tt.Fatal("insert and find compare failed, field: FactorySid")
	}
	if u.GroupId != u2.GroupId {
		tt.Fatal("insert and find compare failed, field: GroupId")
	}
	if u.Email != u2.Email {
		tt.Fatal("insert and find compare failed, field: Email")
	}
	if u.Description != u2.Description {
		tt.Fatal("insert and find compare failed, field: Description")
	}
}

func TestUserUpdate(tt *testing.T) {
	u := NewUser(nil)
	randSetUser(tt, u)
	u.Id = 0
	u.Save()
	id := u.Id
	created := u.Created
	randSetUser(tt, u)
	u.Id = id
	u.Visibly = true
	u.Created = created
	u.Update()
	if u.Created.After(u.Updated) || u.Created.Equal(u.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareUser(tt, u)
}

func TestUserInvisibly(tt *testing.T) {
	u := NewUser(nil)
	randSetUser(tt, u)
	u.Id = 0
	u.Save()
	u.Invisibly()
	u2 := FindUser(nil, u.Id)
	if u2 != nil {
		tt.Fatal("after Invisibly, FindUser() must return nil")
	}
}

func TestUserDelete(tt *testing.T) {
	u := NewUser(nil)
	randSetUser(tt, u)
	u.Id = 0
	u.Save()
	u.Delete()
	u2 := FindUser(nil, u.Id)
	if u2 != nil {
		tt.Fatal("after Invisibly, FindUser() must return nil")
	}
}

func TestUserUnmarshalMap(tt *testing.T) {
	u := NewUser(nil)
	mm := make(map[string]interface{})
	id := randInt()
	created := randTime()
	updated := randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err := u.UnmarshalMap(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	if !compareUserValue(mm["id"], u.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareUserValue(mm["created"], u.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareUserValue(mm["updated"], u.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = u.UnmarshalMap(nil, mm, UserCol.Id)
	if !compareUserValue(mm["id"], u.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareUserValue(mm["created"], u.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareUserValue(mm["updated"], u.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareUserJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareUserValue(jsonValue, fieldValue)
}

func compareUserValue(jsonValue interface{}, fieldValue interface{}) bool {
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

func marshalAndUnmarshalUser(tt *testing.T, u *User) map[string]interface{} {
	bs, err := json.Marshal(u)
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

func TestUserExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	u := NewUser(nil)
	u.SetExt(ext)
	randSetUser(tt, u)
	mm := marshalAndUnmarshalUser(tt, u)
	mmLen := len(mm)
	vv, ok := mm["created"].(string)
	if !ok {
		tt.Fatal("time type must be string in json")
	}
	if !strings.HasSuffix(vv, "+05:00") {
		// tt.Fatal("ext.loc has not affect")
	}
	u.ext.Verbose = "v"
	u.ext.IsComplex = true
	if _, ok := dalVerboses[UserTid]; !ok {
		dalVerboses[UserTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[UserTid][u.ext.Verbose]
	dalVerboses[UserTid][u.ext.Verbose] = []map[db.Col]interface{}{
		{UserCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalUser(tt, u)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareUserValue(id, u.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[UserTid][u.ext.Verbose] = []map[db.Col]interface{}{
		{}, {UserCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalUser(tt, u)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[UserTid][u.ext.Verbose] = origin
}

func TestUserPadding(tt *testing.T) {
	u := NewUser(nil)
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
	u.Padding(key1, val1)
	u.Padding(key2, val2)
	u.Padding(key3, val3)
	u.Padding(key4, val4)
	u.Padding(key5, val5)
	mm := marshalAndUnmarshalUser(tt, u)
	if !compareUserValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareUserValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareUserValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareUserValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareUserValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestUserMarshalJSON(tt *testing.T) {
	u := NewUser(nil)
	randSetUser(tt, u)
	mm := marshalAndUnmarshalUser(tt, u)
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareUserValue(jsonValue, u.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareUserValue(jsonValue, u.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareUserValue(jsonValue, u.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareUserValue(jsonValue, u.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["account"]
	if !compareUserValue(jsonValue, u.Account) {
		tt.Fatal("json Marshal and Unmarshal compare field (Account) failed")
	}
	jsonValue = mm["password"]
	if !compareUserValue(jsonValue, u.Password) {
		tt.Fatal("json Marshal and Unmarshal compare field (Password) failed")
	}
	jsonValue = mm["name"]
	if !compareUserValue(jsonValue, u.Name) {
		tt.Fatal("json Marshal and Unmarshal compare field (Name) failed")
	}
	jsonValue = mm["factory_sid"]
	if !compareUserValue(jsonValue, u.FactorySid) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactorySid) failed")
	}
	jsonValue = mm["group_id"]
	if !compareUserValue(jsonValue, u.GroupId) {
		tt.Fatal("json Marshal and Unmarshal compare field (GroupId) failed")
	}
	jsonValue = mm["email"]
	if !compareUserValue(jsonValue, u.Email) {
		tt.Fatal("json Marshal and Unmarshal compare field (Email) failed")
	}
	jsonValue = mm["description"]
	if !compareUserValue(jsonValue, u.Description) {
		tt.Fatal("json Marshal and Unmarshal compare field (Description) failed")
	}
}

func TestUserMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	u := NewUser(nil)
	u.SetExt(ext)
	randSetUser(tt, u)
	origin := dalVerboses[UserTid][u.ext.Verbose]
	dalVerboses[UserTid][u.ext.Verbose] = []map[db.Col]interface{}{
		{}, {UserCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalUser(tt, u)
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
	u.Padding("id", rint)
	u.Padding("updated", rstr)
	mm = marshalAndUnmarshalUser(tt, u)
	if !compareUserValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareUserValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[UserTid][u.ext.Verbose] = origin
}

func TestUserUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["account"] = randStr(64)
	mm["password"] = randStr(64)
	mm["name"] = randStr(64)
	mm["factory_sid"] = randStr(64)
	mm["group_id"] = randInt32()
	mm["email"] = randStr(128)
	mm["description"] = randStr(256)
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	u, err := UnmarshalUser(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareUserValue(jsonValue, u.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareUserValue(jsonValue, u.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareUserValue(jsonValue, u.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareUserValue(jsonValue, u.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["account"]
	if !compareUserValue(jsonValue, u.Account) {
		tt.Fatal("json Marshal and Unmarshal compare field (Account) failed")
	}
	jsonValue = mm["password"]
	if !compareUserValue(jsonValue, u.Password) {
		tt.Fatal("json Marshal and Unmarshal compare field (Password) failed")
	}
	jsonValue = mm["name"]
	if !compareUserValue(jsonValue, u.Name) {
		tt.Fatal("json Marshal and Unmarshal compare field (Name) failed")
	}
	jsonValue = mm["factory_sid"]
	if !compareUserValue(jsonValue, u.FactorySid) {
		tt.Fatal("json Marshal and Unmarshal compare field (FactorySid) failed")
	}
	jsonValue = mm["group_id"]
	if !compareUserValue(jsonValue, u.GroupId) {
		tt.Fatal("json Marshal and Unmarshal compare field (GroupId) failed")
	}
	jsonValue = mm["email"]
	if !compareUserValue(jsonValue, u.Email) {
		tt.Fatal("json Marshal and Unmarshal compare field (Email) failed")
	}
	jsonValue = mm["description"]
	if !compareUserValue(jsonValue, u.Description) {
		tt.Fatal("json Marshal and Unmarshal compare field (Description) failed")
	}
}

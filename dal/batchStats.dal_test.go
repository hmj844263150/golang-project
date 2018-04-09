package dal

import (
	"encoding/json"
	"espressif.com/chip/factory/db"
	"log"
	"strings"
	"testing"
	"time"
)

func TestBatchStatsSave(tt *testing.T) {
	b := NewBatchStats(nil)
	randSetBatchStats(tt, b)
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
	compareBatchStats(tt, b)
}

func randSetBatchStats(tt *testing.T, b *BatchStats) {
	b.Id = randInt32()
	b.Created = randTime()
	b.Updated = randTime()
	b.Visibly = randBool()
	b.BatchId = randInt32()
	b.Start = randTime()
	b.End = randTime()
	b.Cnt = randInt32()
	b.Success = randInt32()
	b.SuccessPct = randInt32()
	b.RightFirstTime = randInt32()
	b.RightFirstTimePct = randInt32()
	b.Failed = randInt32()
	b.FailedPct = randInt32()
	b.Rejected = randInt32()
	b.RejectedPct = randInt32()
}

func compareBatchStats(tt *testing.T, b *BatchStats) {
	b2 := FindBatchStats(nil, b.Id)
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
	if b.BatchId != b2.BatchId {
		tt.Fatal("insert and find compare failed, field: BatchId")
	}
	if !almostSameTime(b.Start, b2.Start, 1) {
		tt.Fatal("insert and find compare failed, field: Start")
	}
	if !almostSameTime(b.End, b2.End, 1) {
		tt.Fatal("insert and find compare failed, field: End")
	}
	if b.Cnt != b2.Cnt {
		tt.Fatal("insert and find compare failed, field: Cnt")
	}
	if b.Success != b2.Success {
		tt.Fatal("insert and find compare failed, field: Success")
	}
	if b.SuccessPct != b2.SuccessPct {
		tt.Fatal("insert and find compare failed, field: SuccessPct")
	}
	if b.RightFirstTime != b2.RightFirstTime {
		tt.Fatal("insert and find compare failed, field: RightFirstTime")
	}
	if b.RightFirstTimePct != b2.RightFirstTimePct {
		tt.Fatal("insert and find compare failed, field: RightFirstTimePct")
	}
	if b.Failed != b2.Failed {
		tt.Fatal("insert and find compare failed, field: Failed")
	}
	if b.FailedPct != b2.FailedPct {
		tt.Fatal("insert and find compare failed, field: FailedPct")
	}
	if b.Rejected != b2.Rejected {
		tt.Fatal("insert and find compare failed, field: Rejected")
	}
	if b.RejectedPct != b2.RejectedPct {
		tt.Fatal("insert and find compare failed, field: RejectedPct")
	}
}

func TestBatchStatsUpdate(tt *testing.T) {
	b := NewBatchStats(nil)
	randSetBatchStats(tt, b)
	b.Id = 0
	b.Save()
	id := b.Id
	created := b.Created
	randSetBatchStats(tt, b)
	b.Id = id
	b.Visibly = true
	b.Created = created
	b.Update()
	if b.Created.After(b.Updated) || b.Created.Equal(b.Updated) {
		tt.Fatal("after update, Updated must be great than Created")
	}
	compareBatchStats(tt, b)
}

func TestBatchStatsInvisibly(tt *testing.T) {
	b := NewBatchStats(nil)
	randSetBatchStats(tt, b)
	b.Id = 0
	b.Save()
	b.Invisibly()
	b2 := FindBatchStats(nil, b.Id)
	if b2 != nil {
		tt.Fatal("after Invisibly, FindBatchStats() must return nil")
	}
}

func TestBatchStatsDelete(tt *testing.T) {
	b := NewBatchStats(nil)
	randSetBatchStats(tt, b)
	b.Id = 0
	b.Save()
	b.Delete()
	b2 := FindBatchStats(nil, b.Id)
	if b2 != nil {
		tt.Fatal("after Invisibly, FindBatchStats() must return nil")
	}
}

func TestBatchStatsUnmarshalMap(tt *testing.T) {
	b := NewBatchStats(nil)
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
	if !compareBatchStatsValue(mm["id"], b.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if !compareBatchStatsValue(mm["created"], b.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if !compareBatchStatsValue(mm["updated"], b.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
	id = randInt()
	created = randTime()
	updated = randTime()
	mm["id"] = id
	mm["created"] = created
	mm["updated"] = updated
	_, err = b.UnmarshalMap(nil, mm, BatchStatsCol.Id)
	if !compareBatchStatsValue(mm["id"], b.Id) {
		tt.Fatal("UnmarshalMap failed id")
	}
	if compareBatchStatsValue(mm["created"], b.Created) {
		tt.Fatal("UnmarshalMap failed created")
	}
	if compareBatchStatsValue(mm["updated"], b.Updated) {
		tt.Fatal("UnmarshalMap failed updated")
	}
}

func compareBatchStatsJsonField(jsonb []byte, field string, fieldValue interface{}) bool {
	mm := make(map[string]interface{})
	err := json.Unmarshal(jsonb, &mm)
	if err != nil {
		return false
	}
	jsonValue := mm[field]
	return compareBatchStatsValue(jsonValue, fieldValue)
}

func compareBatchStatsValue(jsonValue interface{}, fieldValue interface{}) bool {
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

func marshalAndUnmarshalBatchStats(tt *testing.T, b *BatchStats) map[string]interface{} {
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

func TestBatchStatsExt(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc}
	b := NewBatchStats(nil)
	b.SetExt(ext)
	randSetBatchStats(tt, b)
	mm := marshalAndUnmarshalBatchStats(tt, b)
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
	if _, ok := dalVerboses[BatchStatsTid]; !ok {
		dalVerboses[BatchStatsTid] = map[string][]map[db.Col]interface{}{"v": nil}
	}
	origin := dalVerboses[BatchStatsTid][b.ext.Verbose]
	dalVerboses[BatchStatsTid][b.ext.Verbose] = []map[db.Col]interface{}{
		{BatchStatsCol.Id: struct{}{}}, {},
	}
	mm = marshalAndUnmarshalBatchStats(tt, b)
	if len(mm) != 1 {
		tt.Fatal("ext.includes only include id field, len(mm) != 1")
	}
	id, ok := mm["id"]
	if !ok {
		tt.Fatal("ext.includes only include id field, id, ok := mm[\"id\"]")
	}
	if !compareBatchStatsValue(id, b.Id) {
		tt.Fatal("ext.includes compare failed")
	}
	dalVerboses[BatchStatsTid][b.ext.Verbose] = []map[db.Col]interface{}{
		{}, {BatchStatsCol.Id: struct{}{}},
	}
	mm = marshalAndUnmarshalBatchStats(tt, b)
	if len(mm) != (mmLen - 1) {
		tt.Fatal("ext.excludes only exclude id field, len(mm) != (mmLen - 1)")
	}
	_, ok = mm["id"]
	if ok {
		tt.Fatal("ext.excludes only exclude id field, _, ok := mm[\"id\"]")
	}
	dalVerboses[BatchStatsTid][b.ext.Verbose] = origin
}

func TestBatchStatsPadding(tt *testing.T) {
	b := NewBatchStats(nil)
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
	mm := marshalAndUnmarshalBatchStats(tt, b)
	if !compareBatchStatsValue(mm[key1], val1) {
		tt.Fatal("Padding() string compare failed")
	}
	if !compareBatchStatsValue(mm[key2], val2) {
		tt.Fatal("Padding() float64 compare failed")
	}
	if !compareBatchStatsValue(mm[key3], val3) {
		tt.Fatal("Padding() int compare failed")
	}
	if !compareBatchStatsValue(mm[key4], val4) {
		tt.Fatal("Padding() bool compare failed")
	}
	if !compareBatchStatsValue(mm[key5], val5) {
		tt.Fatal("Padding() time compare failed")
	}
}

func TestBatchStatsMarshalJSON(tt *testing.T) {
	b := NewBatchStats(nil)
	randSetBatchStats(tt, b)
	mm := marshalAndUnmarshalBatchStats(tt, b)
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareBatchStatsValue(jsonValue, b.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareBatchStatsValue(jsonValue, b.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareBatchStatsValue(jsonValue, b.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareBatchStatsValue(jsonValue, b.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["batch_id"]
	if !compareBatchStatsValue(jsonValue, b.BatchId) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchId) failed")
	}
	jsonValue = mm["start"]
	if !compareBatchStatsValue(jsonValue, b.Start) {
		tt.Fatal("json Marshal and Unmarshal compare field (Start) failed")
	}
	jsonValue = mm["end"]
	if !compareBatchStatsValue(jsonValue, b.End) {
		tt.Fatal("json Marshal and Unmarshal compare field (End) failed")
	}
	jsonValue = mm["cnt"]
	if !compareBatchStatsValue(jsonValue, b.Cnt) {
		tt.Fatal("json Marshal and Unmarshal compare field (Cnt) failed")
	}
	jsonValue = mm["success"]
	if !compareBatchStatsValue(jsonValue, b.Success) {
		tt.Fatal("json Marshal and Unmarshal compare field (Success) failed")
	}
	jsonValue = mm["success_pct"]
	if !compareBatchStatsValue(jsonValue, b.SuccessPct) {
		tt.Fatal("json Marshal and Unmarshal compare field (SuccessPct) failed")
	}
	jsonValue = mm["right_first_time"]
	if !compareBatchStatsValue(jsonValue, b.RightFirstTime) {
		tt.Fatal("json Marshal and Unmarshal compare field (RightFirstTime) failed")
	}
	jsonValue = mm["right_first_time_pct"]
	if !compareBatchStatsValue(jsonValue, b.RightFirstTimePct) {
		tt.Fatal("json Marshal and Unmarshal compare field (RightFirstTimePct) failed")
	}
	jsonValue = mm["failed"]
	if !compareBatchStatsValue(jsonValue, b.Failed) {
		tt.Fatal("json Marshal and Unmarshal compare field (Failed) failed")
	}
	jsonValue = mm["failed_pct"]
	if !compareBatchStatsValue(jsonValue, b.FailedPct) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailedPct) failed")
	}
	jsonValue = mm["rejected"]
	if !compareBatchStatsValue(jsonValue, b.Rejected) {
		tt.Fatal("json Marshal and Unmarshal compare field (Rejected) failed")
	}
	jsonValue = mm["rejected_pct"]
	if !compareBatchStatsValue(jsonValue, b.RejectedPct) {
		tt.Fatal("json Marshal and Unmarshal compare field (RejectedPct) failed")
	}
}

func TestBatchStatsMarshalJSONComplex(tt *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tt.Fatal(err)
	}
	ext := &Ext{Loc: loc, IsComplex: true, Verbose: "v"}
	b := NewBatchStats(nil)
	b.SetExt(ext)
	randSetBatchStats(tt, b)
	origin := dalVerboses[BatchStatsTid][b.ext.Verbose]
	dalVerboses[BatchStatsTid][b.ext.Verbose] = []map[db.Col]interface{}{
		{}, {BatchStatsCol.Updated: struct{}{}},
	}
	mm := marshalAndUnmarshalBatchStats(tt, b)
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
	mm = marshalAndUnmarshalBatchStats(tt, b)
	if !compareBatchStatsValue(mm["id"], rint) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite id")
	}
	if !compareBatchStatsValue(mm["updated"], rstr) {
		tt.Fatal("MarshalJSONComplex Padding must overwrite updated")
	}
	dalVerboses[BatchStatsTid][b.ext.Verbose] = origin
}

func TestBatchStatsUnmarshal(tt *testing.T) {
	mm := make(map[string]interface{})
	mm["id"] = randInt32()
	mm["created"] = randTime()
	mm["updated"] = randTime()
	mm["visibly"] = randBool()
	mm["batch_id"] = randInt32()
	mm["start"] = randTime()
	mm["end"] = randTime()
	mm["cnt"] = randInt32()
	mm["success"] = randInt32()
	mm["success_pct"] = randInt32()
	mm["right_first_time"] = randInt32()
	mm["right_first_time_pct"] = randInt32()
	mm["failed"] = randInt32()
	mm["failed_pct"] = randInt32()
	mm["rejected"] = randInt32()
	mm["rejected_pct"] = randInt32()
	bs, err := json.Marshal(mm)
	if err != nil {
		tt.Fatal(err)
	}
	mm = make(map[string]interface{})
	err = json.Unmarshal(bs, &mm)
	if err != nil {
		tt.Fatal(err)
	}
	b, err := UnmarshalBatchStats(nil, mm)
	if err != nil {
		tt.Fatal(err)
	}
	var jsonValue interface{}
	jsonValue = mm["id"]
	if !compareBatchStatsValue(jsonValue, b.Id) {
		tt.Fatal("json Marshal and Unmarshal compare field (Id) failed")
	}
	jsonValue = mm["created"]
	if !compareBatchStatsValue(jsonValue, b.Created) {
		tt.Fatal("json Marshal and Unmarshal compare field (Created) failed")
	}
	jsonValue = mm["updated"]
	if !compareBatchStatsValue(jsonValue, b.Updated) {
		tt.Fatal("json Marshal and Unmarshal compare field (Updated) failed")
	}
	jsonValue = mm["visibly"]
	if !compareBatchStatsValue(jsonValue, b.Visibly) {
		tt.Fatal("json Marshal and Unmarshal compare field (Visibly) failed")
	}
	jsonValue = mm["batch_id"]
	if !compareBatchStatsValue(jsonValue, b.BatchId) {
		tt.Fatal("json Marshal and Unmarshal compare field (BatchId) failed")
	}
	jsonValue = mm["start"]
	if !compareBatchStatsValue(jsonValue, b.Start) {
		tt.Fatal("json Marshal and Unmarshal compare field (Start) failed")
	}
	jsonValue = mm["end"]
	if !compareBatchStatsValue(jsonValue, b.End) {
		tt.Fatal("json Marshal and Unmarshal compare field (End) failed")
	}
	jsonValue = mm["cnt"]
	if !compareBatchStatsValue(jsonValue, b.Cnt) {
		tt.Fatal("json Marshal and Unmarshal compare field (Cnt) failed")
	}
	jsonValue = mm["success"]
	if !compareBatchStatsValue(jsonValue, b.Success) {
		tt.Fatal("json Marshal and Unmarshal compare field (Success) failed")
	}
	jsonValue = mm["success_pct"]
	if !compareBatchStatsValue(jsonValue, b.SuccessPct) {
		tt.Fatal("json Marshal and Unmarshal compare field (SuccessPct) failed")
	}
	jsonValue = mm["right_first_time"]
	if !compareBatchStatsValue(jsonValue, b.RightFirstTime) {
		tt.Fatal("json Marshal and Unmarshal compare field (RightFirstTime) failed")
	}
	jsonValue = mm["right_first_time_pct"]
	if !compareBatchStatsValue(jsonValue, b.RightFirstTimePct) {
		tt.Fatal("json Marshal and Unmarshal compare field (RightFirstTimePct) failed")
	}
	jsonValue = mm["failed"]
	if !compareBatchStatsValue(jsonValue, b.Failed) {
		tt.Fatal("json Marshal and Unmarshal compare field (Failed) failed")
	}
	jsonValue = mm["failed_pct"]
	if !compareBatchStatsValue(jsonValue, b.FailedPct) {
		tt.Fatal("json Marshal and Unmarshal compare field (FailedPct) failed")
	}
	jsonValue = mm["rejected"]
	if !compareBatchStatsValue(jsonValue, b.Rejected) {
		tt.Fatal("json Marshal and Unmarshal compare field (Rejected) failed")
	}
	jsonValue = mm["rejected_pct"]
	if !compareBatchStatsValue(jsonValue, b.RejectedPct) {
		tt.Fatal("json Marshal and Unmarshal compare field (RejectedPct) failed")
	}
}

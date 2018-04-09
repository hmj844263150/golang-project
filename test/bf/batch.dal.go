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

var BatchTid = 1
var _ db.Doer = (*Batch)(nil)
var batchcols = []db.Col{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
var batchfields = []string{"id", "created", "updated", "visibly", "sid", "factory_sid", "name", "desc", "cnt", "remain", "esp_mac_from", "esp_mac_to", "cus_mac_from", "cus_mac_to", "esp_mac_num_from", "esp_mac_num_to", "cus_mac_num_from", "cus_mac_num_to", "is_cus", "success", "right_first_time", "failed", "rejected", "statsed"}

var BatchCol = struct {
	Id, Created, Updated, Visibly, Sid, FactorySid, Name, Desc, Cnt, Remain, EspMacFrom, EspMacTo, CusMacFrom, CusMacTo, EspMacNumFrom, EspMacNumTo, CusMacNumFrom, CusMacNumTo, IsCus, Success, RightFirstTime, Failed, Rejected, Statsed, _ db.Col
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 0}

type Batch struct {
	Id             int
	Created        time.Time
	Updated        time.Time
	Visibly        bool
	Sid            string
	FactorySid     string
	Name           string
	Desc           string
	Cnt            int
	Remain         int
	EspMacFrom     string
	EspMacTo       string
	CusMacFrom     string
	CusMacTo       string
	EspMacNumFrom  int
	EspMacNumTo    int
	CusMacNumFrom  int
	CusMacNumTo    int
	IsCus          bool
	Success        int
	RightFirstTime int
	Failed         int
	Rejected       int
	Statsed        time.Time

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewBatch(ctx context.Context) *Batch {
	now := time.Now()
	b := &Batch{Created: now, Updated: now, Visibly: true}
	b.ext = GetExtFromContext(ctx)
	defaultBatch(ctx, b)
	return b
}

func FindBatch(ctx context.Context, id int) *Batch {
	dos, err := db.Open("Batch").Query(newBatchDest, true, batchSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		batch, _ := do.(*Batch)
		if ext != nil {
			batch.ext = ext
		}
		return batch
	}
	return nil
}

func ListBatch(ctx context.Context, ids ...int) []*Batch {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(batchSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("Batch").Query(newBatchDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	batchs := make([]*Batch, len(dos))
	for ii, do := range dos {
		batch, _ := do.(*Batch)
		if ext != nil {
			batch.ext = ext
		}
		batchs[ii] = batch
	}
	return batchs
}

func ListBatchAll(ctx context.Context, offset int, rowCount int) []*Batch {
	dos, err := db.Open("Batch").Query(newBatchDest, true, batchSqls[8], offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	batchs := make([]*Batch, len(dos))
	for ii, do := range dos {
		batch, _ := do.(*Batch)
		if ext != nil {
			batch.ext = ext
		}
		batchs[ii] = batch
	}
	return batchs
}

func FindBatchBySid(ctx context.Context, sid string) *Batch {
	dos, err := db.Open("Batch").Query(newBatchDest, true, batchSqls[9], sid)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		batch, _ := do.(*Batch)
		if ext != nil {
			batch.ext = ext
		}
		return batch
	}
	return nil
}

func FindBatchBySidEspMacNumRange(ctx context.Context, sid string, macNumFrom int, macNumTo int) *Batch {
	dos, err := db.Open("Batch").Query(newBatchDest, true, batchSqls[10], sid, macNumFrom, macNumTo)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		batch, _ := do.(*Batch)
		if ext != nil {
			batch.ext = ext
		}
		return batch
	}
	return nil
}

func ListBatchByFactorySid(ctx context.Context, factorySid string, offset int, rowCount int) []*Batch {
	dos, err := db.Open("Batch").Query(newBatchDest, true, batchSqls[11], factorySid, offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	batchs := make([]*Batch, len(dos))
	for ii, do := range dos {
		batch, _ := do.(*Batch)
		if ext != nil {
			batch.ext = ext
		}
		batchs[ii] = batch
	}
	return batchs
}

func (b *Batch) Save() error {
	now := time.Now()
	b.Created, b.Updated, b.Visibly = now, now, true
	var id int64
	var err error
	if b.Id == 0 {
		id, _, err = db.Open("Batch").Exec(batchSqls[0], b.Sid, b.FactorySid, b.Name, b.Desc, b.Cnt, b.Remain, b.EspMacFrom, b.EspMacTo, b.CusMacFrom, b.CusMacTo, b.EspMacNumFrom, b.EspMacNumTo, b.CusMacNumFrom, b.CusMacNumTo, b.IsCus, b.Success, b.RightFirstTime, b.Failed, b.Rejected, b.Statsed)
	} else {
		id, _, err = db.Open("Batch").Exec(batchSqls[1], b.Id, b.Sid, b.FactorySid, b.Name, b.Desc, b.Cnt, b.Remain, b.EspMacFrom, b.EspMacTo, b.CusMacFrom, b.CusMacTo, b.EspMacNumFrom, b.EspMacNumTo, b.CusMacNumFrom, b.CusMacNumTo, b.IsCus, b.Success, b.RightFirstTime, b.Failed, b.Rejected, b.Statsed)
	}
	if err != nil {
		return err
	}
	b.Id = int(id)
	return nil
}

func (b *Batch) Update(cs ...db.Col) error {
	if b.Id == 0 {
		return logError("dal.Batch Error: can not update row while id is zero")
	}
	b.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("Batch").Exec(batchSqls[2], b.Visibly, b.Sid, b.FactorySid, b.Name, b.Desc, b.Cnt, b.Remain, b.EspMacFrom, b.EspMacTo, b.CusMacFrom, b.CusMacTo, b.EspMacNumFrom, b.EspMacNumTo, b.CusMacNumFrom, b.CusMacNumTo, b.IsCus, b.Success, b.RightFirstTime, b.Failed, b.Rejected, b.Statsed, b.Id)
		return err
	}
	cols, args, err := colsAndArgsBatch(b, cs...)
	if err != nil {
		return err
	}
	args = append(args, b.Id)
	sqlstr := fmt.Sprintf(batchSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("Batch").Exec(sqlstr, args...)
	return err
}

func (b *Batch) Invisibly() error {
	if b.Id == 0 {
		return logError("dal.Batch Error: can not invisibly row while id is zero")
	}
	b.Updated = time.Now()
	b.Visibly = false
	_, _, err := db.Open("Batch").Exec(batchSqls[4], b.Id)
	return err
}

func (b *Batch) Delete() error {
	if b.Id == 0 {
		return logError("dal.Batch Error: can not delete row while id is zero")
	}
	b.Updated = time.Now()
	_, _, err := db.Open("Batch").Exec(batchSqls[5], b.Id)
	return err
}

func (b *Batch) Valid() error {
	return b.valid()
}

func (b *Batch) SetExt(ext *Ext) {
	b.ext = ext
}

func (b *Batch) Padding(pkey string, pvalue interface{}) {
	if b.ext == nil {
		b.ext = &Ext{Loc: DefaultLoc}
	}
	if b.paddings == nil {
		b.paddings = make(map[string]interface{})
	}
	b.paddings[pkey] = pvalue
	b.ext.IsComplex = true
}

func (b *Batch) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case BatchCol.Id:
			if isColumnName {
				mm["id"] = b.Id
			} else {
				mm["Id"] = b.Id
			}
		case BatchCol.Created:
			if isColumnName {
				mm["created"] = b.Created
			} else {
				mm["Created"] = b.Created
			}
		case BatchCol.Updated:
			if isColumnName {
				mm["updated"] = b.Updated
			} else {
				mm["Updated"] = b.Updated
			}
		case BatchCol.Visibly:
			if isColumnName {
				mm["visibly"] = b.Visibly
			} else {
				mm["Visibly"] = b.Visibly
			}
		case BatchCol.Sid:
			if isColumnName {
				mm["sid"] = b.Sid
			} else {
				mm["Sid"] = b.Sid
			}
		case BatchCol.FactorySid:
			if isColumnName {
				mm["factory_sid"] = b.FactorySid
			} else {
				mm["FactorySid"] = b.FactorySid
			}
		case BatchCol.Name:
			if isColumnName {
				mm["name"] = b.Name
			} else {
				mm["Name"] = b.Name
			}
		case BatchCol.Desc:
			if isColumnName {
				mm["desc"] = b.Desc
			} else {
				mm["Desc"] = b.Desc
			}
		case BatchCol.Cnt:
			if isColumnName {
				mm["cnt"] = b.Cnt
			} else {
				mm["Cnt"] = b.Cnt
			}
		case BatchCol.Remain:
			if isColumnName {
				mm["remain"] = b.Remain
			} else {
				mm["Remain"] = b.Remain
			}
		case BatchCol.EspMacFrom:
			if isColumnName {
				mm["esp_mac_from"] = b.EspMacFrom
			} else {
				mm["EspMacFrom"] = b.EspMacFrom
			}
		case BatchCol.EspMacTo:
			if isColumnName {
				mm["esp_mac_to"] = b.EspMacTo
			} else {
				mm["EspMacTo"] = b.EspMacTo
			}
		case BatchCol.CusMacFrom:
			if isColumnName {
				mm["cus_mac_from"] = b.CusMacFrom
			} else {
				mm["CusMacFrom"] = b.CusMacFrom
			}
		case BatchCol.CusMacTo:
			if isColumnName {
				mm["cus_mac_to"] = b.CusMacTo
			} else {
				mm["CusMacTo"] = b.CusMacTo
			}
		case BatchCol.EspMacNumFrom:
			if isColumnName {
				mm["esp_mac_num_from"] = b.EspMacNumFrom
			} else {
				mm["EspMacNumFrom"] = b.EspMacNumFrom
			}
		case BatchCol.EspMacNumTo:
			if isColumnName {
				mm["esp_mac_num_to"] = b.EspMacNumTo
			} else {
				mm["EspMacNumTo"] = b.EspMacNumTo
			}
		case BatchCol.CusMacNumFrom:
			if isColumnName {
				mm["cus_mac_num_from"] = b.CusMacNumFrom
			} else {
				mm["CusMacNumFrom"] = b.CusMacNumFrom
			}
		case BatchCol.CusMacNumTo:
			if isColumnName {
				mm["cus_mac_num_to"] = b.CusMacNumTo
			} else {
				mm["CusMacNumTo"] = b.CusMacNumTo
			}
		case BatchCol.IsCus:
			if isColumnName {
				mm["is_cus"] = b.IsCus
			} else {
				mm["IsCus"] = b.IsCus
			}
		case BatchCol.Success:
			if isColumnName {
				mm["success"] = b.Success
			} else {
				mm["Success"] = b.Success
			}
		case BatchCol.RightFirstTime:
			if isColumnName {
				mm["right_first_time"] = b.RightFirstTime
			} else {
				mm["RightFirstTime"] = b.RightFirstTime
			}
		case BatchCol.Failed:
			if isColumnName {
				mm["failed"] = b.Failed
			} else {
				mm["Failed"] = b.Failed
			}
		case BatchCol.Rejected:
			if isColumnName {
				mm["rejected"] = b.Rejected
			} else {
				mm["Rejected"] = b.Rejected
			}
		case BatchCol.Statsed:
			if isColumnName {
				mm["statsed"] = b.Statsed
			} else {
				mm["Statsed"] = b.Statsed
			}
		default:
			logError(fmt.Sprintf("dal.Batch Error: unknow column num %d in talbe batch", cc))
		}
	}
	return mm
}

func (b *Batch) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	loc := DefaultLoc
	var numericEnum bool
	if b.ext != nil {
		if b.ext.IsComplex {
			return b.marshalJSONComplex()
		}
		loc = b.ext.Loc
		numericEnum = b.ext.NumericEnum
	}
	var buf bytes.Buffer
	buf.WriteString(`{"id":`)
	buf.WriteString(strconv.FormatInt(int64(b.Id), 10))
	buf.WriteString(`, "created":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(b.Created.Unix(), 10))
	} else {
		b.Created = b.Created.In(loc)
		buf.WriteString(`"` + b.Created.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "updated":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(b.Updated.Unix(), 10))
	} else {
		b.Updated = b.Updated.In(loc)
		buf.WriteString(`"` + b.Updated.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "visibly":`)
	if numericEnum {
		if b.Visibly {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(b.Visibly))
	}
	buf.WriteString(`, "sid":`)
	WriteJsonString(&buf, b.Sid)
	buf.WriteString(`, "factory_sid":`)
	WriteJsonString(&buf, b.FactorySid)
	buf.WriteString(`, "name":`)
	WriteJsonString(&buf, b.Name)
	buf.WriteString(`, "desc":`)
	WriteJsonString(&buf, b.Desc)
	buf.WriteString(`, "cnt":`)
	buf.WriteString(strconv.FormatInt(int64(b.Cnt), 10))
	buf.WriteString(`, "remain":`)
	buf.WriteString(strconv.FormatInt(int64(b.Remain), 10))
	buf.WriteString(`, "esp_mac_from":`)
	WriteJsonString(&buf, b.EspMacFrom)
	buf.WriteString(`, "esp_mac_to":`)
	WriteJsonString(&buf, b.EspMacTo)
	buf.WriteString(`, "cus_mac_from":`)
	WriteJsonString(&buf, b.CusMacFrom)
	buf.WriteString(`, "cus_mac_to":`)
	WriteJsonString(&buf, b.CusMacTo)
	buf.WriteString(`, "esp_mac_num_from":`)
	buf.WriteString(strconv.FormatInt(int64(b.EspMacNumFrom), 10))
	buf.WriteString(`, "esp_mac_num_to":`)
	buf.WriteString(strconv.FormatInt(int64(b.EspMacNumTo), 10))
	buf.WriteString(`, "cus_mac_num_from":`)
	buf.WriteString(strconv.FormatInt(int64(b.CusMacNumFrom), 10))
	buf.WriteString(`, "cus_mac_num_to":`)
	buf.WriteString(strconv.FormatInt(int64(b.CusMacNumTo), 10))
	buf.WriteString(`, "is_cus":`)
	if numericEnum {
		if b.IsCus {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(b.IsCus))
	}
	buf.WriteString(`, "success":`)
	buf.WriteString(strconv.FormatInt(int64(b.Success), 10))
	buf.WriteString(`, "right_first_time":`)
	buf.WriteString(strconv.FormatInt(int64(b.RightFirstTime), 10))
	buf.WriteString(`, "failed":`)
	buf.WriteString(strconv.FormatInt(int64(b.Failed), 10))
	buf.WriteString(`, "rejected":`)
	buf.WriteString(strconv.FormatInt(int64(b.Rejected), 10))
	buf.WriteString(`, "statsed":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(b.Statsed.Unix(), 10))
	} else {
		b.Statsed = b.Statsed.In(loc)
		buf.WriteString(`"` + b.Statsed.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (b *Batch) marshalJSONComplex() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	if b.ext == nil {
		return nil, logError("dal.Batch Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := b.ext.Loc
	numericEnum := b.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[BatchTid]; ok {
		if vvv, ok := vv[b.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := b.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(BatchCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(b.Id), 10))
	}
	isRender = isRenderField(BatchCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.Created.Unix(), 10))
		} else {
			b.Created = b.Created.In(loc)
			buf.WriteString(`"` + b.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(BatchCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.Updated.Unix(), 10))
		} else {
			b.Updated = b.Updated.In(loc)
			buf.WriteString(`"` + b.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(BatchCol.Visibly, "visibly", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "visibly":`)
		if numericEnum {
			if b.Visibly {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(b.Visibly))
		}
	}
	isRender = isRenderField(BatchCol.Sid, "sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "sid":`)
		WriteJsonString(&buf, b.Sid)
	}
	isRender = isRenderField(BatchCol.FactorySid, "factory_sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "factory_sid":`)
		WriteJsonString(&buf, b.FactorySid)
	}
	isRender = isRenderField(BatchCol.Name, "name", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "name":`)
		WriteJsonString(&buf, b.Name)
	}
	isRender = isRenderField(BatchCol.Desc, "desc", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "desc":`)
		WriteJsonString(&buf, b.Desc)
	}
	isRender = isRenderField(BatchCol.Cnt, "cnt", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cnt":`)
		buf.WriteString(strconv.FormatInt(int64(b.Cnt), 10))
	}
	isRender = isRenderField(BatchCol.Remain, "remain", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "remain":`)
		buf.WriteString(strconv.FormatInt(int64(b.Remain), 10))
	}
	isRender = isRenderField(BatchCol.EspMacFrom, "esp_mac_from", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac_from":`)
		WriteJsonString(&buf, b.EspMacFrom)
	}
	isRender = isRenderField(BatchCol.EspMacTo, "esp_mac_to", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac_to":`)
		WriteJsonString(&buf, b.EspMacTo)
	}
	isRender = isRenderField(BatchCol.CusMacFrom, "cus_mac_from", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cus_mac_from":`)
		WriteJsonString(&buf, b.CusMacFrom)
	}
	isRender = isRenderField(BatchCol.CusMacTo, "cus_mac_to", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cus_mac_to":`)
		WriteJsonString(&buf, b.CusMacTo)
	}
	isRender = isRenderField(BatchCol.EspMacNumFrom, "esp_mac_num_from", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac_num_from":`)
		buf.WriteString(strconv.FormatInt(int64(b.EspMacNumFrom), 10))
	}
	isRender = isRenderField(BatchCol.EspMacNumTo, "esp_mac_num_to", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac_num_to":`)
		buf.WriteString(strconv.FormatInt(int64(b.EspMacNumTo), 10))
	}
	isRender = isRenderField(BatchCol.CusMacNumFrom, "cus_mac_num_from", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cus_mac_num_from":`)
		buf.WriteString(strconv.FormatInt(int64(b.CusMacNumFrom), 10))
	}
	isRender = isRenderField(BatchCol.CusMacNumTo, "cus_mac_num_to", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cus_mac_num_to":`)
		buf.WriteString(strconv.FormatInt(int64(b.CusMacNumTo), 10))
	}
	isRender = isRenderField(BatchCol.IsCus, "is_cus", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "is_cus":`)
		if numericEnum {
			if b.IsCus {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(b.IsCus))
		}
	}
	isRender = isRenderField(BatchCol.Success, "success", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "success":`)
		buf.WriteString(strconv.FormatInt(int64(b.Success), 10))
	}
	isRender = isRenderField(BatchCol.RightFirstTime, "right_first_time", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "right_first_time":`)
		buf.WriteString(strconv.FormatInt(int64(b.RightFirstTime), 10))
	}
	isRender = isRenderField(BatchCol.Failed, "failed", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failed":`)
		buf.WriteString(strconv.FormatInt(int64(b.Failed), 10))
	}
	isRender = isRenderField(BatchCol.Rejected, "rejected", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "rejected":`)
		buf.WriteString(strconv.FormatInt(int64(b.Rejected), 10))
	}
	isRender = isRenderField(BatchCol.Statsed, "statsed", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "statsed":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.Statsed.Unix(), 10))
		} else {
			b.Statsed = b.Statsed.In(loc)
			buf.WriteString(`"` + b.Statsed.Format(DefaultTimeFormat) + `"`)
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

func (b *Batch) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalBatch with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalBatch")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = batchcols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case BatchCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			b.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			b.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case BatchCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			b.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case BatchCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			b.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Sid:
			vvv, ok := vv["sid"]
			if !ok {
				continue
			}
			b.Sid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.FactorySid:
			vvv, ok := vv["factory_sid"]
			if !ok {
				continue
			}
			b.FactorySid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Name:
			vvv, ok := vv["name"]
			if !ok {
				continue
			}
			b.Name, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Desc:
			vvv, ok := vv["desc"]
			if !ok {
				continue
			}
			b.Desc, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Cnt:
			vvv, ok := vv["cnt"]
			if !ok {
				continue
			}
			b.Cnt, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Remain:
			vvv, ok := vv["remain"]
			if !ok {
				continue
			}
			b.Remain, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.EspMacFrom:
			vvv, ok := vv["esp_mac_from"]
			if !ok {
				continue
			}
			b.EspMacFrom, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.EspMacTo:
			vvv, ok := vv["esp_mac_to"]
			if !ok {
				continue
			}
			b.EspMacTo, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.CusMacFrom:
			vvv, ok := vv["cus_mac_from"]
			if !ok {
				continue
			}
			b.CusMacFrom, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.CusMacTo:
			vvv, ok := vv["cus_mac_to"]
			if !ok {
				continue
			}
			b.CusMacTo, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.EspMacNumFrom:
			vvv, ok := vv["esp_mac_num_from"]
			if !ok {
				continue
			}
			b.EspMacNumFrom, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.EspMacNumTo:
			vvv, ok := vv["esp_mac_num_to"]
			if !ok {
				continue
			}
			b.EspMacNumTo, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.CusMacNumFrom:
			vvv, ok := vv["cus_mac_num_from"]
			if !ok {
				continue
			}
			b.CusMacNumFrom, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.CusMacNumTo:
			vvv, ok := vv["cus_mac_num_to"]
			if !ok {
				continue
			}
			b.CusMacNumTo, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.IsCus:
			vvv, ok := vv["is_cus"]
			if !ok {
				continue
			}
			b.IsCus, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Success:
			vvv, ok := vv["success"]
			if !ok {
				continue
			}
			b.Success, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.RightFirstTime:
			vvv, ok := vv["right_first_time"]
			if !ok {
				continue
			}
			b.RightFirstTime, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Failed:
			vvv, ok := vv["failed"]
			if !ok {
				continue
			}
			b.Failed, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Rejected:
			vvv, ok := vv["rejected"]
			if !ok {
				continue
			}
			b.Rejected, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchCol.Statsed:
			vvv, ok := vv["statsed"]
			if !ok {
				continue
			}
			b.Statsed, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalBatch(ctx context.Context, vi interface{}, cols ...db.Col) (*Batch, error) {
	b := NewBatch(ctx)
	_, err := b.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return b, err
}

func UnmarshalBatchs(ctx context.Context, vi interface{}, cols ...db.Col) ([]*Batch, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalBatchs with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalBatchs")
	}
	batchs := make([]*Batch, len(vv))
	for ii, vvv := range vv {
		var b *Batch
		b, err = UnmarshalBatch(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		batchs[ii] = b
	}
	return batchs, nil
}

func newBatchDest(cols ...string) (db.Doer, []interface{}, error) {
	b := &Batch{}
	if cols == nil || len(cols) == 0 {
		return b, []interface{}{&b.Id, &b.Created, &b.Updated, &b.Visibly, &b.Sid, &b.FactorySid, &b.Name, &b.Desc, &b.Cnt, &b.Remain, &b.EspMacFrom, &b.EspMacTo, &b.CusMacFrom, &b.CusMacTo, &b.EspMacNumFrom, &b.EspMacNumTo, &b.CusMacNumFrom, &b.CusMacNumTo, &b.IsCus, &b.Success, &b.RightFirstTime, &b.Failed, &b.Rejected, &b.Statsed}, nil
	}
	dest := make([]interface{}, len(cols))
	for ii, col := range cols {
		switch col {
		case "id":
			dest[ii] = &b.Id
		case "created":
			dest[ii] = &b.Created
		case "updated":
			dest[ii] = &b.Updated
		case "visibly":
			dest[ii] = &b.Visibly
		case "sid":
			dest[ii] = &b.Sid
		case "factory_sid":
			dest[ii] = &b.FactorySid
		case "name":
			dest[ii] = &b.Name
		case "desc":
			dest[ii] = &b.Desc
		case "cnt":
			dest[ii] = &b.Cnt
		case "remain":
			dest[ii] = &b.Remain
		case "esp_mac_from":
			dest[ii] = &b.EspMacFrom
		case "esp_mac_to":
			dest[ii] = &b.EspMacTo
		case "cus_mac_from":
			dest[ii] = &b.CusMacFrom
		case "cus_mac_to":
			dest[ii] = &b.CusMacTo
		case "esp_mac_num_from":
			dest[ii] = &b.EspMacNumFrom
		case "esp_mac_num_to":
			dest[ii] = &b.EspMacNumTo
		case "cus_mac_num_from":
			dest[ii] = &b.CusMacNumFrom
		case "cus_mac_num_to":
			dest[ii] = &b.CusMacNumTo
		case "is_cus":
			dest[ii] = &b.IsCus
		case "success":
			dest[ii] = &b.Success
		case "right_first_time":
			dest[ii] = &b.RightFirstTime
		case "failed":
			dest[ii] = &b.Failed
		case "rejected":
			dest[ii] = &b.Rejected
		case "statsed":
			dest[ii] = &b.Statsed
		default:
			return nil, nil, logError("dal.Batch Error: unknow column " + col + " in talbe batch")
		}
	}
	return b, dest, nil
}

func colsAndArgsBatch(b *Batch, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.Batch Error: at least one column to colsAndArgsBatch")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case BatchCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = b.Id
		case BatchCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = b.Created
		case BatchCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = b.Updated
		case BatchCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = b.Visibly
		case BatchCol.Sid:
			cols[ii] = "`sid` = ?"
			args[ii] = b.Sid
		case BatchCol.FactorySid:
			cols[ii] = "`factory_sid` = ?"
			args[ii] = b.FactorySid
		case BatchCol.Name:
			cols[ii] = "`name` = ?"
			args[ii] = b.Name
		case BatchCol.Desc:
			cols[ii] = "`desc` = ?"
			args[ii] = b.Desc
		case BatchCol.Cnt:
			cols[ii] = "`cnt` = ?"
			args[ii] = b.Cnt
		case BatchCol.Remain:
			cols[ii] = "`remain` = ?"
			args[ii] = b.Remain
		case BatchCol.EspMacFrom:
			cols[ii] = "`esp_mac_from` = ?"
			args[ii] = b.EspMacFrom
		case BatchCol.EspMacTo:
			cols[ii] = "`esp_mac_to` = ?"
			args[ii] = b.EspMacTo
		case BatchCol.CusMacFrom:
			cols[ii] = "`cus_mac_from` = ?"
			args[ii] = b.CusMacFrom
		case BatchCol.CusMacTo:
			cols[ii] = "`cus_mac_to` = ?"
			args[ii] = b.CusMacTo
		case BatchCol.EspMacNumFrom:
			cols[ii] = "`esp_mac_num_from` = ?"
			args[ii] = b.EspMacNumFrom
		case BatchCol.EspMacNumTo:
			cols[ii] = "`esp_mac_num_to` = ?"
			args[ii] = b.EspMacNumTo
		case BatchCol.CusMacNumFrom:
			cols[ii] = "`cus_mac_num_from` = ?"
			args[ii] = b.CusMacNumFrom
		case BatchCol.CusMacNumTo:
			cols[ii] = "`cus_mac_num_to` = ?"
			args[ii] = b.CusMacNumTo
		case BatchCol.IsCus:
			cols[ii] = "`is_cus` = ?"
			args[ii] = b.IsCus
		case BatchCol.Success:
			cols[ii] = "`success` = ?"
			args[ii] = b.Success
		case BatchCol.RightFirstTime:
			cols[ii] = "`right_first_time` = ?"
			args[ii] = b.RightFirstTime
		case BatchCol.Failed:
			cols[ii] = "`failed` = ?"
			args[ii] = b.Failed
		case BatchCol.Rejected:
			cols[ii] = "`rejected` = ?"
			args[ii] = b.Rejected
		case BatchCol.Statsed:
			cols[ii] = "`statsed` = ?"
			args[ii] = b.Statsed
		default:
			return nil, nil, logError(fmt.Sprintf("dal.Batch Error: unknow column num %d in talbe batch", cc))
		}
	}
	return cols, args, nil
}

var BatchEnum = struct {
}{}

var batchSqls = []string{
	/*
		CREATE TABLE `batch` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `sid` varchar(64) NOT NULL,
		  `factory_sid` varchar(64) NOT NULL,
		  `name` varchar(64) NOT NULL,
		  `desc` varchar(128) NOT NULL,
		  `cnt` int(11) NOT NULL,
		  `remain` int(11) NOT NULL,
		  `esp_mac_from` varchar(64) NOT NULL,
		  `esp_mac_to` varchar(64) NOT NULL,
		  `cus_mac_from` varchar(64) NOT NULL,
		  `cus_mac_to` varchar(64) NOT NULL,
		  `esp_mac_num_from` int(11) NOT NULL,
		  `esp_mac_num_to` int(11) NOT NULL,
		  `cus_mac_num_from` int(11) NOT NULL,
		  `cus_mac_num_to` int(11) NOT NULL,
		  `is_cus` bool NOT NULL,
		  `success` int(11) NOT NULL,
		  `right_first_time` int(11) NOT NULL,
		  `failed` int(11) NOT NULL,
		  `rejected` int(11) NOT NULL,
		  `statsed` datetime NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into batch(`created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed`) values(now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*1*/ "insert into batch(`id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed`) values(?, now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*2*/ "update batch set updated = now(), `visibly` = ?, `sid` = ?, `factory_sid` = ?, `name` = ?, `desc` = ?, `cnt` = ?, `remain` = ?, `esp_mac_from` = ?, `esp_mac_to` = ?, `cus_mac_from` = ?, `cus_mac_to` = ?, `esp_mac_num_from` = ?, `esp_mac_num_to` = ?, `cus_mac_num_from` = ?, `cus_mac_num_to` = ?, `is_cus` = ?, `success` = ?, `right_first_time` = ?, `failed` = ?, `rejected` = ?, `statsed` = ? where id = ?",
	/*3*/ "update batch set updated = now(), %s where id = ?",
	/*4*/ "update batch set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from batch where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed` from batch where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed` from batch where id in (%s) and visibly = 1",

	/*8*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed` from batch where visibly = 1 order by id desc limit ?, ?",
	/*9*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed` from batch where visibly = 1 and sid = ? limit 0, 1",
	/*10*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed` from batch where visibly = 1 and sid = ? and cus_mac_num_from >= ? and ? <= cus_mac_num_to limit 0, 1",
	/*11*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `factory_sid`, `name`, `desc`, `cnt`, `remain`, `esp_mac_from`, `esp_mac_to`, `cus_mac_from`, `cus_mac_to`, `esp_mac_num_from`, `esp_mac_num_to`, `cus_mac_num_from`, `cus_mac_num_to`, `is_cus`, `success`, `right_first_time`, `failed`, `rejected`, `statsed` from batch where visibly = 1 and factory_sid = ? order by id desc limit ?, ?",
}

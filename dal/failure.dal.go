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

var FailureTid = 10
var _ db.Doer = (*Failure)(nil)
var failurecols = []db.Col{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
var failurefields = []string{"id", "created", "updated", "visibly", "factory_id", "batch_id", "mac_int", "mode", "mac", "latest", "raw", "is_failed", "failure_mode", "failure_mode1", "failure_mode2", "failure_mode3"}

var FailureCol = struct {
	Id, Created, Updated, Visibly, FactoryId, BatchId, MacInt, Mode, Mac, Latest, raw, IsFailed, FailureMode, FailureMode1, FailureMode2, FailureMode3, _ db.Col
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 0}

type Failure struct {
	Id           int
	Created      time.Time
	Updated      time.Time
	Visibly      bool
	FactoryId    int
	BatchId      int
	MacInt       int
	Mode         int
	Mac          string
	Latest       bool
	raw          string
	IsFailed     bool
	FailureMode  int
	FailureMode1 int
	FailureMode2 int
	FailureMode3 int

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewFailure(ctx context.Context) *Failure {
	now := time.Now()
	f := &Failure{Created: now, Updated: now, Visibly: true}
	f.ext = GetExtFromContext(ctx)
	defaultFailure(ctx, f)
	return f
}

func FindFailure(ctx context.Context, id int) *Failure {
	dos, err := db.Open("Failure").Query(newFailureDest, true, failureSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		failure, _ := do.(*Failure)
		if ext != nil {
			failure.ext = ext
		}
		return failure
	}
	return nil
}

func ListFailure(ctx context.Context, ids ...int) []*Failure {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(failureSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("Failure").Query(newFailureDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	failures := make([]*Failure, len(dos))
	for ii, do := range dos {
		failure, _ := do.(*Failure)
		if ext != nil {
			failure.ext = ext
		}
		failures[ii] = failure
	}
	return failures
}

func (f *Failure) Save() error {
	now := time.Now()
	f.Created, f.Updated, f.Visibly = now, now, true
	var id int64
	var err error
	if f.Id == 0 {
		id, _, err = db.Open("Failure").Exec(failureSqls[0], f.FactoryId, f.BatchId, f.MacInt, f.Mode, f.Mac, f.Latest, f.raw, f.IsFailed, f.FailureMode, f.FailureMode1, f.FailureMode2, f.FailureMode3)
	} else {
		id, _, err = db.Open("Failure").Exec(failureSqls[1], f.Id, f.FactoryId, f.BatchId, f.MacInt, f.Mode, f.Mac, f.Latest, f.raw, f.IsFailed, f.FailureMode, f.FailureMode1, f.FailureMode2, f.FailureMode3)
	}
	if err != nil {
		return err
	}
	f.Id = int(id)
	return nil
}

func (f *Failure) Update(cs ...db.Col) error {
	if f.Id == 0 {
		return logError("dal.Failure Error: can not update row while id is zero")
	}
	f.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("Failure").Exec(failureSqls[2], f.Visibly, f.FactoryId, f.BatchId, f.MacInt, f.Mode, f.Mac, f.Latest, f.raw, f.IsFailed, f.FailureMode, f.FailureMode1, f.FailureMode2, f.FailureMode3, f.Id)
		return err
	}
	cols, args, err := colsAndArgsFailure(f, cs...)
	if err != nil {
		return err
	}
	args = append(args, f.Id)
	sqlstr := fmt.Sprintf(failureSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("Failure").Exec(sqlstr, args...)
	return err
}

func (f *Failure) Invisibly() error {
	if f.Id == 0 {
		return logError("dal.Failure Error: can not invisibly row while id is zero")
	}
	f.Updated = time.Now()
	f.Visibly = false
	_, _, err := db.Open("Failure").Exec(failureSqls[4], f.Id)
	return err
}

func (f *Failure) Delete() error {
	if f.Id == 0 {
		return logError("dal.Failure Error: can not delete row while id is zero")
	}
	f.Updated = time.Now()
	_, _, err := db.Open("Failure").Exec(failureSqls[5], f.Id)
	return err
}

func (f *Failure) Valid() error {
	if _, ok := FailureEnum.Mode.RMAP[f.Mode]; !ok {
		return logError("Mode enum must one of " + FailureEnum.Mode.LIST)
	}
	if _, ok := FailureEnum.FailureMode.RMAP[f.FailureMode]; !ok {
		return logError("FailureMode enum must one of " + FailureEnum.FailureMode.LIST)
	}
	if _, ok := FailureEnum.FailureMode1.RMAP[f.FailureMode1]; !ok {
		return logError("FailureMode1 enum must one of " + FailureEnum.FailureMode1.LIST)
	}
	if _, ok := FailureEnum.FailureMode2.RMAP[f.FailureMode2]; !ok {
		return logError("FailureMode2 enum must one of " + FailureEnum.FailureMode2.LIST)
	}
	if _, ok := FailureEnum.FailureMode3.RMAP[f.FailureMode3]; !ok {
		return logError("FailureMode3 enum must one of " + FailureEnum.FailureMode3.LIST)
	}
	return f.valid()
}

func (f *Failure) SetExt(ext *Ext) {
	f.ext = ext
}

func (f *Failure) Padding(pkey string, pvalue interface{}) {
	if f.ext == nil {
		f.ext = &Ext{Loc: DefaultLoc}
	}
	if f.paddings == nil {
		f.paddings = make(map[string]interface{})
	}
	f.paddings[pkey] = pvalue
	f.ext.IsComplex = true
}

func (f *Failure) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case FailureCol.Id:
			if isColumnName {
				mm["id"] = f.Id
			} else {
				mm["Id"] = f.Id
			}
		case FailureCol.Created:
			if isColumnName {
				mm["created"] = f.Created
			} else {
				mm["Created"] = f.Created
			}
		case FailureCol.Updated:
			if isColumnName {
				mm["updated"] = f.Updated
			} else {
				mm["Updated"] = f.Updated
			}
		case FailureCol.Visibly:
			if isColumnName {
				mm["visibly"] = f.Visibly
			} else {
				mm["Visibly"] = f.Visibly
			}
		case FailureCol.FactoryId:
			if isColumnName {
				mm["factory_id"] = f.FactoryId
			} else {
				mm["FactoryId"] = f.FactoryId
			}
		case FailureCol.BatchId:
			if isColumnName {
				mm["batch_id"] = f.BatchId
			} else {
				mm["BatchId"] = f.BatchId
			}
		case FailureCol.MacInt:
			if isColumnName {
				mm["mac_int"] = f.MacInt
			} else {
				mm["MacInt"] = f.MacInt
			}
		case FailureCol.Mode:
			if isColumnName {
				mm["mode"] = f.Mode
			} else {
				mm["Mode"] = f.Mode
			}
		case FailureCol.Mac:
			if isColumnName {
				mm["mac"] = f.Mac
			} else {
				mm["Mac"] = f.Mac
			}
		case FailureCol.Latest:
			if isColumnName {
				mm["latest"] = f.Latest
			} else {
				mm["Latest"] = f.Latest
			}
		case FailureCol.raw:
			if isColumnName {
				mm["raw"] = f.raw
			} else {
				mm["raw"] = f.raw
			}
		case FailureCol.IsFailed:
			if isColumnName {
				mm["is_failed"] = f.IsFailed
			} else {
				mm["IsFailed"] = f.IsFailed
			}
		case FailureCol.FailureMode:
			if isColumnName {
				mm["failure_mode"] = f.FailureMode
			} else {
				mm["FailureMode"] = f.FailureMode
			}
		case FailureCol.FailureMode1:
			if isColumnName {
				mm["failure_mode1"] = f.FailureMode1
			} else {
				mm["FailureMode1"] = f.FailureMode1
			}
		case FailureCol.FailureMode2:
			if isColumnName {
				mm["failure_mode2"] = f.FailureMode2
			} else {
				mm["FailureMode2"] = f.FailureMode2
			}
		case FailureCol.FailureMode3:
			if isColumnName {
				mm["failure_mode3"] = f.FailureMode3
			} else {
				mm["FailureMode3"] = f.FailureMode3
			}
		default:
			logError(fmt.Sprintf("dal.Failure Error: unknow column num %d in talbe failure", cc))
		}
	}
	return mm
}

func (f *Failure) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}
	loc := DefaultLoc
	var numericEnum bool
	var ee string
	if f.ext != nil {
		if f.ext.IsComplex {
			return f.marshalJSONComplex()
		}
		loc = f.ext.Loc
		numericEnum = f.ext.NumericEnum
	}
	var buf bytes.Buffer
	buf.WriteString(`{"id":`)
	buf.WriteString(strconv.FormatInt(int64(f.Id), 10))
	buf.WriteString(`, "created":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(f.Created.Unix(), 10))
	} else {
		f.Created = f.Created.In(loc)
		buf.WriteString(`"` + f.Created.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "updated":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(f.Updated.Unix(), 10))
	} else {
		f.Updated = f.Updated.In(loc)
		buf.WriteString(`"` + f.Updated.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "visibly":`)
	if numericEnum {
		if f.Visibly {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(f.Visibly))
	}
	buf.WriteString(`, "factory_id":`)
	buf.WriteString(strconv.FormatInt(int64(f.FactoryId), 10))
	buf.WriteString(`, "batch_id":`)
	buf.WriteString(strconv.FormatInt(int64(f.BatchId), 10))
	buf.WriteString(`, "mac_int":`)
	buf.WriteString(strconv.FormatInt(int64(f.MacInt), 10))
	buf.WriteString(`, "mode":`)
	if numericEnum {
		buf.WriteString(strconv.FormatInt(int64(f.Mode), 10))
	} else {
		ee = "_"
		if vv, ok := FailureEnum.Mode.RMAP[f.Mode]; ok {
			ee = vv
		}
		buf.WriteString(`"` + ee + `"`)
	}
	buf.WriteString(`, "mac":`)
	WriteJsonString(&buf, f.Mac)
	buf.WriteString(`, "latest":`)
	if numericEnum {
		if f.Latest {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(f.Latest))
	}
	buf.WriteString(`, "raw":`)
	WriteJsonString(&buf, f.raw)
	buf.WriteString(`, "is_failed":`)
	if numericEnum {
		if f.IsFailed {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(f.IsFailed))
	}
	buf.WriteString(`, "failure_mode":`)
	if numericEnum {
		buf.WriteString(strconv.FormatInt(int64(f.FailureMode), 10))
	} else {
		ee = "_"
		if vv, ok := FailureEnum.FailureMode.RMAP[f.FailureMode]; ok {
			ee = vv
		}
		buf.WriteString(`"` + ee + `"`)
	}
	buf.WriteString(`, "failure_mode1":`)
	if numericEnum {
		buf.WriteString(strconv.FormatInt(int64(f.FailureMode1), 10))
	} else {
		ee = "_"
		if vv, ok := FailureEnum.FailureMode1.RMAP[f.FailureMode1]; ok {
			ee = vv
		}
		buf.WriteString(`"` + ee + `"`)
	}
	buf.WriteString(`, "failure_mode2":`)
	if numericEnum {
		buf.WriteString(strconv.FormatInt(int64(f.FailureMode2), 10))
	} else {
		ee = "_"
		if vv, ok := FailureEnum.FailureMode2.RMAP[f.FailureMode2]; ok {
			ee = vv
		}
		buf.WriteString(`"` + ee + `"`)
	}
	buf.WriteString(`, "failure_mode3":`)
	if numericEnum {
		buf.WriteString(strconv.FormatInt(int64(f.FailureMode3), 10))
	} else {
		ee = "_"
		if vv, ok := FailureEnum.FailureMode3.RMAP[f.FailureMode3]; ok {
			ee = vv
		}
		buf.WriteString(`"` + ee + `"`)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (f *Failure) marshalJSONComplex() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}
	if f.ext == nil {
		return nil, logError("dal.Failure Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := f.ext.Loc
	numericEnum := f.ext.NumericEnum
	var ee string
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[FailureTid]; ok {
		if vvv, ok := vv[f.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := f.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(FailureCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(f.Id), 10))
	}
	isRender = isRenderField(FailureCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(f.Created.Unix(), 10))
		} else {
			f.Created = f.Created.In(loc)
			buf.WriteString(`"` + f.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(FailureCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(f.Updated.Unix(), 10))
		} else {
			f.Updated = f.Updated.In(loc)
			buf.WriteString(`"` + f.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(FailureCol.Visibly, "visibly", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "visibly":`)
		if numericEnum {
			if f.Visibly {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(f.Visibly))
		}
	}
	isRender = isRenderField(FailureCol.FactoryId, "factory_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "factory_id":`)
		buf.WriteString(strconv.FormatInt(int64(f.FactoryId), 10))
	}
	isRender = isRenderField(FailureCol.BatchId, "batch_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "batch_id":`)
		buf.WriteString(strconv.FormatInt(int64(f.BatchId), 10))
	}
	isRender = isRenderField(FailureCol.MacInt, "mac_int", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "mac_int":`)
		buf.WriteString(strconv.FormatInt(int64(f.MacInt), 10))
	}
	isRender = isRenderField(FailureCol.Mode, "mode", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "mode":`)
		if numericEnum {
			buf.WriteString(strconv.FormatInt(int64(f.Mode), 10))
		} else {
			ee = "_"
			if vv, ok := FailureEnum.Mode.RMAP[f.Mode]; ok {
				ee = vv
			}
			buf.WriteString(`"` + ee + `"`)
		}
	}
	isRender = isRenderField(FailureCol.Mac, "mac", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "mac":`)
		WriteJsonString(&buf, f.Mac)
	}
	isRender = isRenderField(FailureCol.Latest, "latest", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "latest":`)
		if numericEnum {
			if f.Latest {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(f.Latest))
		}
	}
	isRender = isRenderField(FailureCol.raw, "raw", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "raw":`)
		WriteJsonString(&buf, f.raw)
	}
	isRender = isRenderField(FailureCol.IsFailed, "is_failed", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "is_failed":`)
		if numericEnum {
			if f.IsFailed {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(f.IsFailed))
		}
	}
	isRender = isRenderField(FailureCol.FailureMode, "failure_mode", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failure_mode":`)
		if numericEnum {
			buf.WriteString(strconv.FormatInt(int64(f.FailureMode), 10))
		} else {
			ee = "_"
			if vv, ok := FailureEnum.FailureMode.RMAP[f.FailureMode]; ok {
				ee = vv
			}
			buf.WriteString(`"` + ee + `"`)
		}
	}
	isRender = isRenderField(FailureCol.FailureMode1, "failure_mode1", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failure_mode1":`)
		if numericEnum {
			buf.WriteString(strconv.FormatInt(int64(f.FailureMode1), 10))
		} else {
			ee = "_"
			if vv, ok := FailureEnum.FailureMode1.RMAP[f.FailureMode1]; ok {
				ee = vv
			}
			buf.WriteString(`"` + ee + `"`)
		}
	}
	isRender = isRenderField(FailureCol.FailureMode2, "failure_mode2", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failure_mode2":`)
		if numericEnum {
			buf.WriteString(strconv.FormatInt(int64(f.FailureMode2), 10))
		} else {
			ee = "_"
			if vv, ok := FailureEnum.FailureMode2.RMAP[f.FailureMode2]; ok {
				ee = vv
			}
			buf.WriteString(`"` + ee + `"`)
		}
	}
	isRender = isRenderField(FailureCol.FailureMode3, "failure_mode3", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failure_mode3":`)
		if numericEnum {
			buf.WriteString(strconv.FormatInt(int64(f.FailureMode3), 10))
		} else {
			ee = "_"
			if vv, ok := FailureEnum.FailureMode3.RMAP[f.FailureMode3]; ok {
				ee = vv
			}
			buf.WriteString(`"` + ee + `"`)
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

func (f *Failure) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalFailure with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalFailure")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = failurecols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case FailureCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			f.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			f.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case FailureCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			f.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case FailureCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			f.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.FactoryId:
			vvv, ok := vv["factory_id"]
			if !ok {
				continue
			}
			f.FactoryId, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.BatchId:
			vvv, ok := vv["batch_id"]
			if !ok {
				continue
			}
			f.BatchId, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.MacInt:
			vvv, ok := vv["mac_int"]
			if !ok {
				continue
			}
			f.MacInt, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.Mode:
			vvv, ok := vv["mode"]
			if !ok {
				continue
			}
			f.Mode, err = Enum(vvv, FailureEnum.Mode.MAP, FailureEnum.Mode.RMAP, FailureEnum.Mode.LIST)
			updatedCols = append(updatedCols, col)
		case FailureCol.Mac:
			vvv, ok := vv["mac"]
			if !ok {
				continue
			}
			f.Mac, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.Latest:
			vvv, ok := vv["latest"]
			if !ok {
				continue
			}
			f.Latest, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.raw:
			vvv, ok := vv["raw"]
			if !ok {
				continue
			}
			f.raw, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.IsFailed:
			vvv, ok := vv["is_failed"]
			if !ok {
				continue
			}
			f.IsFailed, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case FailureCol.FailureMode:
			vvv, ok := vv["failure_mode"]
			if !ok {
				continue
			}
			f.FailureMode, err = Enum(vvv, FailureEnum.FailureMode.MAP, FailureEnum.FailureMode.RMAP, FailureEnum.FailureMode.LIST)
			updatedCols = append(updatedCols, col)
		case FailureCol.FailureMode1:
			vvv, ok := vv["failure_mode1"]
			if !ok {
				continue
			}
			f.FailureMode1, err = Enum(vvv, FailureEnum.FailureMode1.MAP, FailureEnum.FailureMode1.RMAP, FailureEnum.FailureMode1.LIST)
			updatedCols = append(updatedCols, col)
		case FailureCol.FailureMode2:
			vvv, ok := vv["failure_mode2"]
			if !ok {
				continue
			}
			f.FailureMode2, err = Enum(vvv, FailureEnum.FailureMode2.MAP, FailureEnum.FailureMode2.RMAP, FailureEnum.FailureMode2.LIST)
			updatedCols = append(updatedCols, col)
		case FailureCol.FailureMode3:
			vvv, ok := vv["failure_mode3"]
			if !ok {
				continue
			}
			f.FailureMode3, err = Enum(vvv, FailureEnum.FailureMode3.MAP, FailureEnum.FailureMode3.RMAP, FailureEnum.FailureMode3.LIST)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalFailure(ctx context.Context, vi interface{}, cols ...db.Col) (*Failure, error) {
	f := NewFailure(ctx)
	_, err := f.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return f, err
}

func UnmarshalFailures(ctx context.Context, vi interface{}, cols ...db.Col) ([]*Failure, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalFailures with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalFailures")
	}
	failures := make([]*Failure, len(vv))
	for ii, vvv := range vv {
		var f *Failure
		f, err = UnmarshalFailure(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		failures[ii] = f
	}
	return failures, nil
}

func newFailureDest(cols ...string) (db.Doer, []interface{}, error) {
	f := &Failure{}
	if cols == nil || len(cols) == 0 {
		return f, []interface{}{&f.Id, &f.Created, &f.Updated, &f.Visibly, &f.FactoryId, &f.BatchId, &f.MacInt, &f.Mode, &f.Mac, &f.Latest, &f.raw, &f.IsFailed, &f.FailureMode, &f.FailureMode1, &f.FailureMode2, &f.FailureMode3}, nil
	}
	dest := make([]interface{}, len(cols))
	for ii, col := range cols {
		switch col {
		case "id":
			dest[ii] = &f.Id
		case "created":
			dest[ii] = &f.Created
		case "updated":
			dest[ii] = &f.Updated
		case "visibly":
			dest[ii] = &f.Visibly
		case "factory_id":
			dest[ii] = &f.FactoryId
		case "batch_id":
			dest[ii] = &f.BatchId
		case "mac_int":
			dest[ii] = &f.MacInt
		case "mode":
			dest[ii] = &f.Mode
		case "mac":
			dest[ii] = &f.Mac
		case "latest":
			dest[ii] = &f.Latest
		case "raw":
			dest[ii] = &f.raw
		case "is_failed":
			dest[ii] = &f.IsFailed
		case "failure_mode":
			dest[ii] = &f.FailureMode
		case "failure_mode1":
			dest[ii] = &f.FailureMode1
		case "failure_mode2":
			dest[ii] = &f.FailureMode2
		case "failure_mode3":
			dest[ii] = &f.FailureMode3
		default:
			return nil, nil, logError("dal.Failure Error: unknow column " + col + " in talbe failure")
		}
	}
	return f, dest, nil
}

func colsAndArgsFailure(f *Failure, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.Failure Error: at least one column to colsAndArgsFailure")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case FailureCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = f.Id
		case FailureCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = f.Created
		case FailureCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = f.Updated
		case FailureCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = f.Visibly
		case FailureCol.FactoryId:
			cols[ii] = "`factory_id` = ?"
			args[ii] = f.FactoryId
		case FailureCol.BatchId:
			cols[ii] = "`batch_id` = ?"
			args[ii] = f.BatchId
		case FailureCol.MacInt:
			cols[ii] = "`mac_int` = ?"
			args[ii] = f.MacInt
		case FailureCol.Mode:
			cols[ii] = "`mode` = ?"
			args[ii] = f.Mode
		case FailureCol.Mac:
			cols[ii] = "`mac` = ?"
			args[ii] = f.Mac
		case FailureCol.Latest:
			cols[ii] = "`latest` = ?"
			args[ii] = f.Latest
		case FailureCol.raw:
			cols[ii] = "`raw` = ?"
			args[ii] = f.raw
		case FailureCol.IsFailed:
			cols[ii] = "`is_failed` = ?"
			args[ii] = f.IsFailed
		case FailureCol.FailureMode:
			cols[ii] = "`failure_mode` = ?"
			args[ii] = f.FailureMode
		case FailureCol.FailureMode1:
			cols[ii] = "`failure_mode1` = ?"
			args[ii] = f.FailureMode1
		case FailureCol.FailureMode2:
			cols[ii] = "`failure_mode2` = ?"
			args[ii] = f.FailureMode2
		case FailureCol.FailureMode3:
			cols[ii] = "`failure_mode3` = ?"
			args[ii] = f.FailureMode3
		default:
			return nil, nil, logError(fmt.Sprintf("dal.Failure Error: unknow column num %d in talbe failure", cc))
		}
	}
	return cols, args, nil
}

var FailureEnum = struct {
	Mode struct {
		LIST              string
		MAP               map[string]int
		RMAP              map[int]string
		ESP8266, ESP32, _ int
	}
	FailureMode struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}
	FailureMode1 struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}
	FailureMode2 struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}
	FailureMode3 struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}
}{
	struct {
		LIST              string
		MAP               map[string]int
		RMAP              map[int]string
		ESP8266, ESP32, _ int
	}{"ESP8266/ESP32", map[string]int{"ESP8266": 1, "ESP32": 2}, map[int]string{1: "ESP8266", 2: "ESP32"}, 1, 2, 0},
	struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}{"DUT_RXRSSI/FB_RXRSSI/FREQ_OFFSET/RXDC/RXIQ/RX_NOISEFLOOR/TXDC/TXIQ/TX_POWER_BACKOFF/TXP_RESULT/TX_VDD33/VDD33", map[string]int{"DUT_RXRSSI": 1, "FB_RXRSSI": 2, "FREQ_OFFSET": 3, "RXDC": 4, "RXIQ": 5, "RX_NOISEFLOOR": 6, "TXDC": 7, "TXIQ": 8, "TX_POWER_BACKOFF": 9, "TXP_RESULT": 10, "TX_VDD33": 11, "VDD33": 12}, map[int]string{1: "DUT_RXRSSI", 2: "FB_RXRSSI", 3: "FREQ_OFFSET", 4: "RXDC", 5: "RXIQ", 6: "RX_NOISEFLOOR", 7: "TXDC", 8: "TXIQ", 9: "TX_POWER_BACKOFF", 10: "TXP_RESULT", 11: "TX_VDD33", 12: "VDD33"}, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0},
	struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}{"DUT_RXRSSI/FB_RXRSSI/FREQ_OFFSET/RXDC/RXIQ/RX_NOISEFLOOR/TXDC/TXIQ/TX_POWER_BACKOFF/TXP_RESULT/TX_VDD33/VDD33", map[string]int{"DUT_RXRSSI": 1, "FB_RXRSSI": 2, "FREQ_OFFSET": 3, "RXDC": 4, "RXIQ": 5, "RX_NOISEFLOOR": 6, "TXDC": 7, "TXIQ": 8, "TX_POWER_BACKOFF": 9, "TXP_RESULT": 10, "TX_VDD33": 11, "VDD33": 12}, map[int]string{1: "DUT_RXRSSI", 2: "FB_RXRSSI", 3: "FREQ_OFFSET", 4: "RXDC", 5: "RXIQ", 6: "RX_NOISEFLOOR", 7: "TXDC", 8: "TXIQ", 9: "TX_POWER_BACKOFF", 10: "TXP_RESULT", 11: "TX_VDD33", 12: "VDD33"}, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0},
	struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}{"DUT_RXRSSI/FB_RXRSSI/FREQ_OFFSET/RXDC/RXIQ/RX_NOISEFLOOR/TXDC/TXIQ/TX_POWER_BACKOFF/TXP_RESULT/TX_VDD33/VDD33", map[string]int{"DUT_RXRSSI": 1, "FB_RXRSSI": 2, "FREQ_OFFSET": 3, "RXDC": 4, "RXIQ": 5, "RX_NOISEFLOOR": 6, "TXDC": 7, "TXIQ": 8, "TX_POWER_BACKOFF": 9, "TXP_RESULT": 10, "TX_VDD33": 11, "VDD33": 12}, map[int]string{1: "DUT_RXRSSI", 2: "FB_RXRSSI", 3: "FREQ_OFFSET", 4: "RXDC", 5: "RXIQ", 6: "RX_NOISEFLOOR", 7: "TXDC", 8: "TXIQ", 9: "TX_POWER_BACKOFF", 10: "TXP_RESULT", 11: "TX_VDD33", 12: "VDD33"}, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0},
	struct {
		LIST                                                                                                                        string
		MAP                                                                                                                         map[string]int
		RMAP                                                                                                                        map[int]string
		DUT_RXRSSI, FB_RXRSSI, FREQ_OFFSET, RXDC, RXIQ, RX_NOISEFLOOR, TXDC, TXIQ, TX_POWER_BACKOFF, TXP_RESULT, TX_VDD33, VDD33, _ int
	}{"DUT_RXRSSI/FB_RXRSSI/FREQ_OFFSET/RXDC/RXIQ/RX_NOISEFLOOR/TXDC/TXIQ/TX_POWER_BACKOFF/TXP_RESULT/TX_VDD33/VDD33", map[string]int{"DUT_RXRSSI": 1, "FB_RXRSSI": 2, "FREQ_OFFSET": 3, "RXDC": 4, "RXIQ": 5, "RX_NOISEFLOOR": 6, "TXDC": 7, "TXIQ": 8, "TX_POWER_BACKOFF": 9, "TXP_RESULT": 10, "TX_VDD33": 11, "VDD33": 12}, map[int]string{1: "DUT_RXRSSI", 2: "FB_RXRSSI", 3: "FREQ_OFFSET", 4: "RXDC", 5: "RXIQ", 6: "RX_NOISEFLOOR", 7: "TXDC", 8: "TXIQ", 9: "TX_POWER_BACKOFF", 10: "TXP_RESULT", 11: "TX_VDD33", 12: "VDD33"}, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0},
}

var failureSqls = []string{
	/*
		CREATE TABLE `failure` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `factory_id` int(11) NOT NULL,
		  `batch_id` int(11) NOT NULL,
		  `mac_int` int(11) NOT NULL,
		  `mode` int(11) NOT NULL,
		  `mac` varchar(0) NOT NULL,
		  `latest` bool NOT NULL,
		  `raw` varchar(4096) NOT NULL,
		  `is_failed` bool NOT NULL,
		  `failure_mode` int(11) NOT NULL,
		  `failure_mode1` int(11) NOT NULL,
		  `failure_mode2` int(11) NOT NULL,
		  `failure_mode3` int(11) NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into failure(`created`, `updated`, `visibly`, `factory_id`, `batch_id`, `mac_int`, `mode`, `mac`, `latest`, `raw`, `is_failed`, `failure_mode`, `failure_mode1`, `failure_mode2`, `failure_mode3`) values(now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*1*/ "insert into failure(`id`, `created`, `updated`, `visibly`, `factory_id`, `batch_id`, `mac_int`, `mode`, `mac`, `latest`, `raw`, `is_failed`, `failure_mode`, `failure_mode1`, `failure_mode2`, `failure_mode3`) values(?, now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*2*/ "update failure set updated = now(), `visibly` = ?, `factory_id` = ?, `batch_id` = ?, `mac_int` = ?, `mode` = ?, `mac` = ?, `latest` = ?, `raw` = ?, `is_failed` = ?, `failure_mode` = ?, `failure_mode1` = ?, `failure_mode2` = ?, `failure_mode3` = ? where id = ?",
	/*3*/ "update failure set updated = now(), %s where id = ?",
	/*4*/ "update failure set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from failure where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `factory_id`, `batch_id`, `mac_int`, `mode`, `mac`, `latest`, `raw`, `is_failed`, `failure_mode`, `failure_mode1`, `failure_mode2`, `failure_mode3` from failure where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `factory_id`, `batch_id`, `mac_int`, `mode`, `mac`, `latest`, `raw`, `is_failed`, `failure_mode`, `failure_mode1`, `failure_mode2`, `failure_mode3` from failure where id in (%s) and visibly = 1",
}

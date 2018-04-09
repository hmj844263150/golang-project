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

var BatchStatsTid = 8
var _ db.Doer = (*BatchStats)(nil)
var batchStatscols = []db.Col{1, 2, 3, 4, 5, 6, 7, 16, 8, 9, 10, 11, 12, 13, 14, 15}
var batchStatsfields = []string{"id", "created", "updated", "visibly", "batch_id", "start", "end", "cnt", "success", "success_pct", "right_first_time", "right_first_time_pct", "failed", "failed_pct", "rejected", "rejected_pct"}

var BatchStatsCol = struct {
	Id, Created, Updated, Visibly, BatchId, Start, End, Cnt, Success, SuccessPct, RightFirstTime, RightFirstTimePct, Failed, FailedPct, Rejected, RejectedPct, _ db.Col
}{1, 2, 3, 4, 5, 6, 7, 16, 8, 9, 10, 11, 12, 13, 14, 15, 0}

type BatchStats struct {
	Id                int
	Created           time.Time
	Updated           time.Time
	Visibly           bool
	BatchId           int
	Start             time.Time
	End               time.Time
	Cnt               int
	Success           int
	SuccessPct        int
	RightFirstTime    int
	RightFirstTimePct int
	Failed            int
	FailedPct         int
	Rejected          int
	RejectedPct       int

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewBatchStats(ctx context.Context) *BatchStats {
	now := time.Now()
	b := &BatchStats{Created: now, Updated: now, Visibly: true}
	b.ext = GetExtFromContext(ctx)
	defaultBatchStats(ctx, b)
	return b
}

func FindBatchStats(ctx context.Context, id int) *BatchStats {
	dos, err := db.Open("BatchStats").Query(newBatchStatsDest, true, batchStatsSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		batchStats, _ := do.(*BatchStats)
		if ext != nil {
			batchStats.ext = ext
		}
		return batchStats
	}
	return nil
}

func ListBatchStats(ctx context.Context, ids ...int) []*BatchStats {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(batchStatsSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("BatchStats").Query(newBatchStatsDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	batchStatss := make([]*BatchStats, len(dos))
	for ii, do := range dos {
		batchStats, _ := do.(*BatchStats)
		if ext != nil {
			batchStats.ext = ext
		}
		batchStatss[ii] = batchStats
	}
	return batchStatss
}

func FindBatchStatsByBatchIdStartEnd(ctx context.Context, batchId int, start time.Time, end time.Time) *BatchStats {
	dos, err := db.Open("BatchStats").Query(newBatchStatsDest, true, batchStatsSqls[8], batchId, start, end)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		batchStats, _ := do.(*BatchStats)
		if ext != nil {
			batchStats.ext = ext
		}
		return batchStats
	}
	return nil
}

func (b *BatchStats) Save() error {
	now := time.Now()
	b.Created, b.Updated, b.Visibly = now, now, true
	var id int64
	var err error
	if b.Id == 0 {
		id, _, err = db.Open("BatchStats").Exec(batchStatsSqls[0], b.BatchId, b.Start, b.End, b.Cnt, b.Success, b.SuccessPct, b.RightFirstTime, b.RightFirstTimePct, b.Failed, b.FailedPct, b.Rejected, b.RejectedPct)
	} else {
		id, _, err = db.Open("BatchStats").Exec(batchStatsSqls[1], b.Id, b.BatchId, b.Start, b.End, b.Cnt, b.Success, b.SuccessPct, b.RightFirstTime, b.RightFirstTimePct, b.Failed, b.FailedPct, b.Rejected, b.RejectedPct)
	}
	if err != nil {
		return err
	}
	b.Id = int(id)
	return nil
}

func (b *BatchStats) Update(cs ...db.Col) error {
	if b.Id == 0 {
		return logError("dal.BatchStats Error: can not update row while id is zero")
	}
	b.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("BatchStats").Exec(batchStatsSqls[2], b.Visibly, b.BatchId, b.Start, b.End, b.Cnt, b.Success, b.SuccessPct, b.RightFirstTime, b.RightFirstTimePct, b.Failed, b.FailedPct, b.Rejected, b.RejectedPct, b.Id)
		return err
	}
	cols, args, err := colsAndArgsBatchStats(b, cs...)
	if err != nil {
		return err
	}
	args = append(args, b.Id)
	sqlstr := fmt.Sprintf(batchStatsSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("BatchStats").Exec(sqlstr, args...)
	return err
}

func (b *BatchStats) Invisibly() error {
	if b.Id == 0 {
		return logError("dal.BatchStats Error: can not invisibly row while id is zero")
	}
	b.Updated = time.Now()
	b.Visibly = false
	_, _, err := db.Open("BatchStats").Exec(batchStatsSqls[4], b.Id)
	return err
}

func (b *BatchStats) Delete() error {
	if b.Id == 0 {
		return logError("dal.BatchStats Error: can not delete row while id is zero")
	}
	b.Updated = time.Now()
	_, _, err := db.Open("BatchStats").Exec(batchStatsSqls[5], b.Id)
	return err
}

func (b *BatchStats) Valid() error {
	return b.valid()
}

func (b *BatchStats) SetExt(ext *Ext) {
	b.ext = ext
}

func (b *BatchStats) Padding(pkey string, pvalue interface{}) {
	if b.ext == nil {
		b.ext = &Ext{Loc: DefaultLoc}
	}
	if b.paddings == nil {
		b.paddings = make(map[string]interface{})
	}
	b.paddings[pkey] = pvalue
	b.ext.IsComplex = true
}

func (b *BatchStats) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case BatchStatsCol.Id:
			if isColumnName {
				mm["id"] = b.Id
			} else {
				mm["Id"] = b.Id
			}
		case BatchStatsCol.Created:
			if isColumnName {
				mm["created"] = b.Created
			} else {
				mm["Created"] = b.Created
			}
		case BatchStatsCol.Updated:
			if isColumnName {
				mm["updated"] = b.Updated
			} else {
				mm["Updated"] = b.Updated
			}
		case BatchStatsCol.Visibly:
			if isColumnName {
				mm["visibly"] = b.Visibly
			} else {
				mm["Visibly"] = b.Visibly
			}
		case BatchStatsCol.BatchId:
			if isColumnName {
				mm["batch_id"] = b.BatchId
			} else {
				mm["BatchId"] = b.BatchId
			}
		case BatchStatsCol.Start:
			if isColumnName {
				mm["start"] = b.Start
			} else {
				mm["Start"] = b.Start
			}
		case BatchStatsCol.End:
			if isColumnName {
				mm["end"] = b.End
			} else {
				mm["End"] = b.End
			}
		case BatchStatsCol.Cnt:
			if isColumnName {
				mm["cnt"] = b.Cnt
			} else {
				mm["Cnt"] = b.Cnt
			}
		case BatchStatsCol.Success:
			if isColumnName {
				mm["success"] = b.Success
			} else {
				mm["Success"] = b.Success
			}
		case BatchStatsCol.SuccessPct:
			if isColumnName {
				mm["success_pct"] = b.SuccessPct
			} else {
				mm["SuccessPct"] = b.SuccessPct
			}
		case BatchStatsCol.RightFirstTime:
			if isColumnName {
				mm["right_first_time"] = b.RightFirstTime
			} else {
				mm["RightFirstTime"] = b.RightFirstTime
			}
		case BatchStatsCol.RightFirstTimePct:
			if isColumnName {
				mm["right_first_time_pct"] = b.RightFirstTimePct
			} else {
				mm["RightFirstTimePct"] = b.RightFirstTimePct
			}
		case BatchStatsCol.Failed:
			if isColumnName {
				mm["failed"] = b.Failed
			} else {
				mm["Failed"] = b.Failed
			}
		case BatchStatsCol.FailedPct:
			if isColumnName {
				mm["failed_pct"] = b.FailedPct
			} else {
				mm["FailedPct"] = b.FailedPct
			}
		case BatchStatsCol.Rejected:
			if isColumnName {
				mm["rejected"] = b.Rejected
			} else {
				mm["Rejected"] = b.Rejected
			}
		case BatchStatsCol.RejectedPct:
			if isColumnName {
				mm["rejected_pct"] = b.RejectedPct
			} else {
				mm["RejectedPct"] = b.RejectedPct
			}
		default:
			logError(fmt.Sprintf("dal.BatchStats Error: unknow column num %d in talbe batch_stats", cc))
		}
	}
	return mm
}

func (b *BatchStats) MarshalJSON() ([]byte, error) {
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
	buf.WriteString(`, "batch_id":`)
	buf.WriteString(strconv.FormatInt(int64(b.BatchId), 10))
	buf.WriteString(`, "start":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(b.Start.Unix(), 10))
	} else {
		b.Start = b.Start.In(loc)
		buf.WriteString(`"` + b.Start.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "end":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(b.End.Unix(), 10))
	} else {
		b.End = b.End.In(loc)
		buf.WriteString(`"` + b.End.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "cnt":`)
	buf.WriteString(strconv.FormatInt(int64(b.Cnt), 10))
	buf.WriteString(`, "success":`)
	buf.WriteString(strconv.FormatInt(int64(b.Success), 10))
	buf.WriteString(`, "success_pct":`)
	buf.WriteString(strconv.FormatInt(int64(b.SuccessPct), 10))
	buf.WriteString(`, "right_first_time":`)
	buf.WriteString(strconv.FormatInt(int64(b.RightFirstTime), 10))
	buf.WriteString(`, "right_first_time_pct":`)
	buf.WriteString(strconv.FormatInt(int64(b.RightFirstTimePct), 10))
	buf.WriteString(`, "failed":`)
	buf.WriteString(strconv.FormatInt(int64(b.Failed), 10))
	buf.WriteString(`, "failed_pct":`)
	buf.WriteString(strconv.FormatInt(int64(b.FailedPct), 10))
	buf.WriteString(`, "rejected":`)
	buf.WriteString(strconv.FormatInt(int64(b.Rejected), 10))
	buf.WriteString(`, "rejected_pct":`)
	buf.WriteString(strconv.FormatInt(int64(b.RejectedPct), 10))
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (b *BatchStats) marshalJSONComplex() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	if b.ext == nil {
		return nil, logError("dal.BatchStats Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := b.ext.Loc
	numericEnum := b.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[BatchStatsTid]; ok {
		if vvv, ok := vv[b.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := b.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(BatchStatsCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(b.Id), 10))
	}
	isRender = isRenderField(BatchStatsCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.Created.Unix(), 10))
		} else {
			b.Created = b.Created.In(loc)
			buf.WriteString(`"` + b.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(BatchStatsCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.Updated.Unix(), 10))
		} else {
			b.Updated = b.Updated.In(loc)
			buf.WriteString(`"` + b.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(BatchStatsCol.Visibly, "visibly", includes, excludes, paddings)
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
	isRender = isRenderField(BatchStatsCol.BatchId, "batch_id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "batch_id":`)
		buf.WriteString(strconv.FormatInt(int64(b.BatchId), 10))
	}
	isRender = isRenderField(BatchStatsCol.Start, "start", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "start":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.Start.Unix(), 10))
		} else {
			b.Start = b.Start.In(loc)
			buf.WriteString(`"` + b.Start.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(BatchStatsCol.End, "end", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "end":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(b.End.Unix(), 10))
		} else {
			b.End = b.End.In(loc)
			buf.WriteString(`"` + b.End.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(BatchStatsCol.Cnt, "cnt", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "cnt":`)
		buf.WriteString(strconv.FormatInt(int64(b.Cnt), 10))
	}
	isRender = isRenderField(BatchStatsCol.Success, "success", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "success":`)
		buf.WriteString(strconv.FormatInt(int64(b.Success), 10))
	}
	isRender = isRenderField(BatchStatsCol.SuccessPct, "success_pct", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "success_pct":`)
		buf.WriteString(strconv.FormatInt(int64(b.SuccessPct), 10))
	}
	isRender = isRenderField(BatchStatsCol.RightFirstTime, "right_first_time", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "right_first_time":`)
		buf.WriteString(strconv.FormatInt(int64(b.RightFirstTime), 10))
	}
	isRender = isRenderField(BatchStatsCol.RightFirstTimePct, "right_first_time_pct", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "right_first_time_pct":`)
		buf.WriteString(strconv.FormatInt(int64(b.RightFirstTimePct), 10))
	}
	isRender = isRenderField(BatchStatsCol.Failed, "failed", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failed":`)
		buf.WriteString(strconv.FormatInt(int64(b.Failed), 10))
	}
	isRender = isRenderField(BatchStatsCol.FailedPct, "failed_pct", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "failed_pct":`)
		buf.WriteString(strconv.FormatInt(int64(b.FailedPct), 10))
	}
	isRender = isRenderField(BatchStatsCol.Rejected, "rejected", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "rejected":`)
		buf.WriteString(strconv.FormatInt(int64(b.Rejected), 10))
	}
	isRender = isRenderField(BatchStatsCol.RejectedPct, "rejected_pct", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "rejected_pct":`)
		buf.WriteString(strconv.FormatInt(int64(b.RejectedPct), 10))
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

func (b *BatchStats) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalBatchStats with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalBatchStats")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = batchStatscols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case BatchStatsCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			b.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			b.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			b.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			b.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.BatchId:
			vvv, ok := vv["batch_id"]
			if !ok {
				continue
			}
			b.BatchId, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Start:
			vvv, ok := vv["start"]
			if !ok {
				continue
			}
			b.Start, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.End:
			vvv, ok := vv["end"]
			if !ok {
				continue
			}
			b.End, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Cnt:
			vvv, ok := vv["cnt"]
			if !ok {
				continue
			}
			b.Cnt, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Success:
			vvv, ok := vv["success"]
			if !ok {
				continue
			}
			b.Success, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.SuccessPct:
			vvv, ok := vv["success_pct"]
			if !ok {
				continue
			}
			b.SuccessPct, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.RightFirstTime:
			vvv, ok := vv["right_first_time"]
			if !ok {
				continue
			}
			b.RightFirstTime, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.RightFirstTimePct:
			vvv, ok := vv["right_first_time_pct"]
			if !ok {
				continue
			}
			b.RightFirstTimePct, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Failed:
			vvv, ok := vv["failed"]
			if !ok {
				continue
			}
			b.Failed, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.FailedPct:
			vvv, ok := vv["failed_pct"]
			if !ok {
				continue
			}
			b.FailedPct, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.Rejected:
			vvv, ok := vv["rejected"]
			if !ok {
				continue
			}
			b.Rejected, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case BatchStatsCol.RejectedPct:
			vvv, ok := vv["rejected_pct"]
			if !ok {
				continue
			}
			b.RejectedPct, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalBatchStats(ctx context.Context, vi interface{}, cols ...db.Col) (*BatchStats, error) {
	b := NewBatchStats(ctx)
	_, err := b.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return b, err
}

func UnmarshalBatchStatss(ctx context.Context, vi interface{}, cols ...db.Col) ([]*BatchStats, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalBatchStatss with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalBatchStatss")
	}
	batchStatss := make([]*BatchStats, len(vv))
	for ii, vvv := range vv {
		var b *BatchStats
		b, err = UnmarshalBatchStats(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		batchStatss[ii] = b
	}
	return batchStatss, nil
}

func newBatchStatsDest(cols ...string) (db.Doer, []interface{}, error) {
	b := &BatchStats{}
	if cols == nil || len(cols) == 0 {
		return b, []interface{}{&b.Id, &b.Created, &b.Updated, &b.Visibly, &b.BatchId, &b.Start, &b.End, &b.Cnt, &b.Success, &b.SuccessPct, &b.RightFirstTime, &b.RightFirstTimePct, &b.Failed, &b.FailedPct, &b.Rejected, &b.RejectedPct}, nil
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
		case "batch_id":
			dest[ii] = &b.BatchId
		case "start":
			dest[ii] = &b.Start
		case "end":
			dest[ii] = &b.End
		case "cnt":
			dest[ii] = &b.Cnt
		case "success":
			dest[ii] = &b.Success
		case "success_pct":
			dest[ii] = &b.SuccessPct
		case "right_first_time":
			dest[ii] = &b.RightFirstTime
		case "right_first_time_pct":
			dest[ii] = &b.RightFirstTimePct
		case "failed":
			dest[ii] = &b.Failed
		case "failed_pct":
			dest[ii] = &b.FailedPct
		case "rejected":
			dest[ii] = &b.Rejected
		case "rejected_pct":
			dest[ii] = &b.RejectedPct
		default:
			return nil, nil, logError("dal.BatchStats Error: unknow column " + col + " in talbe batch_stats")
		}
	}
	return b, dest, nil
}

func colsAndArgsBatchStats(b *BatchStats, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.BatchStats Error: at least one column to colsAndArgsBatchStats")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case BatchStatsCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = b.Id
		case BatchStatsCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = b.Created
		case BatchStatsCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = b.Updated
		case BatchStatsCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = b.Visibly
		case BatchStatsCol.BatchId:
			cols[ii] = "`batch_id` = ?"
			args[ii] = b.BatchId
		case BatchStatsCol.Start:
			cols[ii] = "`start` = ?"
			args[ii] = b.Start
		case BatchStatsCol.End:
			cols[ii] = "`end` = ?"
			args[ii] = b.End
		case BatchStatsCol.Cnt:
			cols[ii] = "`cnt` = ?"
			args[ii] = b.Cnt
		case BatchStatsCol.Success:
			cols[ii] = "`success` = ?"
			args[ii] = b.Success
		case BatchStatsCol.SuccessPct:
			cols[ii] = "`success_pct` = ?"
			args[ii] = b.SuccessPct
		case BatchStatsCol.RightFirstTime:
			cols[ii] = "`right_first_time` = ?"
			args[ii] = b.RightFirstTime
		case BatchStatsCol.RightFirstTimePct:
			cols[ii] = "`right_first_time_pct` = ?"
			args[ii] = b.RightFirstTimePct
		case BatchStatsCol.Failed:
			cols[ii] = "`failed` = ?"
			args[ii] = b.Failed
		case BatchStatsCol.FailedPct:
			cols[ii] = "`failed_pct` = ?"
			args[ii] = b.FailedPct
		case BatchStatsCol.Rejected:
			cols[ii] = "`rejected` = ?"
			args[ii] = b.Rejected
		case BatchStatsCol.RejectedPct:
			cols[ii] = "`rejected_pct` = ?"
			args[ii] = b.RejectedPct
		default:
			return nil, nil, logError(fmt.Sprintf("dal.BatchStats Error: unknow column num %d in talbe batch_stats", cc))
		}
	}
	return cols, args, nil
}

var BatchStatsEnum = struct {
}{}

var batchStatsSqls = []string{
	/*
		CREATE TABLE `batch_stats` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `batch_id` int(11) NOT NULL,
		  `start` datetime NOT NULL,
		  `end` datetime NOT NULL,
		  `cnt` int(11) NOT NULL,
		  `success` int(11) NOT NULL,
		  `success_pct` int(11) NOT NULL,
		  `right_first_time` int(11) NOT NULL,
		  `right_first_time_pct` int(11) NOT NULL,
		  `failed` int(11) NOT NULL,
		  `failed_pct` int(11) NOT NULL,
		  `rejected` int(11) NOT NULL,
		  `rejected_pct` int(11) NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into batch_stats(`created`, `updated`, `visibly`, `batch_id`, `start`, `end`, `cnt`, `success`, `success_pct`, `right_first_time`, `right_first_time_pct`, `failed`, `failed_pct`, `rejected`, `rejected_pct`) values(now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*1*/ "insert into batch_stats(`id`, `created`, `updated`, `visibly`, `batch_id`, `start`, `end`, `cnt`, `success`, `success_pct`, `right_first_time`, `right_first_time_pct`, `failed`, `failed_pct`, `rejected`, `rejected_pct`) values(?, now(), now(), 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	/*2*/ "update batch_stats set updated = now(), `visibly` = ?, `batch_id` = ?, `start` = ?, `end` = ?, `cnt` = ?, `success` = ?, `success_pct` = ?, `right_first_time` = ?, `right_first_time_pct` = ?, `failed` = ?, `failed_pct` = ?, `rejected` = ?, `rejected_pct` = ? where id = ?",
	/*3*/ "update batch_stats set updated = now(), %s where id = ?",
	/*4*/ "update batch_stats set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from batch_stats where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `batch_id`, `start`, `end`, `cnt`, `success`, `success_pct`, `right_first_time`, `right_first_time_pct`, `failed`, `failed_pct`, `rejected`, `rejected_pct` from batch_stats where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `batch_id`, `start`, `end`, `cnt`, `success`, `success_pct`, `right_first_time`, `right_first_time_pct`, `failed`, `failed_pct`, `rejected`, `rejected_pct` from batch_stats where id in (%s) and visibly = 1",

	/*8*/ "select `id`, `created`, `updated`, `visibly`, `batch_id`, `start`, `end`, `cnt`, `success`, `success_pct`, `right_first_time`, `right_first_time_pct`, `failed`, `failed_pct`, `rejected`, `rejected_pct` from batch_stats where visibly = 1 and batch_id = ? and start = ? and end = ?",
}

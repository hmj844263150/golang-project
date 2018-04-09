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

var FactoryTid = 1
var _ db.Doer = (*Factory)(nil)
var factorycols = []db.Col{1, 2, 3, 4, 5, 6, 7, 8, 9}
var factoryfields = []string{"id", "created", "updated", "visibly", "sid", "name", "location", "token", "is_staff"}

var FactoryCol = struct {
	Id, Created, Updated, Visibly, Sid, Name, Location, Token, IsStaff, _ db.Col
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

type Factory struct {
	Id       int
	Created  time.Time
	Updated  time.Time
	Visibly  bool
	Sid      string
	Name     string
	Location string
	Token    string
	IsStaff  bool

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewFactory(ctx context.Context) *Factory {
	now := time.Now()
	f := &Factory{Created: now, Updated: now, Visibly: true}
	f.ext = GetExtFromContext(ctx)
	defaultFactory(ctx, f)
	return f
}

func FindFactory(ctx context.Context, id int) *Factory {
	dos, err := db.Open("Factory").Query(newFactoryDest, true, factorySqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		return factory
	}
	return nil
}

func ListFactory(ctx context.Context, ids ...int) []*Factory {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(factorySqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("Factory").Query(newFactoryDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	factorys := make([]*Factory, len(dos))
	for ii, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		factorys[ii] = factory
	}
	return factorys
}

func FindFactoryBySid(ctx context.Context, sid string) *Factory {
	dos, err := db.Open("Factory").Query(newFactoryDest, true, factorySqls[8], sid)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		return factory
	}
	return nil
}

func ListFactoryAll(ctx context.Context, offset int, rowCount int) []*Factory {
	dos, err := db.Open("Factory").Query(newFactoryDest, true, factorySqls[9], offset, rowCount)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	factorys := make([]*Factory, len(dos))
	for ii, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		factorys[ii] = factory
	}
	return factorys
}

func FindFactoryByToken(ctx context.Context, token string) *Factory {
	dos, err := db.Open("Factory").Query(newFactoryDest, true, factorySqls[10], token)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		return factory
	}
	return nil
}

func FindFactoryByName(ctx context.Context, name string) *Factory {
	dos, err := db.Open("Factory").Query(newFactoryDest, true, factorySqls[11], name)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		return factory
	}
	return nil
}

func FindFactoryStaff(ctx context.Context) *Factory {
	dos, err := db.Open("Factory").Query(newFactoryDest, true, factorySqls[12])
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		factory, _ := do.(*Factory)
		if ext != nil {
			factory.ext = ext
		}
		return factory
	}
	return nil
}

func (f *Factory) Save() error {
	now := time.Now()
	f.Created, f.Updated, f.Visibly = now, now, true
	var id int64
	var err error
	if f.Id == 0 {
		id, _, err = db.Open("Factory").Exec(factorySqls[0], f.Sid, f.Name, f.Location, f.Token, f.IsStaff)
	} else {
		id, _, err = db.Open("Factory").Exec(factorySqls[1], f.Id, f.Sid, f.Name, f.Location, f.Token, f.IsStaff)
	}
	if err != nil {
		return err
	}
	f.Id = int(id)
	return nil
}

func (f *Factory) Update(cs ...db.Col) error {
	if f.Id == 0 {
		return logError("dal.Factory Error: can not update row while id is zero")
	}
	f.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("Factory").Exec(factorySqls[2], f.Visibly, f.Sid, f.Name, f.Location, f.Token, f.IsStaff, f.Id)
		return err
	}
	cols, args, err := colsAndArgsFactory(f, cs...)
	if err != nil {
		return err
	}
	args = append(args, f.Id)
	sqlstr := fmt.Sprintf(factorySqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("Factory").Exec(sqlstr, args...)
	return err
}

func (f *Factory) Invisibly() error {
	if f.Id == 0 {
		return logError("dal.Factory Error: can not invisibly row while id is zero")
	}
	f.Updated = time.Now()
	f.Visibly = false
	_, _, err := db.Open("Factory").Exec(factorySqls[4], f.Id)
	return err
}

func (f *Factory) Delete() error {
	if f.Id == 0 {
		return logError("dal.Factory Error: can not delete row while id is zero")
	}
	f.Updated = time.Now()
	_, _, err := db.Open("Factory").Exec(factorySqls[5], f.Id)
	return err
}

func (f *Factory) Valid() error {
	return f.valid()
}

func (f *Factory) SetExt(ext *Ext) {
	f.ext = ext
}

func (f *Factory) Padding(pkey string, pvalue interface{}) {
	if f.ext == nil {
		f.ext = &Ext{Loc: DefaultLoc}
	}
	if f.paddings == nil {
		f.paddings = make(map[string]interface{})
	}
	f.paddings[pkey] = pvalue
	f.ext.IsComplex = true
}

func (f *Factory) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case FactoryCol.Id:
			if isColumnName {
				mm["id"] = f.Id
			} else {
				mm["Id"] = f.Id
			}
		case FactoryCol.Created:
			if isColumnName {
				mm["created"] = f.Created
			} else {
				mm["Created"] = f.Created
			}
		case FactoryCol.Updated:
			if isColumnName {
				mm["updated"] = f.Updated
			} else {
				mm["Updated"] = f.Updated
			}
		case FactoryCol.Visibly:
			if isColumnName {
				mm["visibly"] = f.Visibly
			} else {
				mm["Visibly"] = f.Visibly
			}
		case FactoryCol.Sid:
			if isColumnName {
				mm["sid"] = f.Sid
			} else {
				mm["Sid"] = f.Sid
			}
		case FactoryCol.Name:
			if isColumnName {
				mm["name"] = f.Name
			} else {
				mm["Name"] = f.Name
			}
		case FactoryCol.Location:
			if isColumnName {
				mm["location"] = f.Location
			} else {
				mm["Location"] = f.Location
			}
		case FactoryCol.Token:
			if isColumnName {
				mm["token"] = f.Token
			} else {
				mm["Token"] = f.Token
			}
		case FactoryCol.IsStaff:
			if isColumnName {
				mm["is_staff"] = f.IsStaff
			} else {
				mm["IsStaff"] = f.IsStaff
			}
		default:
			logError(fmt.Sprintf("dal.Factory Error: unknow column num %d in talbe factory", cc))
		}
	}
	return mm
}

func (f *Factory) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}
	loc := DefaultLoc
	var numericEnum bool
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
	buf.WriteString(`, "sid":`)
	WriteJsonString(&buf, f.Sid)
	buf.WriteString(`, "name":`)
	WriteJsonString(&buf, f.Name)
	buf.WriteString(`, "location":`)
	WriteJsonString(&buf, f.Location)
	buf.WriteString(`, "token":`)
	WriteJsonString(&buf, f.Token)
	buf.WriteString(`, "is_staff":`)
	if numericEnum {
		if f.IsStaff {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(f.IsStaff))
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (f *Factory) marshalJSONComplex() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}
	if f.ext == nil {
		return nil, logError("dal.Factory Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := f.ext.Loc
	numericEnum := f.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[FactoryTid]; ok {
		if vvv, ok := vv[f.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := f.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(FactoryCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(f.Id), 10))
	}
	isRender = isRenderField(FactoryCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(f.Created.Unix(), 10))
		} else {
			f.Created = f.Created.In(loc)
			buf.WriteString(`"` + f.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(FactoryCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(f.Updated.Unix(), 10))
		} else {
			f.Updated = f.Updated.In(loc)
			buf.WriteString(`"` + f.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(FactoryCol.Visibly, "visibly", includes, excludes, paddings)
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
	isRender = isRenderField(FactoryCol.Sid, "sid", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "sid":`)
		WriteJsonString(&buf, f.Sid)
	}
	isRender = isRenderField(FactoryCol.Name, "name", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "name":`)
		WriteJsonString(&buf, f.Name)
	}
	isRender = isRenderField(FactoryCol.Location, "location", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "location":`)
		WriteJsonString(&buf, f.Location)
	}
	isRender = isRenderField(FactoryCol.Token, "token", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "token":`)
		WriteJsonString(&buf, f.Token)
	}
	isRender = isRenderField(FactoryCol.IsStaff, "is_staff", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "is_staff":`)
		if numericEnum {
			if f.IsStaff {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(f.IsStaff))
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

func (f *Factory) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalFactory with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalFactory")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = factorycols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case FactoryCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			f.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			f.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			f.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			f.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Sid:
			vvv, ok := vv["sid"]
			if !ok {
				continue
			}
			f.Sid, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Name:
			vvv, ok := vv["name"]
			if !ok {
				continue
			}
			f.Name, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Location:
			vvv, ok := vv["location"]
			if !ok {
				continue
			}
			f.Location, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case FactoryCol.Token:
			vvv, ok := vv["token"]
			if !ok {
				continue
			}
			f.Token, err = String(vvv)
			updatedCols = append(updatedCols, col)
		case FactoryCol.IsStaff:
			vvv, ok := vv["is_staff"]
			if !ok {
				continue
			}
			f.IsStaff, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalFactory(ctx context.Context, vi interface{}, cols ...db.Col) (*Factory, error) {
	f := NewFactory(ctx)
	_, err := f.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return f, err
}

func UnmarshalFactorys(ctx context.Context, vi interface{}, cols ...db.Col) ([]*Factory, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalFactorys with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalFactorys")
	}
	factorys := make([]*Factory, len(vv))
	for ii, vvv := range vv {
		var f *Factory
		f, err = UnmarshalFactory(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		factorys[ii] = f
	}
	return factorys, nil
}

func newFactoryDest(cols ...string) (db.Doer, []interface{}, error) {
	f := &Factory{}
	if cols == nil || len(cols) == 0 {
		return f, []interface{}{&f.Id, &f.Created, &f.Updated, &f.Visibly, &f.Sid, &f.Name, &f.Location, &f.Token, &f.IsStaff}, nil
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
		case "sid":
			dest[ii] = &f.Sid
		case "name":
			dest[ii] = &f.Name
		case "location":
			dest[ii] = &f.Location
		case "token":
			dest[ii] = &f.Token
		case "is_staff":
			dest[ii] = &f.IsStaff
		default:
			return nil, nil, logError("dal.Factory Error: unknow column " + col + " in talbe factory")
		}
	}
	return f, dest, nil
}

func colsAndArgsFactory(f *Factory, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.Factory Error: at least one column to colsAndArgsFactory")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case FactoryCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = f.Id
		case FactoryCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = f.Created
		case FactoryCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = f.Updated
		case FactoryCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = f.Visibly
		case FactoryCol.Sid:
			cols[ii] = "`sid` = ?"
			args[ii] = f.Sid
		case FactoryCol.Name:
			cols[ii] = "`name` = ?"
			args[ii] = f.Name
		case FactoryCol.Location:
			cols[ii] = "`location` = ?"
			args[ii] = f.Location
		case FactoryCol.Token:
			cols[ii] = "`token` = ?"
			args[ii] = f.Token
		case FactoryCol.IsStaff:
			cols[ii] = "`is_staff` = ?"
			args[ii] = f.IsStaff
		default:
			return nil, nil, logError(fmt.Sprintf("dal.Factory Error: unknow column num %d in talbe factory", cc))
		}
	}
	return cols, args, nil
}

var FactoryEnum = struct {
}{}

var factorySqls = []string{
	/*
		CREATE TABLE `factory` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `sid` varchar(64) NOT NULL,
		  `name` varchar(64) NOT NULL,
		  `location` varchar(64) NOT NULL,
		  `token` varchar(64) NOT NULL,
		  `is_staff` bool NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into factory(`created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff`) values(now(), now(), 1, ?, ?, ?, ?, ?)",
	/*1*/ "insert into factory(`id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff`) values(?, now(), now(), 1, ?, ?, ?, ?, ?)",
	/*2*/ "update factory set updated = now(), `visibly` = ?, `sid` = ?, `name` = ?, `location` = ?, `token` = ?, `is_staff` = ? where id = ?",
	/*3*/ "update factory set updated = now(), %s where id = ?",
	/*4*/ "update factory set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from factory where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where id in (%s) and visibly = 1",

	/*8*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where visibly = 1 and sid = ? limit 0, 1",
	/*9*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where visibly = 1 order by id desc limit ?, ?",
	/*10*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where visibly = 1 and token = ? limit 0, 1",
	/*11*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where visibly = 1 and name = ? limit 0, 1",
	/*12*/ "select `id`, `created`, `updated`, `visibly`, `sid`, `name`, `location`, `token`, `is_staff` from factory where visibly = 1 and is_staff = 1 limit 0, 1",
}

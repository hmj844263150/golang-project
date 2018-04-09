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

var ModuleTid = 1
var _ db.Doer = (*Module)(nil)
var modulecols = []db.Col{1, 2, 3, 4, 8}
var modulefields = []string{"id", "created", "updated", "visibly", "esp_mac"}

var ModuleCol = struct {
	Id, Created, Updated, Visibly, EspMac, _ db.Col
}{1, 2, 3, 4, 8, 0}

type Module struct {
	Id      int
	Created time.Time
	Updated time.Time
	Visibly bool
	EspMac  string

	// ext, not persistent field
	ext      *Ext
	paddings map[string]interface{}
}

func NewModule(ctx context.Context) *Module {
	now := time.Now()
	m := &Module{Created: now, Updated: now, Visibly: true}
	m.ext = GetExtFromContext(ctx)
	defaultModule(ctx, m)
	return m
}

func FindModule(ctx context.Context, id int) *Module {
	dos, err := db.Open("Module").Query(newModuleDest, true, moduleSqls[6], id)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	for _, do := range dos {
		module, _ := do.(*Module)
		if ext != nil {
			module.ext = ext
		}
		return module
	}
	return nil
}

func ListModule(ctx context.Context, ids ...int) []*Module {
	holders := make([]string, len(ids))
	generic := make([]interface{}, len(ids))
	for ii, id := range ids {
		holders[ii] = "?"
		generic[ii] = id
	}
	sql := fmt.Sprintf(moduleSqls[7], strings.Join(holders, ", "))
	dos, err := db.Open("Module").Query(newModuleDest, true, sql, generic...)
	if err != nil {
		return nil
	}
	ext := GetExtFromContext(ctx)
	modules := make([]*Module, len(dos))
	for ii, do := range dos {
		module, _ := do.(*Module)
		if ext != nil {
			module.ext = ext
		}
		modules[ii] = module
	}
	return modules
}

func (m *Module) Save() error {
	now := time.Now()
	m.Created, m.Updated, m.Visibly = now, now, true
	var id int64
	var err error
	if m.Id == 0 {
		id, _, err = db.Open("Module").Exec(moduleSqls[0], m.EspMac)
	} else {
		id, _, err = db.Open("Module").Exec(moduleSqls[1], m.Id, m.EspMac)
	}
	if err != nil {
		return err
	}
	m.Id = int(id)
	return nil
}

func (m *Module) Update(cs ...db.Col) error {
	if m.Id == 0 {
		return logError("dal.Module Error: can not update row while id is zero")
	}
	m.Updated = time.Now()
	if len(cs) == 0 {
		_, _, err := db.Open("Module").Exec(moduleSqls[2], m.Visibly, m.EspMac, m.Id)
		return err
	}
	cols, args, err := colsAndArgsModule(m, cs...)
	if err != nil {
		return err
	}
	args = append(args, m.Id)
	sqlstr := fmt.Sprintf(moduleSqls[3], strings.Join(cols, ", "))
	_, _, err = db.Open("Module").Exec(sqlstr, args...)
	return err
}

func (m *Module) Invisibly() error {
	if m.Id == 0 {
		return logError("dal.Module Error: can not invisibly row while id is zero")
	}
	m.Updated = time.Now()
	m.Visibly = false
	_, _, err := db.Open("Module").Exec(moduleSqls[4], m.Id)
	return err
}

func (m *Module) Delete() error {
	if m.Id == 0 {
		return logError("dal.Module Error: can not delete row while id is zero")
	}
	m.Updated = time.Now()
	_, _, err := db.Open("Module").Exec(moduleSqls[5], m.Id)
	return err
}

func (m *Module) Valid() error {
	return m.valid()
}

func (m *Module) SetExt(ext *Ext) {
	m.ext = ext
}

func (m *Module) Padding(pkey string, pvalue interface{}) {
	if m.ext == nil {
		m.ext = &Ext{Loc: DefaultLoc}
	}
	if m.paddings == nil {
		m.paddings = make(map[string]interface{})
	}
	m.paddings[pkey] = pvalue
	m.ext.IsComplex = true
}

func (m *Module) AsMap(isColumnName bool, cs ...db.Col) map[string]interface{} {
	mm := make(map[string]interface{})
	for _, cc := range cs {
		switch cc {
		case ModuleCol.Id:
			if isColumnName {
				mm["id"] = m.Id
			} else {
				mm["Id"] = m.Id
			}
		case ModuleCol.Created:
			if isColumnName {
				mm["created"] = m.Created
			} else {
				mm["Created"] = m.Created
			}
		case ModuleCol.Updated:
			if isColumnName {
				mm["updated"] = m.Updated
			} else {
				mm["Updated"] = m.Updated
			}
		case ModuleCol.Visibly:
			if isColumnName {
				mm["visibly"] = m.Visibly
			} else {
				mm["Visibly"] = m.Visibly
			}
		case ModuleCol.EspMac:
			if isColumnName {
				mm["esp_mac"] = m.EspMac
			} else {
				mm["EspMac"] = m.EspMac
			}
		default:
			logError(fmt.Sprintf("dal.Module Error: unknow column num %d in talbe module", cc))
		}
	}
	return mm
}

func (m *Module) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	loc := DefaultLoc
	var numericEnum bool
	if m.ext != nil {
		if m.ext.IsComplex {
			return m.marshalJSONComplex()
		}
		loc = m.ext.Loc
		numericEnum = m.ext.NumericEnum
	}
	var buf bytes.Buffer
	buf.WriteString(`{"id":`)
	buf.WriteString(strconv.FormatInt(int64(m.Id), 10))
	buf.WriteString(`, "created":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(m.Created.Unix(), 10))
	} else {
		m.Created = m.Created.In(loc)
		buf.WriteString(`"` + m.Created.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "updated":`)
	if loc == Epoch {
		buf.WriteString(strconv.FormatInt(m.Updated.Unix(), 10))
	} else {
		m.Updated = m.Updated.In(loc)
		buf.WriteString(`"` + m.Updated.Format(DefaultTimeFormat) + `"`)
	}
	buf.WriteString(`, "visibly":`)
	if numericEnum {
		if m.Visibly {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	} else {
		buf.WriteString(strconv.FormatBool(m.Visibly))
	}
	buf.WriteString(`, "esp_mac":`)
	WriteJsonString(&buf, m.EspMac)
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (m *Module) marshalJSONComplex() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	if m.ext == nil {
		return nil, logError("dal.Module Error: can not marshalJSONComplex with .ext == nil")
	}
	loc := m.ext.Loc
	numericEnum := m.ext.NumericEnum
	var includes, excludes map[db.Col]interface{}
	if vv, ok := dalVerboses[ModuleTid]; ok {
		if vvv, ok := vv[m.ext.Verbose]; ok {
			includes, excludes = vvv[0], vvv[1]
		}
	}
	paddings := m.paddings
	var buf bytes.Buffer
	var isRender bool
	isRender = isRenderField(ModuleCol.Id, "id", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "id":`)
		buf.WriteString(strconv.FormatInt(int64(m.Id), 10))
	}
	isRender = isRenderField(ModuleCol.Created, "created", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "created":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(m.Created.Unix(), 10))
		} else {
			m.Created = m.Created.In(loc)
			buf.WriteString(`"` + m.Created.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(ModuleCol.Updated, "updated", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "updated":`)
		if loc == Epoch {
			buf.WriteString(strconv.FormatInt(m.Updated.Unix(), 10))
		} else {
			m.Updated = m.Updated.In(loc)
			buf.WriteString(`"` + m.Updated.Format(DefaultTimeFormat) + `"`)
		}
	}
	isRender = isRenderField(ModuleCol.Visibly, "visibly", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "visibly":`)
		if numericEnum {
			if m.Visibly {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		} else {
			buf.WriteString(strconv.FormatBool(m.Visibly))
		}
	}
	isRender = isRenderField(ModuleCol.EspMac, "esp_mac", includes, excludes, paddings)
	if isRender {
		buf.WriteString(`, "esp_mac":`)
		WriteJsonString(&buf, m.EspMac)
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

func (m *Module) UnmarshalMap(ctx context.Context, vi interface{}, cols ...db.Col) ([]db.Col, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalModule with null value")
	}
	vv, ok := vi.(map[string]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalModule")
	}
	updatedCols := []db.Col{}
	if len(cols) == 0 {
		cols = modulecols
	}
	loc := DefaultLoc
	for _, col := range cols {
		switch col {
		case ModuleCol.Id:
			vvv, ok := vv["id"]
			if !ok {
				continue
			}
			m.Id, err = Int(vvv)
			updatedCols = append(updatedCols, col)
		case ModuleCol.Created:
			vvv, ok := vv["created"]
			if !ok {
				continue
			}
			m.Created, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case ModuleCol.Updated:
			vvv, ok := vv["updated"]
			if !ok {
				continue
			}
			m.Updated, err = Time(vvv, loc)
			updatedCols = append(updatedCols, col)
		case ModuleCol.Visibly:
			vvv, ok := vv["visibly"]
			if !ok {
				continue
			}
			m.Visibly, err = Bool(vvv)
			updatedCols = append(updatedCols, col)
		case ModuleCol.EspMac:
			vvv, ok := vv["esp_mac"]
			if !ok {
				continue
			}
			m.EspMac, err = String(vvv)
			updatedCols = append(updatedCols, col)
		}
		if err != nil {
			return nil, err
		}
	}
	return cols, nil
}

func UnmarshalModule(ctx context.Context, vi interface{}, cols ...db.Col) (*Module, error) {
	m := NewModule(ctx)
	_, err := m.UnmarshalMap(ctx, vi, cols...)
	if err != nil {
		return nil, err
	}
	return m, err
}

func UnmarshalModules(ctx context.Context, vi interface{}, cols ...db.Col) ([]*Module, error) {
	var err error
	if vi == nil {
		return nil, logError("can not UnmarshalModules with null value")
	}
	vv, ok := vi.([]interface{})
	if !ok {
		return nil, logError("type asserted failed when UnmarshalModules")
	}
	modules := make([]*Module, len(vv))
	for ii, vvv := range vv {
		var m *Module
		m, err = UnmarshalModule(ctx, vvv, cols...)
		if err != nil {
			return nil, err
		}
		modules[ii] = m
	}
	return modules, nil
}

func newModuleDest(cols ...string) (db.Doer, []interface{}, error) {
	m := &Module{}
	if cols == nil || len(cols) == 0 {
		return m, []interface{}{&m.Id, &m.Created, &m.Updated, &m.Visibly, &m.EspMac}, nil
	}
	dest := make([]interface{}, len(cols))
	for ii, col := range cols {
		switch col {
		case "id":
			dest[ii] = &m.Id
		case "created":
			dest[ii] = &m.Created
		case "updated":
			dest[ii] = &m.Updated
		case "visibly":
			dest[ii] = &m.Visibly
		case "esp_mac":
			dest[ii] = &m.EspMac
		default:
			return nil, nil, logError("dal.Module Error: unknow column " + col + " in talbe module")
		}
	}
	return m, dest, nil
}

func colsAndArgsModule(m *Module, cs ...db.Col) ([]string, []interface{}, error) {
	len := len(cs)
	if len == 0 {
		return nil, nil, logError("dal.Module Error: at least one column to colsAndArgsModule")
	}
	cols := make([]string, len)
	args := make([]interface{}, len)
	for ii, cc := range cs {
		switch cc {
		case ModuleCol.Id:
			cols[ii] = "`id` = ?"
			args[ii] = m.Id
		case ModuleCol.Created:
			cols[ii] = "`created` = ?"
			args[ii] = m.Created
		case ModuleCol.Updated:
			cols[ii] = "`updated` = ?"
			args[ii] = m.Updated
		case ModuleCol.Visibly:
			cols[ii] = "`visibly` = ?"
			args[ii] = m.Visibly
		case ModuleCol.EspMac:
			cols[ii] = "`esp_mac` = ?"
			args[ii] = m.EspMac
		default:
			return nil, nil, logError(fmt.Sprintf("dal.Module Error: unknow column num %d in talbe module", cc))
		}
	}
	return cols, args, nil
}

var ModuleEnum = struct {
}{}

var moduleSqls = []string{
	/*
		CREATE TABLE `module` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `created` datetime NOT NULL,
		  `updated` datetime NOT NULL,
		  `visibly` bool NOT NULL,
		  `esp_mac` varchar(64) NOT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	*/
	/*0*/ "insert into module(`created`, `updated`, `visibly`, `esp_mac`) values(now(), now(), 1, ?)",
	/*1*/ "insert into module(`id`, `created`, `updated`, `visibly`, `esp_mac`) values(?, now(), now(), 1, ?)",
	/*2*/ "update module set updated = now(), `visibly` = ?, `esp_mac` = ? where id = ?",
	/*3*/ "update module set updated = now(), %s where id = ?",
	/*4*/ "update module set visibly = 0, updated = now() where id = ?",
	/*5*/ "delete from module where id = ?",
	/*6*/ "select `id`, `created`, `updated`, `visibly`, `esp_mac` from module where id = ? and visibly = 1",
	/*7*/ "select `id`, `created`, `updated`, `visibly`, `esp_mac` from module where id in (%s) and visibly = 1",
}

package rpc

import (
	"errors"
	"espressif.com/chip/factory/dal"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Ctx     context.Context
	Nonce   int
	Version Version
	Method  Method
	Path    string
	Token   string
	Tzname  string
	Loc     *time.Location
	Mode    Mode

	Proto  Proto
	Format Format

	Action   string
	Start    time.Time
	End      time.Time
	Offset   int
	RowCount int

	Args map[string]interface{}
	Meta map[string]interface{}
	Get  map[string]interface{}
	Body map[string]interface{}

	Factory *dal.Factory

	R *http.Request
}

func NewRequest() *Request {
	r := &Request{Tzname: "Asia/Shanghai", Loc: UTC8, Mode: None}
	return r
}

func (r *Request) Auth(requireStaff bool) {
	factory := dal.FindFactoryByToken(r.Ctx, r.Token)
	if factory == nil {
		panic("403 factory not exists")
	}
	if requireStaff && !factory.IsStaff {
		panic("403 require staff")
	}
	r.Factory = factory
}

func (r *Request) SetMethod(method string) error {
	var ok bool
	r.Method, ok = methodMap[method]
	if !ok {
		r.Method = Get
	}
	return nil
}

func (r *Request) SetMode(mode string) error {
	var ok bool
	r.Mode, ok = modeMap[mode]
	if !ok {
		r.Mode = None
	}
	return nil
}

func (r *Request) HasAction() bool {
	return r.Action != ""
}

func (r *Request) IsParamTrue(param string) bool {
	return false
}

func (r *Request) GetString(arg string, m map[string]interface{}) (string, error) {
	var ok bool
	if m == nil {
		return "", nilMap
	}
	v, ok := m[arg]
	if !ok {
		return "", errors.New("cat not GetString")
	}
	vv, ok := v.(string)
	if !ok {
		return "", errors.New("cat not GetString")
	}
	return vv, nil
}

func (r *Request) GetInt(arg string, m map[string]interface{}) (int, error) {
	if m == nil {
		return 0, nilMap
	}
	v, ok := m[arg]
	if !ok {
		return 0, errors.New("cat not GetInt")
	}
	vv, ok := v.(float64)
	if ok {
		return int(vv), nil
	}
	vvv, ok := v.(int)
	if ok {
		return vvv, nil
	}
	vvvv, ok := v.(string)
	if ok {
		i, err := strconv.ParseInt(vvvv, 10, 0)
		if err == nil {
			return int(i), err
		}
	}
	return 0, errors.New("cat not GetInt")
}

func (r *Request) GetFloat64(arg string, m map[string]interface{}) (float64, error) {
	if m == nil {
		return 0, nilMap
	}
	v, ok := m[arg]
	if !ok {
		return 0, errors.New("cat not GetFloat64")
	}
	vv, ok := v.(float64)
	if !ok {
		return 0, errors.New("cat not GetFloat64")
	}
	return vv, nil
}

func (r *Request) GetTime(arg string, m map[string]interface{}) (time.Time, error) {
	if m == nil {
		return time.Now(), nilMap
	}
	v, err := r.GetString(arg, m)
	if err != nil {
		return time.Now(), err
	}
	t, err := dal.Time(v, r.Loc)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}

func (r *Request) GetBool(arg string, m map[string]interface{}) (bool, error) {
	if m == nil {
		return false, nilMap
	}
	v, ok := m[arg]
	if !ok {
		return false, errors.New("cat not GetBool")
	}
	vv, ok := v.(bool)
	if ok {
		return vv, nil
	}
	vvv, ok := v.(string)
	if ok {
		return vvv == "true", nil
	}
	return false, errors.New("cat not GetBool")
}

func (r *Request) GetMap(arg string, m map[string]interface{}) (map[string]interface{}, error) {
	if m == nil {
		return nil, nilMap
	}
	v, ok := m[arg]
	if !ok {
		return nil, errors.New("cat not GetMap")
	}
	vv, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.New("cat not GetMap")
	}
	return vv, nil
}

func (r *Request) GetSlice(arg string, m map[string]interface{}) ([]interface{}, error) {
	if m == nil {
		return nil, nilMap
	}
	v, ok := m[arg]
	if !ok {
		return nil, errors.New("cat not GetSlice")
	}
	vv, ok := v.([]interface{})
	if !ok {
		return nil, errors.New("cat not GetSlice")
	}
	return vv, nil
}

func (r *Request) GetSliceMap(arg string, m map[string]interface{}) ([]map[string]interface{}, error) {
	if m == nil {
		return nil, nilMap
	}
	v, ok := m[arg]
	if !ok {
		return nil, errors.New("cat not GetSliceMap")
	}
	vv, ok := v.([]interface{})
	if !ok {
		return nil, errors.New("cat not GetSliceMap")
	}
	vvv := []map[string]interface{}{}
	for _, x := range vv {
		m, ok := x.(map[string]interface{})
		if !ok {
			continue
		}
		vvv = append(vvv, m)
	}
	return vvv, nil
}

func (r *Request) GetValue(arg string, m map[string]interface{}) interface{} {
	if m == nil {
		return nil
	}

	var v interface{}
	var ok bool

	s := strings.Split(arg, ".")
	for i, k := range s {
		v, ok = m[k]
		if !ok {
			return nil
		}
		if i == len(s)-1 {
			return v
		}
		switch vv := v.(type) {
		case map[string]interface{}:
			m = vv
		default:
			return nil
		}
	}
	return nil
}

func (r *Request) FillRange(limitRange bool, timeRange bool) {
	if limitRange {
		r.Offset, r.RowCount = 0, 500
		offset, err := r.GetInt("offset", r.Get)
		if err == nil {
			r.Offset = offset
		}
		rowCount, err := r.GetInt("row_count", r.Get)
		if err == nil {
			r.RowCount = rowCount
		}
		if r.RowCount > 1000 {
			r.RowCount = 1000
		}
	}
	if timeRange {
		r.Start, r.End = dal.MinTime, dal.MaxTime
		start, err := r.GetTime("start", r.Get)
		if err == nil {
			r.Start = start
		}
		end, err := r.GetTime("end", r.Get)
		if err == nil {
			r.End = end
		}
		r.Start = r.Start.In(UTC8)
		r.End = r.End.In(UTC8)
	}
}

func (r *Request) Fill(args []string, params ...interface{}) error {
	if len(args) != len(params) {
		return errors.New("args, params length not match")
	}

	var err error
	var ppp int64
	for i, p := range params {
		arg := args[i]
		value := r.GetValue(arg, r.Body)
		if value == nil {
			continue
		}
		switch pp := p.(type) {
		case *string:
			str, ok := value.(string)
			if ok {
				*pp = str
			}
		case *int:
			*pp, err = dal.Int(value)
		case *time.Time:
			*pp, err = dal.Time(value, r.Loc)
		case *map[string]interface{}:
			m, ok := value.(map[string]interface{})
			if ok {
				*pp = m
			}
		case *float64:
			*pp, err = dal.Float64(value)
		case *float32:
			*pp, err = dal.Float32(value)
		case *int64:
			*pp, err = dal.Int64(value)
		case *int32:
			*pp, err = dal.Int32(value)
		case *int16:
			ppp, err = dal.Int64(value)
			*pp = int16(ppp)
		case *int8:
			ppp, err = dal.Int64(value)
			*pp = int8(ppp)
		case *uint64:
			ppp, err = dal.Int64(value)
			*pp = uint64(ppp)
		case *uint32:
			ppp, err = dal.Int64(value)
			*pp = uint32(ppp)
		case *uint16:
			ppp, err = dal.Int64(value)
			*pp = uint16(ppp)
		case *uint8:
			ppp, err = dal.Int64(value)
			*pp = uint8(ppp)
		}
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

package dal

import (
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"espressif.com/chip/factory/db"
	"golang.org/x/net/context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var MinTime, _ = time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
var MaxTime, _ = time.Parse("2006-01-02 15:04:05", "9999-01-01 00:00:00")

type Pader interface {
	Padding(pkey string, pvalue interface{})
}

var dalVerboses = map[int]map[string][]map[db.Col]interface{}{}

var hexstr = "0123456789abcdef"
var random = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var Epoch = &time.Location{}
var DefaultLoc, _ = time.LoadLocation("Asia/Shanghai")
var DefaultTimeFormat = "2006-01-02 15:04:05"

func Randstr(l int) string {
	half := l / 2
	odd := (l%2 != 0)
	if odd {
		half = half + 1
	}
	b := make([]byte, half)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	s := hex.EncodeToString(b)
	if odd {
		return s[0 : len(s)-1]
	}
	return s
}

type Ext struct {
	Loc         *time.Location
	IsComplex   bool
	NumericEnum bool
	Verbose     string
	Trace       []byte
}

func (e *Ext) Copy() *Ext {
	return &Ext{Loc: e.Loc, Verbose: e.Verbose, Trace: e.Trace}
}

func GetExtFromContext(ctx context.Context) *Ext {
	if ctx == nil {
		return nil
	}
	v := ctx.Value("ext")
	if v == nil {
		return nil
	}
	if vv, ok := v.(*Ext); ok {
		return vv
	}
	return nil
}

func Enum(v interface{}, p map[string]int, r map[int]string, m string) (int, error) {
	str, ok := v.(string)
	if ok {
		if num, ok := p[str]; ok {
			return num, nil
		}
	}
	num, err := Int(v)
	if err == nil {
		if _, ok := r[num]; ok {
			return num, nil
		}
	}
	return 0, logError("can not conv to enum, only support: " + m)
}

func Time(v interface{}, loc *time.Location) (time.Time, error) {
	var t time.Time
	var err error
	tt, ok := v.(time.Time)
	if ok {
		return tt, nil
	}
	str, ok := v.(string)
	if ok {
		switch len(str) {
		case 25, 20:
			t, err = time.ParseInLocation(time.RFC3339, str, loc)
		case 19:
			t, err = time.ParseInLocation("2006-01-02 15:04:05", str, loc)
		default:
			return t, errors.New("can not conv to time, unsupported string format")
		}
		return t, err
	}
	num, err := Int64(v)
	if err == nil {
		t = time.Unix(num, 0).In(loc)
		return t, nil
	}
	return t, errors.New("conv to time failed")
}

func Bool(v interface{}) (bool, error) {
	b, ok := v.(bool)
	if ok {
		return b, nil
	}
	str, ok := v.(string)
	if ok {
		return str == "true", nil
	}
	num, err := Int64(v)
	if err == nil {
		return num == int64(1), nil
	}
	return false, errors.New("conv to bool failed")
}

func Int(v interface{}) (int, error) {
	switch vv := v.(type) {
	case float64:
		return int(vv), nil
	case float32:
		return int(vv), nil
	case string:
		vvv, err := strconv.ParseInt(vv, 10, 0)
		if err != nil {
			return 0, err
		}
		return int(vvv), nil
	case int:
		return int(vv), nil
	case int64:
		return int(vv), nil
	case int32:
		return int(vv), nil
	case int16:
		return int(vv), nil
	case int8:
		return int(vv), nil
	case uint:
		return int(vv), nil
	case uint64:
		return int(vv), nil
	case uint32:
		return int(vv), nil
	case uint16:
		return int(vv), nil
	case uint8:
		return int(vv), nil
	}
	return 0, logError("can not conv to int")
}

func Int64(v interface{}) (int64, error) {
	switch vv := v.(type) {
	case float64:
		return int64(vv), nil
	case float32:
		return int64(vv), nil
	case int:
		return int64(vv), nil
	case int64:
		return int64(vv), nil
	case int32:
		return int64(vv), nil
	case int16:
		return int64(vv), nil
	case int8:
		return int64(vv), nil
	case uint:
		return int64(vv), nil
	case uint64:
		return int64(vv), nil
	case uint32:
		return int64(vv), nil
	case uint16:
		return int64(vv), nil
	case uint8:
		return int64(vv), nil
	}
	return 0, logError("can not conv to int64")
}

func Int32(v interface{}) (int32, error) {
	switch vv := v.(type) {
	case float64:
		return int32(vv), nil
	case float32:
		return int32(vv), nil
	case int:
		return int32(vv), nil
	case int64:
		return int32(vv), nil
	case int32:
		return int32(vv), nil
	case int16:
		return int32(vv), nil
	case int8:
		return int32(vv), nil
	case uint64:
		return int32(vv), nil
	case uint32:
		return int32(vv), nil
	case uint16:
		return int32(vv), nil
	case uint8:
		return int32(vv), nil
	}
	return 0, logError("can not conv to int32")
}

func Int16(v interface{}) (int16, error) {
	switch vv := v.(type) {
	case float64:
		return int16(vv), nil
	case float32:
		return int16(vv), nil
	case int:
		return int16(vv), nil
	case int64:
		return int16(vv), nil
	case int32:
		return int16(vv), nil
	case int16:
		return int16(vv), nil
	case int8:
		return int16(vv), nil
	case uint64:
		return int16(vv), nil
	case uint32:
		return int16(vv), nil
	case uint16:
		return int16(vv), nil
	case uint8:
		return int16(vv), nil
	}
	return 0, logError("can not conv to int16")
}

func Int8(v interface{}) (int8, error) {
	switch vv := v.(type) {
	case float64:
		return int8(vv), nil
	case float32:
		return int8(vv), nil
	case int:
		return int8(vv), nil
	case int64:
		return int8(vv), nil
	case int32:
		return int8(vv), nil
	case int16:
		return int8(vv), nil
	case int8:
		return int8(vv), nil
	case uint64:
		return int8(vv), nil
	case uint32:
		return int8(vv), nil
	case uint16:
		return int8(vv), nil
	case uint8:
		return int8(vv), nil
	}
	return 0, logError("can not conv to int8")
}

func Uint64(v interface{}) (uint64, error) {
	switch vv := v.(type) {
	case float64:
		return uint64(vv), nil
	case float32:
		return uint64(vv), nil
	case int:
		return uint64(vv), nil
	case int64:
		return uint64(vv), nil
	case int32:
		return uint64(vv), nil
	case int16:
		return uint64(vv), nil
	case int8:
		return uint64(vv), nil
	case uint:
		return uint64(vv), nil
	case uint64:
		return uint64(vv), nil
	case uint32:
		return uint64(vv), nil
	case uint16:
		return uint64(vv), nil
	case uint8:
		return uint64(vv), nil
	}
	return 0, logError("can not conv to uint64")
}

func Uint32(v interface{}) (uint32, error) {
	switch vv := v.(type) {
	case float64:
		return uint32(vv), nil
	case float32:
		return uint32(vv), nil
	case int:
		return uint32(vv), nil
	case int64:
		return uint32(vv), nil
	case int32:
		return uint32(vv), nil
	case int16:
		return uint32(vv), nil
	case int8:
		return uint32(vv), nil
	case uint64:
		return uint32(vv), nil
	case uint32:
		return uint32(vv), nil
	case uint16:
		return uint32(vv), nil
	case uint8:
		return uint32(vv), nil
	}
	return 0, logError("can not conv to uint32")
}

func Uint16(v interface{}) (uint16, error) {
	switch vv := v.(type) {
	case float64:
		return uint16(vv), nil
	case float32:
		return uint16(vv), nil
	case int:
		return uint16(vv), nil
	case int64:
		return uint16(vv), nil
	case int32:
		return uint16(vv), nil
	case int16:
		return uint16(vv), nil
	case int8:
		return uint16(vv), nil
	case uint64:
		return uint16(vv), nil
	case uint32:
		return uint16(vv), nil
	case uint16:
		return uint16(vv), nil
	case uint8:
		return uint16(vv), nil
	}
	return 0, logError("can not conv to uint16")
}

func Uint8(v interface{}) (uint8, error) {
	switch vv := v.(type) {
	case float64:
		return uint8(vv), nil
	case float32:
		return uint8(vv), nil
	case int:
		return uint8(vv), nil
	case int64:
		return uint8(vv), nil
	case int32:
		return uint8(vv), nil
	case int16:
		return uint8(vv), nil
	case int8:
		return uint8(vv), nil
	case uint64:
		return uint8(vv), nil
	case uint32:
		return uint8(vv), nil
	case uint16:
		return uint8(vv), nil
	case uint8:
		return uint8(vv), nil
	}
	return 0, logError("can not conv to uint8")
}

func Float64(v interface{}) (float64, error) {
	switch vv := v.(type) {
	case float64:
		return float64(vv), nil
	case float32:
		return float64(vv), nil
	}
	return float64(0), logError("can not conv to float64")
}

func Float32(v interface{}) (float32, error) {
	switch vv := v.(type) {
	case float64:
		return float32(vv), nil
	case float32:
		return float32(vv), nil
	}
	return float32(0), logError("can not conv to float32")
}

func String(v interface{}) (string, error) {
	str, ok := v.(string)
	if ok {
		return str, nil
	}
	return "", logError("can not conv to string")
}

func Join(v interface{}, sep string) string {
	switch vv := v.(type) {
	case []int:
		xs := make([]string, len(vv))
		for i, x := range vv {
			xs[i] = strconv.FormatInt(int64(x), 10)
		}
		return strings.Join(xs, sep)
	case []string:
		return strings.Join(vv, sep)
	}
	return ""
}

func General(v interface{}) (string, error) {
	switch vv := v.(type) {
	case string:
		return `"` + vv + `"`, nil
	case float64:
		return strconv.FormatFloat(vv, 'f', 2, 64), nil
	case float32:
		return strconv.FormatFloat(float64(vv), 'f', 2, 32), nil
	case int:
		return strconv.FormatInt(int64(vv), 10), nil
	case int64:
		return strconv.FormatInt(vv, 10), nil
	case int32:
		return strconv.FormatInt(int64(vv), 10), nil
	case int16:
		return strconv.FormatInt(int64(vv), 10), nil
	case int8:
		return strconv.FormatInt(int64(vv), 10), nil
	case uint:
		return strconv.FormatInt(int64(vv), 10), nil
	case uint64:
		return strconv.FormatInt(int64(vv), 10), nil
	case uint32:
		return strconv.FormatInt(int64(vv), 10), nil
	case uint16:
		return strconv.FormatInt(int64(vv), 10), nil
	case uint8:
		return strconv.FormatInt(int64(vv), 10), nil
	case []byte:
		return `"` + string(vv) + `"`, nil
	case json.Marshaler:
		bs, err := vv.MarshalJSON()
		if err != nil {
			logError(err.Error())
			return `""`, err
		}
		return string(bs), nil
	default:
		bs, err := json.Marshal(v)
		if err != nil {
			logError(err.Error())
			return `""`, err
		}
		return string(bs), nil
	}
	return "", logError("can not conv to general string")
}

func logError(errstr string) error {
	err := errors.New(errstr)
	log.Println(err)
	return err
}

func randEnumStr(rr map[int]string) string {
	r := randInt16()
	if r < 0 {
		r = -r
	}
	s := make([]string, len(rr))
	i := 0
	for _, v := range rr {
		s[i] = v
		i = i + 1
	}
	return s[r%(len(s))]
}

func randEnumInt(rr map[int]string) int {
	r := randInt16()
	if r < 0 {
		r = -r
	}
	s := make([]int, len(rr))
	i := 0
	for k, _ := range rr {
		s[i] = k
		i = i + 1
	}
	return s[r%(len(s))]
}

func randInt16() int {
	return int(int16(rand.Int()))
}

func randInt32() int {
	return int(int32(rand.Int()))
}

func randInt() int {
	return rand.Int()
}

func randFloat64() float64 {
	return rand.Float64()
}

func randFloat32() float32 {
	return rand.Float32()
}

func randStr(n int) string {
	n = (randInt() % n)
	if n < 0 {
		n = -n
	}
	if n < 7 {
		n = 7
	}
	l := len(random)
	bs := make([]byte, n)
	for i := 0; i < n; i++ {
		bs[i] = random[(randInt() % l)]
	}
	return string(bs)
}

func almostSameTime(t1 time.Time, t2 time.Time, seconds int) bool {
	return int(math.Abs(float64(t1.Unix()-t2.Unix()))) < seconds
}

func almostSameFloat(f1 float64, f2 float64, f float64) bool {
	return math.Abs(f1-f2) < f
}

func randTime() time.Time {
	num := 1460359340 + (randInt() % 100000)
	return time.Unix(int64(num), 0)
}

func randBool() bool {
	num := randInt()
	return num%2 == 0
}

func isRenderField(col db.Col, columnName string, includes map[db.Col]interface{}, excludes map[db.Col]interface{}, paddings map[string]interface{}) bool {
	var ok bool
	if len(includes) > 0 {
		if _, ok = includes[col]; !ok {
			return false
		}
	}
	if len(excludes) > 0 {
		if _, ok = excludes[col]; ok {
			return false
		}
	}
	if len(paddings) > 0 {
		if _, ok = paddings[columnName]; ok {
			return false
		}
	}
	return true
}

func WriteJsonString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as <, > and &. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hexstr[b>>4])
				buf.WriteByte(hexstr[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hexstr[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
}

func init() {
	rand.Seed(time.Now().Unix())
}

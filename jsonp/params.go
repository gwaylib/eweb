package jsonp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gwaylib/errors"
	"github.com/shopspring/decimal"
)

var (
	UNIX_TIME_NO_SET = time.Time{}.Unix()
)

type Params map[string]interface{}

func ParseParams(data []byte) (Params, error) {
	params := Params{}
	if err := json.Unmarshal(data, &params); err != nil {
		return params, errors.As(err, string(data))
	}
	return params, nil
}

func ParseParamsByIO(r io.Reader) (Params, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.As(err)
	}
	return ParseParams(body)
}

func (p Params) JsonData() []byte {
	data, _ := json.Marshal(p)
	return data
}

// Obsoleted, call Set
func (p Params) Add(key, value string) {
	p.Set(key, value)
}

// Obsoleted, call SetParams
func (p Params) AddParams(key string, param Params) {
	p.SetParams(key, param)
}

// Obsoleted, call SetAny
func (p Params) AddAny(key string, param interface{}) {
	p.SetAny(key, param)
}

func (p Params) Set(key, value string) {
	p[key] = value
}
func (p Params) SetParams(key string, param Params) {
	p[key] = param
}
func (p Params) SetAny(key string, param interface{}) {
	p[key] = param
}

func (p Params) HasKey(key string) bool {
	_, ok := p[key]
	return ok
}

func (p Params) TrimString(key string) string {
	return strings.TrimSpace(p.String(key))
}

func (p Params) String(key string) string {
	v, ok := p[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return fmt.Sprint(v)
}
func (p Params) Bool(key string) bool {
	v, ok := p[key]
	if !ok {
		return false
	}
	f, ok := v.(bool)
	if ok {
		return f
	}
	return false
}
func (p Params) Float64(key string, noDataRet, errRet float64) float64 {
	v, ok := p[key]
	if !ok {
		return noDataRet
	}
	switch f := v.(type) {
	case float32:
		return float64(f)
	case float64:
		return float64(f)
	case string:
		num, err := json.Number(f).Float64()
		if err != nil {
			return errRet
		}
		return num
	default:
		num, err := json.Number(fmt.Sprint(v)).Float64()
		if err != nil {
			return errRet
		}
		return num
	}
}

func (p Params) Int64(key string, noDataRet, errRet int64) int64 {
	v, ok := p[key]
	if !ok {
		return noDataRet
	}
	switch i := v.(type) {
	case int8:
		return int64(i)
	case int16:
		return int64(i)
	case int32:
		return int64(i)
	case int64:
		return int64(i)
	case string:
		val, err := json.Number(i).Int64()
		if err != nil {
			return errRet
		}
		return val
	}
	return int64(p.Float64(key, float64(noDataRet), float64(errRet)))
}
func (p Params) Time(key string, layoutOpt ...string) time.Time {
	layout := time.RFC3339Nano
	if len(layoutOpt) > 0 {
		layout = layoutOpt[0]
	}
	s, ok := p[key]
	if !ok {
		return time.Time{}
	}
	t, _ := time.Parse(layout, s.(string))
	//return t.In(time.FixedZone("UTC", 8*60*60))
	return t
}

func (p Params) Decimal(key string, noDataRet, errRet float64) decimal.Decimal {
	v, ok := p[key]
	if !ok {
		return decimal.NewFromFloat(noDataRet)
	}
	switch d := v.(type) {
	case int8:
		return decimal.NewFromInt(int64(d))
	case int16:
		return decimal.NewFromInt(int64(d))
	case int32:
		return decimal.NewFromInt(int64(d))
	case int64:
		return decimal.NewFromInt(d)
	case float32:
		return decimal.NewFromFloat(float64(d))
	case float64:
		return decimal.NewFromFloat(d)
	case string:
		val, err := decimal.NewFromString(d)
		if err != nil {
			return decimal.NewFromFloat(errRet)
		}
		return val
	default:
		val, err := decimal.NewFromString(fmt.Sprint(v))
		if err != nil {
			return decimal.NewFromFloat(errRet)
		}
		return val
	}
}

func (p Params) Email(key string) string {
	email := p.String(key)
	if strings.Index(email, "@") < 1 {
		return key
	}
	for _, r := range email {
		if r > 255 {
			return key
		}
	}
	return email
}
func (p Params) Params(key string) Params {
	v, ok := p[key]
	if !ok {
		return Params{}
	}
	switch s := v.(type) {
	case map[string]interface{}:
		return Params(s)
	case Params:
		return Params(s)
	default:
		return Params{}
	}
}
func (p Params) Any(key string) interface{} {
	return p[key]
}

func (p Params) StringArray(key string) []string {
	s, ok := p[key]
	if !ok {
		return []string{}
	}
	arr, ok := s.([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, len(arr))
	for i, a := range arr {
		result[i] = fmt.Sprint(a)
	}
	return result
}
func (p Params) ParamsArray(key string) []Params {
	s, ok := p[key]
	if !ok {
		return []Params{}
	}
	arr, ok := s.([]interface{})
	if !ok {
		return []Params{}
	}
	result := []Params{}
	for _, a := range arr {
		p, ok := a.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, p)
	}
	return result
}
func (p Params) AnyArray(key string) []interface{} {
	s, ok := p[key]
	if !ok {
		return []interface{}{}
	}
	arr, ok := s.([]interface{})
	if !ok {
		return []interface{}{}
	}
	return arr
}

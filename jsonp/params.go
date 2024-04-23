package jsonp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gwaylib/errors"
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

func (p Params) TrimString(key string) string {
	return strings.TrimSpace(p.String(key))
}

func (p Params) String(key string) string {
	s, ok := p[key]
	if !ok {
		return ""
	}
	return fmt.Sprint(s)
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
	f, ok := v.(float64)
	if ok {
		return f
	}
	if len(fmt.Sprint(v)) == 0 {
		return noDataRet
	}
	out, err := strconv.ParseFloat(fmt.Sprint(v), 64)
	if err != nil {
		return errRet
	}
	return out
}

func (p Params) Int64(key string, noDataRet, errRet int64) int64 {
	v, ok := p[key]
	if !ok {
		return noDataRet
	}
	i, ok := v.(int64)
	if ok {
		return i
	}
	// fix big number
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

func (p Params) Email(key string) string {
	email := p.String(key)
	if strings.Index(email, "@") < 1 {
		return ""
	}
	for _, r := range email {
		if r > 255 {
			return ""
		}
	}
	return email
}
func (p Params) Params(key string) Params {
	s, ok := p[key]
	if !ok {
		return Params{}
	}
	sParams, ok := s.(map[string]interface{})
	if !ok {
		return Params{}
	}
	return Params(sParams)
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

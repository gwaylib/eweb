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
		return params, errors.As(err)
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

func (p Params) Add(key, value string) {
	p[key] = value
}
func (p Params) AddParams(key string, param Params) {
	p[key] = param
}
func (p Params) AddAny(key string, param interface{}) {
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
func (p Params) Int64(key string, noDataRet, errRet int64) int64 {
	s, ok := p[key]
	if !ok {
		return noDataRet
	}
	if len(fmt.Sprint(s)) == 0 {
		return noDataRet
	}
	i, err := strconv.ParseInt(fmt.Sprint(s), 10, 64)
	if err != nil {
		return errRet
	}
	return i
}
func (p Params) Float64(key string, noDataRet, errRet float64) float64 {
	s, ok := p[key]
	if !ok {
		return noDataRet
	}
	if len(fmt.Sprint(s)) == 0 {
		return noDataRet
	}
	i, err := strconv.ParseFloat(fmt.Sprint(s), 64)
	if err != nil {
		return errRet
	}
	return i
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
func (p Params) Any(key string) interface{} {
	return p[key]
}

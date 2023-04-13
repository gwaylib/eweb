package jsonp

import "testing"

func TestParams(t *testing.T) {
	in := `{
  "goods":[{
	"kind":"0",
	"num":"0"
  }]
}`
	p, err := ParseParams([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	goods := p.ParamsArray("goods")
	if len(goods) == 0 {
		t.Fatalf("expect 1, but 0")
	}
}

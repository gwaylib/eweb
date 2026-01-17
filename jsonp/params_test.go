package jsonp

import (
	"testing"
)

func TestParams(t *testing.T) {
	in := `{
  "goods":[{
	"kind":"0",
	"num":"0",
	"int":1024,
	"float":1.024,
	"bool":true,
	"obj":{"Balance":9979.48,"Credit":0,"Equity":9979.48,"Margin":0,"MarginFree":9979.48}
  }]
}`
	p, err := ParseParams([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	goods := p.ParamsArray("goods")
	if len(goods) == 0 {
		t.Fatal("expect 1, but 0")
	}
	if goods[0].Int64("int", 0, -1) != 1024 {
		t.Fatalf("expect 1024, but:%v", goods[0].Int64("int", 0, -1))
	}
	if goods[0].Float64("float", 0, -1) != 1.024 {
		t.Fatalf("expect 1.024, but:%v", goods[0].Float64("int", 0, -1))
	}
	if !goods[0].Bool("bool") {
		t.Fatal("expect true, but false")
	}
	if goods[0].Params("obj").Float64("Equity", 0, 0) != 9979.48 {
		t.Fatal("obj test failed")
	}
}

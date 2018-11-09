package main

import (
	"encoding/json"
	"fmt"
	"time"

	r "github.com/xy02/rdgo"
)

type m map[string]interface{}

func main() {
	f := r.Filter{TagKey: "json"}
	f.Input(&Msg{"ok", Detail{Age: 19}})
	str := `{"$f":{"Name":{"$eq":"ok"}, "detail.age":{"$gt":18,"$lt":31}}}`
	expr := r.Expr{}
	json.Unmarshal([]byte(str), &expr)
	// c := r.Condition{"Name": r.Is("ok"), "detail.age": r.Gt(18)}
	h1 := &r.Listener{OnData: func(data r.Data) {
		fmt.Printf("h1, %+v\n", data)
	}}
	h2 := &r.Listener{OnData: func(data r.Data) {
		fmt.Printf("h2, %+v\n", data)
	}}
	h3 := &r.Listener{OnData: func(data r.Data) {
		fmt.Printf("h3, %+v\n", data)
	}}
	f.Select(expr, h1)
	f.Select(expr, h2)
	f.Select(r.Expr{"$gt": 13}, h3)
	h1.Destroy()
	f.Input(Msg{"ok", Detail{Age: 22}})
	f.Input(Msg{"ok", Detail{Age: 18}})
	f.Input(m{"Name": "ok", "detail": m{"age": 30}})
	f.Input(12)
	f.Input(14)
	time.Sleep(time.Second)
}

type Msg struct {
	Name   string
	Detail Detail `json:"detail"`
}

type Detail struct {
	Name string
	Age  int `json:"age"`
}

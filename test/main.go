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
	str := `{"Name":{"$eq":"ok"}, "detail.age":{"$gt":18, "$lte":30}}`
	c := r.Condition{}
	json.Unmarshal([]byte(str), &c)
	// c := r.Condition{"Name": r.Is("ok"), "detail.age": r.Gt(18)}
	h1 := &r.Listener{OnData: func(data r.Data) {
		fmt.Printf("h1, %+v\n", data)
	}}
	h2 := &r.Listener{OnData: func(data r.Data) {
		fmt.Printf("h2, %+v\n", data)
	}}

	f.Select(c, h1)
	f.Select(c, h2)
	f.Select(c, h1)
	h1.Destroy()
	f.Input(&Msg{"ok", Detail{Age: 17}})
	f.Input(m{"Name": "ok", "detail": m{"age": 30}})
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

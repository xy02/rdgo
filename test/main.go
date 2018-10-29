package main

import (
	"fmt"
	"time"

	"github.com/xy02/rdgo"
)

func main() {
	f := rdgo.Filter{}
	f.Input(Msg{"ok"})
	c := rdgo.Condition{"Name": rdgo.Is("ok")}
	h1 := &rdgo.Listener{OnData: func(data rdgo.Data) {
		fmt.Printf("h1, %+v\n", data)
	}}
	h2 := &rdgo.Listener{OnData: func(data rdgo.Data) {
		fmt.Printf("h2, %+v\n", data)
	}}

	f.Select(c, h1)
	f.Select(c, h2)
	h1.Destroy()
	f.Input(Msg{"ok"})

	f.Input(Msg{"ok"})
	time.Sleep(time.Second)
	fmt.Printf("h1, last: %+v\n", h1.GetLast())
}

type Msg struct {
	Name string
}

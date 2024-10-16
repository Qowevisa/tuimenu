package main

import (
	"fmt"

	"git.qowevisa.me/Qowevisa/tuimenu/simple"
)

func main() {
	m := simple.CreateMenu(
		simple.WithCustomTitle("My Custom Title"),
	)
	nameName, err := m.AddCommand("0", "NameName", func(m *simple.Menu) error {
		fmt.Printf("hello, world\n")
		return nil
	})
	if err != nil {
		panic(err)
	}
	gr := nameName.AddGroupingCommand("gr1", "Some grouping")
	gr.AddExecCommand("some1", "Some func that does funny stuff", func(m *simple.Menu) error {
		fmt.Printf("Some funcy stuff\n")
		return nil
	})
	m.AddCommand("lalalla", "something", simple.EmptyAction)
	m.Start()
}

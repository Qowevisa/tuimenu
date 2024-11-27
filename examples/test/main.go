package main

import (
	"fmt"
	"log"
	"time"

	"git.qowevisa.me/qowevisa/tuimenu/simple"
)

func main() {
	m := simple.CreateMenu(
		simple.WithCustomTitle("My Custom Title"),
		simple.WithUsageOfEscapeCodes(),
	)
	// Using this function will redirect every log.PrintX func in THIS file
	//  to m.Log buffer that will Flush everything at the start of next iteration
	m.RedirectLogOutputToBufferedLogger()
	nameName, err := m.AddCommand("0", "NameName", func(m *simple.Menu) error {
		log.Printf("hello, world\n")
		time.Sleep(time.Second * 3)
		m.Log.Logf("some data = %d\n", 42)
		time.Sleep(time.Second * 1)
		m.Log.Logf("some data = %d\n", 43)
		time.Sleep(time.Second * 1)
		m.Log.Logf("some data = %d\n", 44)
		log.Printf("asdas")
		time.Sleep(time.Second * 1)
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

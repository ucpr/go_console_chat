package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/rivo/tview"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	form := tview.NewForm().AddInputField("Message", "", 0, nil, nil)
	form = form.AddButton("Submit", func() {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		txt := form.GetFormItem(0).(*tview.InputField).GetText()

		fmt.Fprintf(textView, "%v ** %s\n", timestamp, txt)
	})

	flex := tview.NewFlex().AddItem(textView, 0, 4, false).AddItem(form, 0, 2, false)

	fmt.Fprintf(textView, "%s\n", "hogehoge")

	textView.SetBorder(true).SetTitle("Go Console Chat")
	if err := app.SetRoot(flex, true).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

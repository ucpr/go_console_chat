package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			formated_txt := fmt.Sprintf("%s\n", message)
			fmt.Fprintf(textView, formated_txt)
		}
	}()

	form := tview.NewForm().AddInputField("Message", "", 0, nil, nil)

	form = form.AddButton("Submit", func() {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		txt := form.GetFormItem(0).(*tview.InputField).GetText()

		formated_txt := fmt.Sprintf("%v ** %s\n", timestamp, txt)

		err = c.WriteMessage(websocket.TextMessage, []byte(formated_txt))
		if err != nil {
			log.Println("write:", err)
			return
		}

		// writeしたメッセージはここでは描画しない
		// broadcastで送られてくるやつを描画する
	})

	flex := tview.NewFlex().AddItem(textView, 0, 4, false).AddItem(form, 0, 2, false)
	textView.SetBorder(true).SetTitle("Go Console Chat")
	if err = app.SetRoot(flex, true).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

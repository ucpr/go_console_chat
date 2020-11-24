package main

import (
	"flag"
	"fmt"
	"go_console_chat/model"
	"log"
	"net/url"
	"os/user"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Application struct {
	conn *websocket.Conn
	tui  *tview.Application
}

func NewApplication() *Application {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return &Application{
		conn: c,
		tui:  tview.NewApplication(),
	}
}

func username() string {
	cu, err := user.Current()
	if err != nil {
		log.Println(err)
		return "anonymous"
	}
	return cu.Username
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	done := make(chan struct{})

	name := username()

	app := NewApplication()
	defer app.conn.Close()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.tui.Draw()
		})

	go func() {
		defer close(done)
		for {
			_, message, err := app.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			fmt.Fprintln(textView, string(message))
		}
	}()

	form := tview.NewForm().AddInputField("Message", "", 0, nil, nil)
	form = form.AddButton("Submit", func() {
		txt := form.GetFormItem(0).(*tview.InputField).GetText()
		msg := model.NewMessage(name, txt)

		err := app.conn.WriteMessage(websocket.TextMessage, []byte(msg.ToText()))
		if err != nil {
			log.Println("write:", err)
			return
		}
		// broadcastで送られてくるやつを描画する
	})

	flex := tview.NewFlex().AddItem(textView, 0, 4, false).AddItem(form, 0, 2, false)
	textView.SetBorder(true).SetTitle("Go Console Chat")
	if err := app.tui.SetRoot(flex, true).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

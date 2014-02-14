package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/bitly/go-notify"
	"github.com/kballard/goirc/irc"
	"log"
	"math/rand"
	"net/http"
)

type NewUrl struct {
	Url string
}

type WebClientMessage struct {
	Type    string
	Message string
}

// a websocket connection to irc
type WebClient struct {
	nick string

	config  irc.Config
	irc     irc.SafeConn
	ircquit chan bool

	ws *websocket.Conn
}

func NewWebClient(ws *websocket.Conn) *WebClient {
	nick := fmt.Sprintf("comicchat%5d", rand.Intn(99999))
	wc := &WebClient{
		config: irc.Config{
			Host:     "chat.freenode.net",
			Nick:     nick,
			User:     "comicchat",
			RealName: "Comic Chat",
		},
		ws: ws,
	}

	wc.config.Init = func(hr irc.HandlerRegistry) {
		hr.AddHandler(irc.CONNECTED, wc.loggedin)
	}

	return wc
}

func (w *WebClient) loggedin(conn *irc.Conn, line irc.Line) {
	conn.Join([]string{"#comicchat"}, nil)
}

func (w *WebClient) reader() {
	var err error
	w.irc, err = irc.Connect(w.config)
	if err != nil {
		w.ws.Close()
	}
	for {
		var m WebClientMessage
		err = websocket.JSON.Receive(w.ws, &m)
		if err != nil {
			log.Printf("%s read: %s", w.ws.RemoteAddr(), err)
			break
		}

		if w.irc == nil {
			continue
		}

		switch m.Type {
		case "nick":
			w.irc.Nick(m.Message)
		case "privmsg":
			w.irc.Privmsg("#comicchat", m.Message)
		case "action":
			w.irc.Action("#comicchat", m.Message)
		case "connect":
		}
	}
}

func (w *WebClient) writer() {
	imgurlchan := make(chan interface{})
	notify.Start("newimage", imgurlchan)
	defer notify.Stop("newimage", imgurlchan)

loop:
	for {
		select {
		case url := <-imgurlchan:
			var m WebClientMessage
			m.Type = "newimage"
			m.Message = url.(string)
			if err := websocket.JSON.Send(w.ws, &m); err != nil {
				log.Printf("%s closed: %s", w.ws.RemoteAddr(), err)
				break loop
			}
		}
	}
}

func dohttp() {
	http.Handle("/new", websocket.Handler(newimagehandler))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	http.Handle("/comic/", http.StripPrefix("/comic/", http.FileServer(http.Dir("comic/"))))
	if err := http.ListenAndServe(":8899", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func newimagehandler(ws *websocket.Conn) {
	wc := NewWebClient(ws)
	go wc.writer()
	wc.reader()
}

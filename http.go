package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/bitly/go-notify"
	"github.com/kballard/goirc/irc"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	channel = "#comicchat"
)

type WebClientMessage struct {
	Type    string
	Message string
}

const (
	// fresh connection
	ClientNew int32 = iota
	// everything ok - ws and irc connected
	ClientOk
	// someone closed - kill everything
	ClientClosed
)

// a websocket connection to irc
type WebClient struct {
	mu      sync.Mutex
	remote  string
	nick    string
	status  int32
	config  irc.Config
	irc     irc.SafeConn
	ircquit chan bool
	ws      *websocket.Conn
}

func NewWebClient(remote string, ws *websocket.Conn) *WebClient {
	nick := fmt.Sprintf("comicchat%5.5d", rand.Intn(99999))
	wc := &WebClient{
		remote: remote,
		nick:   nick,
		status: ClientNew,
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

func (w WebClient) String() string {
	return fmt.Sprintf("%s %s", w.remote, w.nick)
}

func (w *WebClient) GetStatus() int32 {
	return atomic.LoadInt32(&w.status)
}

func (w *WebClient) SetStatus(s int32) int32 {
	return atomic.SwapInt32(&w.status, s)
}

// called when irc client connects
func (w *WebClient) loggedin(conn *irc.Conn, line irc.Line) {
	conn.Join([]string{"#comicchat"}, nil)
}

func (w *WebClient) Send(typ, msg string) error {
	if w.GetStatus() != ClientClosed {
		var m WebClientMessage

		w.mu.Lock()
		defer w.mu.Unlock()

		m.Type = typ
		m.Message = msg

		if err := websocket.JSON.Send(w.ws, &m); err != nil {
			return err
		}
	}

	return nil
}

// websocket read handler
func (w *WebClient) reader() {
	var err error

	w.mu.Lock()
	w.irc, err = irc.Connect(w.config)
	if err != nil {
		log.Printf("%s %s", w, err)
		w.ws.Close()
		w.mu.Unlock()
		return
	}
	w.mu.Unlock()

	w.Send("connected", "")

	defer func() {
		w.irc.Quit("disconnected")
	}()

	for w.GetStatus() != ClientClosed {
		var m WebClientMessage
		err = websocket.JSON.Receive(w.ws, &m)
		if err != nil {
			log.Printf("%s read: %s", w, err)
			w.SetStatus(ClientClosed)
			continue
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

// websocket write handler
func (w *WebClient) writer() {
	imgurlchan := make(chan interface{})
	notify.Start("newimage", imgurlchan)
	defer notify.Stop("newimage", imgurlchan)

	for w.GetStatus() != ClientClosed {
		select {
		case url := <-imgurlchan:
			if err := w.Send("newimage", url.(string)); err != nil {
				log.Printf("%s closed: %s", w, err)
				w.SetStatus(ClientClosed)
			}
		case <-time.After(30 * time.Second):
			if err := w.Send("ping", fmt.Sprintf("%d", time.Now().UnixNano())); err != nil {
				log.Printf("%s closed: %s", w, err)
				w.SetStatus(ClientClosed)
			}
		}
	}
}

func dohttp() {
	http.HandleFunc("/", indexhandler)
	http.HandleFunc("/new", websockethandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/comic/", http.StripPrefix("/comic/", http.FileServer(http.Dir("comic/"))))
	if err := http.ListenAndServe(":8899", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type Page struct {
	Faces []string
}

func indexhandler(w http.ResponseWriter, r *http.Request) {
	//create a new template
	indextmpl := template.Must(template.New("index.tpl").ParseFiles("template/index.tpl"))

	page := Page{
		Faces: facekeys,
	}
	if err := indextmpl.Execute(w, page); err != nil {
		log.Printf("%s", err)
	}
}

func websockethandler(w http.ResponseWriter, r *http.Request) {
	// get remote client's address. we may be behind a proxy (nginx)
	rem := r.RemoteAddr
	if strings.Split(rem, ":")[0] == "127.0.0.1" {
		rem = r.Header.Get("X-Real-IP")
	}

	websocket.Handler(func(ws *websocket.Conn) {
		wc := NewWebClient(rem, ws)
		log.Printf("%s connected", wc)
		go wc.writer()
		wc.reader()
		log.Printf("%s disconnected", wc)
	}).ServeHTTP(w, r)
}

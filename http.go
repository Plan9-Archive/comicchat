package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"github.com/bitly/go-notify"
	"log"
	"net/http"
)

type NewUrl struct {
	Url string
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
	imgurlchan := make(chan interface{})
	notify.Start("newimage", imgurlchan)
	defer notify.Stop("newimage", imgurlchan)

loop:
	for {
		select {
		case url := <-imgurlchan:
			var u NewUrl
			u.Url = url.(string)
			je := json.NewEncoder(ws)
			if err := je.Encode(&u); err != nil {
				log.Printf("%s closed: %s", ws.RemoteAddr(), err)
				break loop
			}
		}
	}
}

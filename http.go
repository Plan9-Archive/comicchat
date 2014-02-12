package main

import (
	"fmt"
	"github.com/bitly/go-notify"
	"log"
	"net/http"
)

func dohttp() {
	http.HandleFunc("/new", newimagehandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	http.Handle("/comic/", http.StripPrefix("/comic/", http.FileServer(http.Dir("comic/"))))
	if err := http.ListenAndServe(":4005", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func newimagehandler(c http.ResponseWriter, req *http.Request) {
	imgurlchan := make(chan interface{})
	notify.Start("newimage", imgurlchan)
	defer notify.Stop("newimage", imgurlchan)

	fmt.Fprintf(c, "%s\n", <-imgurlchan)
}

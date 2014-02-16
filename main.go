package main

import (
	"fmt"
	"github.com/bitly/go-notify"
	"github.com/kballard/goirc/irc"
	"log"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	quit := make(chan bool, 1)
	config := irc.Config{
		Host: "chat.freenode.net",

		Nick:     "maozedong",
		User:     "comicchat",
		RealName: "comicchat",

		Init: func(hr irc.HandlerRegistry) {
			hr.AddHandler(irc.CONNECTED, h_LoggedIn)
			hr.AddHandler(irc.DISCONNECTED, func(*irc.Conn, irc.Line) {
				log.Println("disconnected")
				quit <- true
			})
			hr.AddHandler("PRIVMSG", h_PRIVMSG)
			hr.AddHandler(irc.ACTION, h_ACTION)
		},
	}

	log.Println("Connecting")
	if _, err := irc.Connect(config); err != nil {
		fmt.Println("error:", err)
		quit <- true
	}

	go dohttp()

	<-quit
	log.Println("Goodbye")
}

func h_LoggedIn(conn *irc.Conn, line irc.Line) {
	conn.Join([]string{"#comicchat"}, nil)
}

var cnt = 0

func h_PRIVMSG(conn *irc.Conn, line irc.Line) {
	log.Printf("[%s] %s> %s\n", line.Args[0], line.Src, line.Args[1])
	if line.Args[1] == "!quit" {
		conn.Quit("")
	}
	i := makeusercomic(line.Src.Nick, line.Args[1])
	saveToPngFile(fmt.Sprintf("comic/%d.png", cnt), i)
	url := fmt.Sprintf("/comic/%d.png", cnt)

	notify.Post("newimage", url)

	cnt++
}

func h_ACTION(conn *irc.Conn, line irc.Line) {
	log.Printf("[%s] %s %s\n", line.Dst, line.Src, line.Args[0])
	f := strings.Fields(line.Args[0])
	if _, ok := faces[f[0]]; ok {
		userface[line.Src.Nick] = f[0]
	}
}

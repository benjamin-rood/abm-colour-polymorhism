package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/benjamin-rood/goio"
	"golang.org/x/net/websocket"
)

func networkError(err error, c chan struct{}) {
	log.Println(err)
	if err.Error() == "use of closed network connection" {
		close(c)
	}
}

func modelError(err error, c chan struct{}) {
	log.Println(err)
	// do something with the error value
	close(c)
}

func dataError(err error, c chan struct{}) {
	log.Println(err)
	// do something with the error value
	close(c)
}

const (
	sweepFreq   = time.Minute
	deathPeriod = time.Hour * 24
)

var socketUsers map[string]Client

func sweepSocketClients() {
	sweeper := time.NewTicker(sweepFreq)
	select {
	case <-sweeper.C:
		for uid, client := range socketUsers {
			if client.Dead {
				delete(socketUsers, uid)
				continue
			}
			if time.Since(client.Stamp) >= deathPeriod {
				delete(socketUsers, uid)
			}
		}
	}
}

func wsSession(ws *websocket.Conn) {
	uid := getUIDString()
	c := NewClient(ws, uid)
	socketUsers[uid] = c
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
		delete(socketUsers, uid)
	}()
	wsCh := make(chan struct{})
	go wsReader(ws, c.Im, wsCh)
	go wsWriter(ws, c.Om, wsCh)
	go c.Controller()
	c.Start()       //	demo
	c.Monitor(wsCh) //	keep alive
}

func wsReader(ws *websocket.Conn, in chan<- goio.InMsg, wsCh chan struct{}) {
	defer func() {
		//	clean up
	}()
	var msg goio.InMsg
	for {
		select {
		case <-wsCh:
			return
		default:
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				log.Println("Disconnected User.")
				close(wsCh)
				return
			}
			in <- msg //	gets picked up by Model controller function
		}
	}
}

func wsWriter(ws *websocket.Conn, out <-chan goio.OutMsg, quit <-chan struct{}) {
	defer func() {
		// clean up
	}()

	for {
		select {
		case <-quit:
			return
		case msg := <-out:
			websocket.JSON.Send(ws, msg)
		}
	}
}

func getUIDString() string {
again:
	id := randomdata.SillyName()
	_, exists := socketUsers[id]
	if exists {
		goto again
	}
	return id
}

func main() {
	rand.Seed(time.Now().UnixNano())
	go sweepSocketClients()
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", websocket.Handler(wsSession))
	http.ListenAndServe(":8080", nil)
}

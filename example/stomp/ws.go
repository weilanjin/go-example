package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-stomp/stomp/v3/frame"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		// mt, message, err := c.ReadMessage()
		_, r, err := c.NextReader()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fr := frame.NewReader(r)
		f, err := fr.Read()
		if err != nil {
			log.Println("read frame:", err)
			break
		}
		log.Printf("recv:")
		log.Printf("%+v", f)

		w, err := c.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println("NextWriter:", err)
			break
		}
		fw := frame.NewWriter(w)

		f = frame.New(frame.CONNECTED)
		h := frame.NewHeader()
		h.Add(frame.Version, "1.1")
		h.Add(frame.HeartBeat, "10000,0")
		h.Add("user-name", "")
		f.Header = h

		err = fw.Write(f)
		w.Close()
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

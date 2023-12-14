package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	port     = 9990
	upgrader = websocket.Upgrader{}
)

func main() {
	fmt.Println("Websocket Server!!!")

	http.HandleFunc("/echo", echo)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening on %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Print("read: ", err)
			break
		}
		log.Printf("recv: %s", message)
		if err := c.WriteMessage(mt, message); err != nil {
			log.Println("write: ", err)
			break
		}
	}
}

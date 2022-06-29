package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/libkarl/golang-chat-project/pkg/websocket"
)

// tohle je nadefinovaná funkce websocketu
func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request){
	fmt.Println("web socket end point reached")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client{
		Conn:  conn,
		Pool:  pool, 
	}
	// posílám vytvořeného klienta do kanálu v pool s názvem Register
	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	// web socket používá go rutinu
	go pool.Start()
	// zde je nadefinovananý end point web socketu
	http.HandleFunc("ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
}

func main() {
	
	fmt.Println("Starting the server on port 9000...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":9000", nil))
}
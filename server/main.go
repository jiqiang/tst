package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jiqiang/tst/server/message"
)

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second
)

//testtest
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func receiver(ws *websocket.Conn, out chan string) {
	defer ws.Close()
	ws.SetReadLimit(1024)
	ws.SetReadDeadline(time.Time{})
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := ws.ReadMessage()
		data := fmt.Sprintf("MESSAGE: %s", string(msg))
		out <- data
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func sender(ws *websocket.Conn, in chan string) {
	defer ws.Close()
	for {
		select {
		case message := <-in:
			err := ws.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}

func timer(out chan string) {
	for {
		data := message.Timer{
			Type: "TIMER",
			Time: time.Now().Format(time.ANSIC),
		}
		dataByte, _ := json.Marshal(data)
		out <- string(dataByte)
		time.Sleep(1 * time.Second)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()
	c := make(chan string)
	go sender(ws, c)
	go timer(c)
	go updateAssets(c)
	receiver(ws, c)
}

func getMockAssets() message.Assets {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	assets := []message.Asset{
		message.Asset{Name: "asset-1", TimeElapsedSinceLastUpdate: r.Intn(30)},
		message.Asset{Name: "asset-2", TimeElapsedSinceLastUpdate: r.Intn(30)},
		message.Asset{Name: "asset-3", TimeElapsedSinceLastUpdate: r.Intn(30)},
	}

	data := message.Assets{
		Type:   "ASSETS",
		Assets: assets,
	}

	return data
}

func updateAssets(out chan string) {
	for {
		data := getMockAssets()
		assetsStr, _ := json.Marshal(data)
		out <- string(assetsStr)
		time.Sleep(3 * time.Second)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWs)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("ui/"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:1234",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

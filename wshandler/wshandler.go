package wshandler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Wshandler struct {
	upgrader             websocket.Upgrader
	eventMapper          map[string]func()
	onConnectionCallback func(*Socket)
}

func New() *Wshandler {
	return &Wshandler{
		upgrader: websocket.Upgrader{},
	}
}

func (ws *Wshandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	socket := &Socket{
		client,
		make(chan *Message),
		make(map[string]func(string)),
	}

	if callback := ws.onConnectionCallback; callback != nil {
		go socket.read()
		go socket.write()
		callback(socket)
	}

}

func (ws *Wshandler) OnClientConnect(callback func(*Socket)) {
	ws.onConnectionCallback = callback
}

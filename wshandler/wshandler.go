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
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := NewSocket(conn)
	go client.read()
	go client.write()

	if callback := ws.onConnectionCallback; callback != nil {
		callback(client)
	}

}

func (ws *Wshandler) OnClientConnect(callback func(*Socket)) {
	ws.onConnectionCallback = callback
}

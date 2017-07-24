package wshandler

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Socket struct {
	client      *websocket.Conn
	send        chan *Message
	eventMapper map[string]func(data string)
	forEveryEvent
}

func New(conn *websocket.Conn) *Socket {
	return &Socket{
		eventMapper: make(map[string]func(string)),
		client:      conn,
		send:        make(chan *Message),
	}
}

func (s *Socket) On(event string, callback func(data string)) {
	s.eventMapper[event] = callback
}

func (s *Socket) Every(callback func(data string)) {

}

func (s *Socket) Default(callback func(data string)) {

}

func (s *Socket) Emit(event string, data string) {
	s.send <- &Message{
		Event: event,
		Data:  data,
	}
}

func (s *Socket) read() {
	for {
		msg := &Message{}
		err := s.client.ReadJSON(msg)

		if err != nil {
			s.client.Close()
			fmt.Println("SOCKET ERROR:", err)
			return
		}

		if callback, ok := s.eventMapper[msg.Event]; ok {
			callback(msg.Data)
		} else {

		}

	}
}

func (s *Socket) write() {
	for {
		select {
		case msg := <-s.send:
			err := s.client.WriteJSON(msg)
			if err != nil {
				fmt.Println("SOCKET WRITE ERROR", err)
				s.client.Close()
				return
			}
		}
	}
}

package wshandler

import (
	"github.com/gorilla/websocket"
)

type Socket struct {
	client          *websocket.Conn
	send            chan *Message
	eventMapper     map[string]func(data interface{})
	forEveryEvent   func(*Message)
	forDefaultEvent func(*Message)
}

func NewSocket(conn *websocket.Conn) *Socket {
	return &Socket{
		eventMapper: make(map[string]func(interface{})),
		client:      conn,
		send:        make(chan *Message),
	}
}

func (s *Socket) On(event string, callback func(data interface{})) {
	s.eventMapper[event] = callback
}

func (s *Socket) OnEveryMessage(callback func(message *Message)) {
	s.forEveryEvent = callback
}

func (s *Socket) OnDefaultMessage(callback func(message *Message)) {
	s.forDefaultEvent = callback
}

func (s *Socket) Emit(event string, data interface{}) {
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
			break
		}

		s.dispatchMessageToCallback(msg)

	}
}

func (s *Socket) write() {
	defer s.client.Close()
	defer close(s.send)
	for {
		select {
		case msg := <-s.send:
			err := s.client.WriteJSON(msg)
			if err != nil {
				break
			}
		}
	}
}

func (s *Socket) dispatchMessageToCallback(message *Message) {
	if eventHandler, ok := s.eventMapper[message.Event]; ok {
		eventHandler(message.Data)

	} else if defaultHandler := s.forDefaultEvent; defaultHandler != nil {
		defaultHandler(message)
	}
	if callback := s.forEveryEvent; callback != nil {
		callback(message)
	}
}

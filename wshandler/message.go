package wshandler

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

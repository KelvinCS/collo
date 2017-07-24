package main

import (
	"collo/filewatcher"
	"collo/wshandler"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	watcher := filewatcher.New("./", handleFileChange)
	watcher.Start()

	ws := wshandler.New()
	ws.OnClientConnect(func(socket *wshandler.Socket) {

		fmt.Println("CLIENT CONNECTED")
		socket.Emit("Hello", "Teste")
		socket.On("foo", func(data interface{}) {
			fmt.Println(data)
		})

		socket.OnDefaultMessage(func(msg *wshandler.Message) {
			fmt.Println("Default: ", msg)
		})

		socket.OnEveryMessage(func(msg *wshandler.Message) {
			fmt.Println("Every: ", msg)
		})
	})

	http.HandleFunc("/", homeController)
	http.Handle("/ws", ws)
	http.ListenAndServe(":3000", nil)
}

func handleFileChange(path string, eventName string) {
	fmt.Println("Arquivo mudou:", path)
}

func homeController(w http.ResponseWriter, r *http.Request) {
	homeTemplate, _ := template.ParseFiles("./static/index.html")
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

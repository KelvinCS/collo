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
		socket.Emit("Hello", "Batatapalha")
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

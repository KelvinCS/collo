package main

import (
	"collo/filewatcher"
	"collo/wshandler"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var fileChange chan string

func main() {
	fmt.Println(os.Args[1])
	fileChange = make(chan string)

	watcher := filewatcher.New(os.Args[1], handleFileChange)
	watcher.Start()

	ws := wshandler.New()

	ws.OnClientConnect(func(socket *wshandler.Socket) {
		go func() {
			for {
				select {
				case file := <-fileChange:
					socket.Emit("CHANGE", file)
				}
			}
		}()
	})

	fs := http.FileServer(http.Dir(os.Args[1]))

	http.Handle("/watcher/", http.StripPrefix("/watcher", fs))
	http.HandleFunc("/", homeController)
	http.Handle("/ws", ws)
	http.ListenAndServe(":3000", nil)
}

func handleFileChange(path string, eventName string) {
	fileChange <- path
}

func homeController(w http.ResponseWriter, r *http.Request) {
	homeTemplate, _ := template.ParseFiles("./static/index.html")
	homeTemplate.Execute(w, r.Host)
}

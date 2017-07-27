package main

import (
	"collo/filewatcher"
	"collo/wshandler"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os/exec"
	"time"
)

var fileChange = make(chan string)
var prefix string

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	prefix = fmt.Sprintf("%d", rand.Int())

	watchDir := flag.String("w", "./", "Root dir to watch for modifications")
	port := flag.String("p", ":3000", "Port to listen server. Default is port :3000")
	flag.Parse()

	watcher := filewatcher.New(*watchDir, handleFileChange)
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

	fs := http.FileServer(http.Dir(*watchDir))

	path := "/watcher" + prefix
	fmt.Println(path)
	http.HandleFunc("/", homeController)
	http.Handle(path+"/", http.StripPrefix(path, fs))
	http.Handle("/ws", ws)

	fmt.Println("Open your web browser and go to `http://localhost" + *port + "`")
	exec.Command("google-chrome-stable", "http://localhost"+*port).Run()

	err := http.ListenAndServe(*port, nil)
	fmt.Println(err)
}

func handleFileChange(path string, eventName string) {
	fmt.Println(path, eventName)
	fileChange <- path
}

func homeController(w http.ResponseWriter, r *http.Request) {
	homeTemplate, _ := template.New("").Parse(HomeTemplate)
	homeTemplate.Execute(w, struct {
		Prefix string
		Host   string
	}{
		prefix,
		r.Host,
	})
}

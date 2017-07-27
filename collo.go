package main

import (
	"collo/filewatcher"
	"collo/wshandler"
	"flag"
	"fmt"
	"html/template"
	"net/http"
)

var fileChange = make(chan string)

func main() {

	watchDir := flag.String("watch", "./", "Root dir to watch for modifications")
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

	http.HandleFunc("/", homeController)
	http.Handle("/watcher/", http.StripPrefix("/watcher", fs))
	http.Handle("/ws", ws)

	fmt.Println("Open your web browser and go to `http://localhost:3000`")
	http.ListenAndServe(":3000", nil)
}

func handleFileChange(path string, eventName string) {
	fmt.Println(path)
	fileChange <- path
}

func homeController(w http.ResponseWriter, r *http.Request) {
	homeTemplate, _ := template.New("").Parse(HomeTemplate)
	homeTemplate.Execute(w, r.Host)
}

const HomeTemplate = `
	
<!DOCTYPE html>
<html lang="en">
  <head>
    <title></title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
  </head>
  <style>
    body {
      margin: 0;
      padding: 0;
      overflow: hidden;
    }
    iframe {
      margin: 0;
      width: 100vw;
      height: 100vh;
    }
  
  </style>
  <body>
    <iframe id="view" src="http://{{.}}/watcher" frameborder="0"></iframe>
  </body>
  <script>
    let iframe = document.querySelector("#view");
    let view = iframe.contentWindow || iframe.contentDocument.document || iframe.contentDocument;

    
		let socket = new WebSocket("ws://{{.}}/ws")
		
    socket.onmessage = (message) => {
      console.log(message)
      iframe.src = iframe.src
    }

    socket.onerror = (err) => {
      console.log(err)
    }

    socket.onopen = (q) => {
      console.log(q)
    }
  </script>
</html>

`

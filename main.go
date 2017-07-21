package main

import (
	"collo/filewatcher"
	"fmt"
)

func main() {
	watcher := filewatcher.New("./", handleFileChange)
	watcher.Start()
	<-make(chan bool)
}

func handleFileChange(path string, eventName string) {
	fmt.Println("Arquivo mudou:", path)
}

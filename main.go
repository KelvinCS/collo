package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	done := make(chan bool)

	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if isCreateDirEvent(event) {
					watcher.Add(event.Name)
				}
				fmt.Println(event)

			case err := <-watcher.Errors:
				fmt.Println(err)
			}
		}
	}()

	<-done
}

/*
* TODO: error treatment
 */
func isCreateDirEvent(event fsnotify.Event) bool {
	if event.Op.String() != "CREATE" {
		return false
	}
	path := event.Name
	file, _ := os.Open(path)
	stat, _ := file.Stat()
	return stat.IsDir()
}

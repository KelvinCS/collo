package filewatcher

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
)

type Callback func(path string, eventName string)

type Watcher struct {
	baseDir  string
	watcher  *fsnotify.Watcher
	callback Callback
}

/* New: Create new watcher to watch modifications in a direcatory, recursively
 */
func New(baseDir string, callback Callback) *Watcher {
	watcher, _ := fsnotify.NewWatcher()
	return &Watcher{
		baseDir,
		watcher,
		callback,
	}
}

func (w *Watcher) Start() error {
	err := w.walkAndAddEveryDir(w.baseDir)
	go w.watchToModifications()

	fmt.Println("Stated")
	return err

}

func (w *Watcher) walkAndAddEveryDir(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})

	return err
}

func (w *Watcher) watchToModifications() {
	defer w.Close()

	for {
		select {
		case event := <-w.watcher.Events:
			w.handleModification(event)

		case err := <-w.watcher.Errors:
			fmt.Println(err)
		}
	}
}

func (w *Watcher) handleModification(event fsnotify.Event) {
	eventName := event.Op.String()
	path := event.Name
	if eventName == "CREATE" {
		w.watcher.Add(path)
	}
	if eventName == "RENAME" {
		w.watcher.Remove(path)
	}
	w.callback(path, eventName)
}

func (w *Watcher) Close() {
	w.watcher.Close()
}

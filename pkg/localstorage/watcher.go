package localstorage

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func watchFileChanges(filePath string, changeChan chan bool) {
	// watcher is used to watch for changes to the file
	// it continuously listens for events
	// when a change is detected, it sends a boolean value to the changeChan
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatalf("Failed to add file to watcher: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			// check if the channel is closed
			if !ok {
				return
			}
			// check if the file was modified
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("File %s was modified", event.Name)
				changeChan <- true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

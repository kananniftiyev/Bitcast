package filemonitor

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func WatchFolder(folderPath string) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	if err = watcher.Add(folderPath); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Print the type of event
			log.Println("event:", event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}

func ListOfFiles(folderPath string) {
	dir, err := os.Open(folderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer dir.Close()

	fileInfos, err := dir.Readdir(-1) // -1 means return all entries
	if err != nil {
		fmt.Println("Error reading folder contents:", err)
		return
	}

	if len(fileInfos) == 0 {
		fmt.Println("Folder is empty")
		return
	}

	// Print the names of the files
	fmt.Println("Files in", folderPath, ":")
	for _, fileInfo := range fileInfos {
		fmt.Println(fileInfo.Name())
	}
}

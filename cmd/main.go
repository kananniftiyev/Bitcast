package main

import filemonitor "fileguard/internal/file"

// Change folder path to test.
const folderPath = `C:\\Users\\kenan\\Desktop\\test`

func main() {
	filemonitor.ListOfFiles(folderPath)
	filemonitor.WatchFolder(folderPath)
}

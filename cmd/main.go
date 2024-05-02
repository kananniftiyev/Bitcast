package main

import (
	"fileguard/internal/server"
	"log"
)

// Change folder path to test.
const folderPath = `C:\Users\kenan\Desktop\test`

func main() {
	// CreateWindowsService(folderPath)
	/* newDB, err := db.NewDatabase("fileguard-cf4d3")
	if err != nil {
		log.Fatal(err)
	}

	_ = newDB*/

	//auth.LoginViaGoogle()

	err := server.UploadFile("C:/Users/kenan/Documents/GitHub/fileguard/internal/server/test.txt", "asd")
	if err != nil {
		log.Println(err)
	}
}

package main

import (
	"fileguard/internal/server"
	"log"
	filemonitor "fileguard/internal/file"
	"github.com/kardianos/service"
	"log"
	"os"
)

// Change folder path to test.
const folderPath = `C:\Users\kenan\Desktop\test`

func main() {
	//auth.LoginViaGoogle()

	err := server.UploadFile("C:/Users/kenan/Documents/GitHub/fileguard/internal/server/test.txt", "asd")
	if err != nil {
		log.Println(err)
	}
}

package main

import (
	"fileguard/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

// Change folder path to test.
const folderPath = `C:\Users\kenan\Desktop\test`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	//auth.LoginViaGoogle()

	/*err := server.UploadFile("C:/Users/kenan/Documents/GitHub/fileguard/internal/server/test.txt", "asd")
	if err != nil {
		log.Println(err)
	}*/

	err := server.DownloadAllFiles("", "C:/Users/kenan/Documents/GitHub/fileguard")
	if err != nil {
		fmt.Println(err)
	}
}

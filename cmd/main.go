package main

import (
	"fileguard/internal/auth"
)

// Change folder path to test.
const folderPath = `C:\Users\kenan\Desktop\test`

func main() {
	auth.LoginViaGoogle()

	/*err := storage.UploadFile("C:/Users/kenan/Documents/GitHub/fileguard/internal/storage/test.txt", "asd")
	if err != nil {
		log.Println(err)
	}*/

	/*err := storage.DownloadAllFiles("", "C:/Users/kenan/Documents/GitHub/fileguard")
	if err != nil {
		fmt.Println(err)
	}*/
}

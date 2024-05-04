package main

import (
	"fileguard/internal/auth"
)

// Change folder path to test.
const folderPath = `C:\Users\kenan\Desktop\test`

func main() {
	auth.LoginViaGoogle()

	/*f, err := os.Getwd()

	storage, err := firebase.NewStorage()

	if err != nil {
		fmt.Println(err)
	}

	err = storage.DownloadFile("x/file.txt", f+"\\file.txt")
	if err != nil {
		fmt.Println(err)
	}*/
}

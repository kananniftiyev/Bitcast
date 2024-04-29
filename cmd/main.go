package main

import "fileguard/internal/auth"

// Change folder path to test.
const folderPath = `C:\Users\kenan\Desktop\test`

func main() {
	// CreateWindowsService(folderPath)
	/* newDB, err := db.NewDatabase("fileguard-cf4d3")
	if err != nil {
		log.Fatal(err)
	}

	_ = newDB*/

	auth.LoginViaGoogle()

}

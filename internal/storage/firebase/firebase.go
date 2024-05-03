// Package firebase provides functionalities to interact with Firebase services such as Firestore and Storage.
package firebase

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Add Google Drive and firebase storage
// Give users choice to either implement their own google drive or use firebase storage by us.
// 200MB each user for firebase storage.

const maxFolderSize = 200 * 1024 * 1024

func initStorage() (*storage.BucketHandle, error) {
	config := &firebase.Config{
		StorageBucket: "fileguard-cf4d3.appspot.com",
	}

	opt := option.WithCredentialsFile("C:/Users/kenan/Documents/GitHub/fileguard/fileguard.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// TODO: Detect file content difference.
// TODO: Encode user infos from Google OAuth2
func UploadFile(localFilePath string, userToken string) error {
	bucket, err := initStorage()
	if err != nil {
		return err
	}
	ctx := context.Background()

	folderPath := "x"

	totalFolderSize, err := GetFolderSize(bucket, folderPath, ctx)

	if err != nil {
		return err
	}

	if totalFolderSize >= maxFolderSize {
		return errors.New("Cannot surpass max folder size")
	}

	// Open the file
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file content into a byte slice
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Username will be folder name
	dst := "y/" + file.Name()

	// Create a writer
	w := bucket.Object(dst).NewWriter(ctx)

	// Write the content to the writer
	if _, err := w.Write(content); err != nil {
		return err
	}

	// Close the writer
	if err := w.Close(); err != nil {
		return err
	}

	log.Println("File uploaded successfully!")
	return nil
}

// TODO: Redesign this better.
func DownloadFile(objectPath string, localPath string) error {
	bucket, err := initStorage()

	if err != nil {
		println(err)
	}

	ctx := context.Background()
	rc, err := bucket.Object(objectPath).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create local file
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy object data to local file
	if _, err := io.Copy(file, rc); err != nil {
		return err
	}

	return nil
}

func DownloadAllFiles(folderPath string, localPath string) error {
	if folderPath == "" {
		return errors.New("You should add Folder Path")
	}

	bucket, err := initStorage()
	if err != nil {
		return err
	}

	it := bucket.Objects(context.Background(), &storage.Query{Prefix: folderPath})

	for {
		objs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}
		fmt.Println(objs.Name)
		fileDir := filepath.Join(localPath, filepath.Base(objs.Name))
		err = DownloadFile(objs.Name, fileDir)
		if err != nil {
			return err
		}

		log.Printf("Downloaded file: %s\n", objs.Name)
	}

	return nil

}

func GetFolderSize(bucket *storage.BucketHandle, folderPathInStorage string, ctx context.Context) (int64, error) {
	var totalFolderSize int64
	it := bucket.Objects(ctx, &storage.Query{Prefix: folderPathInStorage})

	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return 0, err
		}

		totalFolderSize += objAttrs.Size
	}

	return totalFolderSize, nil
}

package server

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Add Google Drive and firebase storage
// Give users choice to either implement their own google drive or use firebase storage by us.
// 500MB each user for firebase storage.

const maxFolderSize = 200 * 1024 * 1024

func initStorage() (*storage.BucketHandle, error) {
	config := &firebase.Config{
		StorageBucket: "fileguard-cf4d3.appspot.com",
	}

	opt := option.WithCredentialsFile("C:/Users/kenan/Documents/GitHub/fileguard/internal/db/fileguard.json")
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

// TODO: Encode user infos from Google OAuth2
func UploadFile(filePath string, userToken string) error {
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
	file, err := os.Open(filePath)
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
	dst := "x/file.txt"

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

func DownloadFile(objectPath string, localPath string) (*os.File, error) {
	bucket, err := initStorage()

	if err != nil {
		println(err)
	}

	ctx := context.Background()
	rc, err := bucket.Object(objectPath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	// Create local file
	file, err := os.Create(localPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Copy object data to local file
	if _, err := io.Copy(file, rc); err != nil {
		return nil, err
	}

	return file, nil
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

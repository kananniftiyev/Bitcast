// Package provides way to CRUD operations with File Collection

package db

import (
	"cloud.google.com/go/firestore"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"os"
	"time"
)

// TODO: find way to duplicate file upload. Normaly 2 file names can be in firestore so...
func (db *Database) CreateNewFileRecord(f *os.File, userID string) error {
	fileInfo, err := f.Stat()

	if err != nil {
		return err
	}

	w := map[string]interface{}{
		"file_id":       uuid.New().String(),
		"file_name":     f.Name(),
		"created_at":    time.Now(),
		"file_size":     fileInfo.Size(),
		"user_id":       userID,
		"last_modified": fileInfo.ModTime(),
	}
	_, r, err := db.Client.Collection("Files").Add(db.ctx, w)
	if err != nil {
		return err
	}
	fmt.Println(r)
	return nil
}

// TODO: Query code is same for user and file so move it to utils.
func (db *Database) GetFileRecordById(fileID string) (map[string]interface{}, error) {
	query := db.Client.Collection("Files").Where("file_id", "==", fileID).Limit(1)

	iter := query.Documents(db.ctx)
	defer iter.Stop()

	doc, err := iter.Next()

	if err == iterator.Done {
		return nil, errors.New("File record not found.")
	}
	if err != nil {
		return nil, err
	}

	return doc.Data(), err
}

func (db *Database) DeleteFileRecordById(fileID string) error {
	query := db.Client.Collection("Files").Where("file_id", "==", fileID).Limit(1)

	iter := query.Documents(db.ctx)
	defer iter.Stop()

	doc, err := iter.Next()

	if err == iterator.Done {
		return errors.New("File record not found.")
	}
	if err != nil {
		return err
	}

	_, err = doc.Ref.Delete(db.ctx)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateFileRecordById(fileID string, newData map[string]interface{}) error {
	query := db.Client.Collection("Files").Where("file_id", "==", fileID).Limit(1)

	iter := query.Documents(db.ctx)
	defer iter.Stop()

	doc, err := iter.Next()

	if err == iterator.Done {
		return errors.New("File record not found.")
	}
	if err != nil {
		return err
	}

	var updates []firestore.Update
	for key, value := range newData {
		update := firestore.Update{Path: key, Value: value}
		updates = append(updates, update)
	}

	_, err = doc.Ref.Update(db.ctx, updates)

	if err != nil {
		return err
	}

	return nil
}

// Package provides way to CRUD operations with File Collection

package db

import (
	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"os"
	"time"
)

// TODO: find way to duplicate file upload. Normaly 2 file names can be in firestore so...
func (db *Database) CreateNewFileRecord(f *os.File, userID string) error {
	fileInfo, err := f.Stat()

	if err != nil {
		return err
	}

	// TODO: bytes, MB, what fix it.
	w := map[string]interface{}{
		"file_id":       uuid.New().String(),
		"file_name":     f.Name(),
		"created_at":    time.Now(),
		"file_size":     fileInfo.Size(),
		"user_id":       userID,
		"last_modified": fileInfo.ModTime(),
	}
	_, _, err = db.Client.Collection("Files").Add(db.ctx, w)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Query code is same for user and file so move it to utils.
func (db *Database) GetFileRecordById(fileID string) (map[string]interface{}, error) {
	doc, err := db.getRecord("Files", "file_id", fileID)
	if err != nil {
		return nil, err
	}

	return doc.Data(), nil
}

func (db *Database) DeleteFileRecordById(fileID string) error {
	doc, err := db.getRecord("Files", "file_id", fileID)
	if err != nil {
		return err
	}

	if _, err = doc.Ref.Delete(db.ctx); err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateFileRecordById(fileID string, newData map[string]interface{}) error {
	doc, err := db.getRecord("Files", "file_id", fileID)
	if err != nil {
		return err
	}

	var updates []firestore.Update
	for key, value := range newData {
		update := firestore.Update{Path: key, Value: value}
		updates = append(updates, update)
	}

	if _, err = doc.Ref.Update(db.ctx, updates); err != nil {
		return err
	}

	return nil
}

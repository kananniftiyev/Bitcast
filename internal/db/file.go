// Package provides way to CRUD operations with File Collection

package db

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

// Human-readable file size formatting
func formatFileSize(size int64) string {
	const (
		KB = 1 << (10 * 1) // 1024 bytes
		MB = 1 << (10 * 2) // 1024 KB
		GB = 1 << (10 * 3) // 1024 MB
		TB = 1 << (10 * 4) // 1024 GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

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
		"file_size":     formatFileSize(fileInfo.Size()),
		"user_id":       userID,
		"last_modified": fileInfo.ModTime(),
	}
	_, _, err = db.Client.Collection("Files").Add(db.ctx, w)
	if err != nil {
		return err
	}
	return nil
}

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

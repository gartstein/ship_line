// Package bolt provides a BoltDB-based implementation of the PackRepository interface.
// It allows for storing, retrieving, updating, and deleting pack sizes.
package bolt

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("packSizes")

// BoltStorage implements the PackRepository interface.
type BoltStorage struct {
	db *bolt.DB
}

// NewBoltStorage opens or creates a BoltDB database at the specified path and
// ensures the required bucket exists.
func NewBoltStorage(path string) (*BoltStorage, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	// Ensure bucket exists.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &BoltStorage{db: db}, nil
}

// GetPackSizes retrieves the stored pack sizes from the BoltDB.
// If no pack sizes are stored, it returns an empty slice.
func (b *BoltStorage) GetPackSizes() ([]int, error) {
	var sizes []int
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		data := bucket.Get([]byte("packSizes"))
		if data == nil {
			return nil // not stored yet
		}
		return json.Unmarshal(data, &sizes)
	})
	return sizes, err
}

// SetPackSizes stores the given pack sizes in BoltDB.
// It merges with existing sizes and removes duplicates.
func (b *BoltStorage) SetPackSizes(sizes []int) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve existing pack sizes
		var existingSizes []int
		data := bucket.Get([]byte("packSizes"))
		if data != nil {
			if err := json.Unmarshal(data, &existingSizes); err != nil {
				return err
			}
		}

		// Merge and remove duplicates
		sizeSet := make(map[int]struct{})
		for _, s := range existingSizes {
			sizeSet[s] = struct{}{}
		}
		for _, s := range sizes {
			sizeSet[s] = struct{}{}
		}

		// Convert map keys back to a slice
		updatedSizes := make([]int, 0, len(sizeSet))
		for s := range sizeSet {
			updatedSizes = append(updatedSizes, s)
		}

		// Store the updated pack sizes
		newData, err := json.Marshal(updatedSizes)
		if err != nil {
			return err
		}

		return bucket.Put([]byte("packSizes"), newData)
	})
}

// DeletePackSize removes a specific pack size from storage if it exists.
func (b *BoltStorage) DeletePackSize(size int) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve existing pack sizes
		var existingSizes []int
		data := bucket.Get([]byte("packSizes"))
		if data != nil {
			if err := json.Unmarshal(data, &existingSizes); err != nil {
				return err
			}
		}

		// Filter out the size to be deleted
		updatedSizes := make([]int, 0, len(existingSizes))
		for _, s := range existingSizes {
			if s != size {
				updatedSizes = append(updatedSizes, s)
			}
		}

		// Store the updated pack sizes
		newData, err := json.Marshal(updatedSizes)
		if err != nil {
			return err
		}

		return bucket.Put([]byte("packSizes"), newData)
	})
}

// Close closes the Bolt database.
func (b *BoltStorage) Close() error {
	return b.db.Close()
}

package bolt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoltStorage_GetAndSetPackSizes(t *testing.T) {
	// Create a temporary file for the Bolt DB.
	tempDB := "test_pack_sizes.db"
	defer os.Remove(tempDB)

	storage, err := NewBoltStorage(tempDB)
	assert.NoError(t, err)
	defer storage.Close()

	// Initially, GetPackSizes should return an empty slice.
	sizes, err := storage.GetPackSizes()
	assert.NoError(t, err)
	assert.Empty(t, sizes)

	// Set some pack sizes.
	expected := []int{1000, 500, 250}
	err = storage.SetPackSizes(expected)
	assert.NoError(t, err)

	// Get them back and verify.
	sizes, err = storage.GetPackSizes()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, sizes)
}

func TestBoltStorage_ErrorAfterClose(t *testing.T) {
	t.Run("GetPackSizes after close", func(t *testing.T) {
		// Create a temporary DB file.
		tempDB := "test_fail_get.db"
		defer os.Remove(tempDB)

		storage, err := NewBoltStorage(tempDB)
		assert.NoError(t, err)

		// Close the DB.
		err = storage.Close()
		assert.NoError(t, err)

		// Now, calling GetPackSizes should return an error.
		_, err = storage.GetPackSizes()
		assert.Error(t, err, "expected error when calling GetPackSizes on closed DB")
	})

	t.Run("SetPackSizes after close", func(t *testing.T) {
		// Create a temporary DB file.
		tempDB := "test_fail_set.db"
		defer os.Remove(tempDB)

		storage, err := NewBoltStorage(tempDB)
		assert.NoError(t, err)

		// Close the DB.
		err = storage.Close()
		assert.NoError(t, err)

		// Now, calling SetPackSizes should return an error.
		err = storage.SetPackSizes([]int{1000, 500})
		assert.Error(t, err, "expected error when calling SetPackSizes on closed DB")
	})
}

func TestBoltStorage_DeletePackSize(t *testing.T) {
	t.Run("DeleteExistingSize", func(t *testing.T) {
		tmpFile := filepath.Join(os.TempDir(), "test_existing.db")
		defer os.Remove(tmpFile)

		storage, err := NewBoltStorage(tmpFile)
		if err != nil {
			t.Fatalf("failed to create BoltStorage: %v", err)
		}
		defer storage.Close()

		// Initialize pack sizes.
		initialSizes := []int{1, 2, 3, 4}
		if err := storage.SetPackSizes(initialSizes); err != nil {
			t.Fatalf("failed to set pack sizes: %v", err)
		}

		// Delete an existing size.
		if err := storage.DeletePackSize(3); err != nil {
			t.Fatalf("DeletePackSize failed: %v", err)
		}

		sizes, err := storage.GetPackSizes()
		if err != nil {
			t.Fatalf("failed to get pack sizes: %v", err)
		}

		expected := []int{4, 2, 1}
		assert.ElementsMatch(t, expected, sizes, fmt.Sprintf("expected sizes %v, got %v", expected, sizes))
	})

	t.Run("DeleteNonExistingSize", func(t *testing.T) {
		tmpFile := filepath.Join(os.TempDir(), "test_non_existing.db")
		defer os.Remove(tmpFile)

		storage, err := NewBoltStorage(tmpFile)
		if err != nil {
			t.Fatalf("failed to create BoltStorage: %v", err)
		}
		defer storage.Close()

		// Initialize pack sizes without the size to be deleted.
		initialSizes := []int{1, 2, 4}
		if err := storage.SetPackSizes(initialSizes); err != nil {
			t.Fatalf("failed to set pack sizes: %v", err)
		}

		// Attempt to delete a non-existing size.
		if err := storage.DeletePackSize(3); err != nil {
			t.Fatalf("DeletePackSize failed: %v", err)
		}

		sizes, err := storage.GetPackSizes()
		if err != nil {
			t.Fatalf("failed to get pack sizes: %v", err)
		}

		expected := []int{1, 2, 4}
		assert.ElementsMatch(t, expected, sizes, fmt.Sprintf("expected sizes %v, got %v", expected, sizes))
	})

	t.Run("DeleteAlreadyDeletedSize", func(t *testing.T) {
		tmpFile := filepath.Join(os.TempDir(), "test_already_deleted.db")
		defer os.Remove(tmpFile)

		storage, err := NewBoltStorage(tmpFile)
		if err != nil {
			t.Fatalf("failed to create BoltStorage: %v", err)
		}
		defer storage.Close()

		// Initialize pack sizes.
		initialSizes := []int{1, 2, 3, 4}
		if err := storage.SetPackSizes(initialSizes); err != nil {
			t.Fatalf("failed to set pack sizes: %v", err)
		}

		// Delete size 3 the first time.
		if err := storage.DeletePackSize(3); err != nil {
			t.Fatalf("first DeletePackSize failed: %v", err)
		}
		// Delete size 3 a second time.
		if err := storage.DeletePackSize(3); err != nil {
			t.Fatalf("second DeletePackSize failed: %v", err)
		}

		sizes, err := storage.GetPackSizes()
		if err != nil {
			t.Fatalf("failed to get pack sizes: %v", err)
		}

		expected := []int{1, 2, 4}
		assert.ElementsMatch(t, expected, sizes, fmt.Sprintf("expected sizes %v, got %v", expected, sizes))
	})
}

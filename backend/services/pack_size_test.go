package services

import (
	"errors"
	"reflect"
	"ship_line/swagger"
	"testing"
)

// Mock repository for testing
type mockPackRepo struct {
	sizes []int
	err   error
}

func (m *mockPackRepo) GetPackSizes() ([]int, error) {
	return m.sizes, m.err
}

func (m *mockPackRepo) SetPackSizes(sizes []int) error {
	if m.err != nil {
		return m.err
	}
	m.sizes = sizes
	return nil
}

func (m *mockPackRepo) DeletePackSize(size int) error {
	return nil
}

func TestGetPackSizes(t *testing.T) {
	t.Run("ValidPackSizes", func(t *testing.T) {
		mockRepo := &mockPackRepo{sizes: []int{50, 30, 10}}
		service := NewPackService(mockRepo)

		expectedResponse := &swagger.PackSizesPayload{
			PackSizes: []int{50, 30, 10},
		}

		resp, err := service.GetPackSizes()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(resp, expectedResponse) {
			t.Errorf("expected response: %+v, got: %+v", expectedResponse, resp)
		}
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := &mockPackRepo{err: errors.New("database error")}
		service := NewPackService(mockRepo)

		_, err := service.GetPackSizes()
		if err == nil || err.Error() != "failed to retrieve pack sizes" {
			t.Errorf("expected repository error, got: %v", err)
		}
	})
}

func TestUpdatePackSizes(t *testing.T) {
	mockRepo := &mockPackRepo{}
	service := NewPackService(mockRepo)

	t.Run("ValidPackSizes", func(t *testing.T) {
		err := service.UpdatePackSizes([]int{50, 30, 10})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("EmptyPackSizes", func(t *testing.T) {
		err := service.UpdatePackSizes([]int{})
		if err == nil || err.Error() != "pack sizes cannot be empty" {
			t.Errorf("expected error 'pack sizes cannot be empty', got: %v", err)
		}
	})

	t.Run("NegativePackSize", func(t *testing.T) {
		err := service.UpdatePackSizes([]int{10, -5, 20})
		if err == nil || err.Error() != "invalid pack size: -5 (must be positive)" {
			t.Errorf("expected error for negative pack size, got: %v", err)
		}
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepoWithError := &mockPackRepo{err: errors.New("database error")}
		serviceWithError := NewPackService(mockRepoWithError)

		err := serviceWithError.UpdatePackSizes([]int{10, 20})
		if err == nil || err.Error() != "database error" {
			t.Errorf("expected database error, got: %v", err)
		}
	})
}

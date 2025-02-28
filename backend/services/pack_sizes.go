package services

import (
	"errors"
	"fmt"
	"ship_line/swagger"
	"sort"
)

// GetPackSizes retrieves the available pack sizes.
func (ps *PackService) GetPackSizes() (*swagger.PackSizesPayload, error) {
	sizes, err := ps.repo.GetPackSizes()
	if err != nil {
		return nil, errors.New("failed to retrieve pack sizes")
	}

	// Return a structured response using Swagger model
	response := &swagger.PackSizesPayload{
		PackSizes: sizes,
	}
	return response, nil
}

// UpdatePackSizes updates the available pack sizes in the repository.
func (ps *PackService) UpdatePackSizes(newSizes []int) error {
	if len(newSizes) == 0 {
		return errors.New("pack sizes cannot be empty")
	}
	for _, size := range newSizes {
		if size <= 0 {
			return fmt.Errorf("invalid pack size: %d (must be positive)", size)
		}
	}

	// Sort pack sizes in descending order to maintain consistency
	sort.Sort(sort.IntSlice(newSizes))

	// Save the new pack sizes
	return ps.repo.SetPackSizes(newSizes)
}

// DeletePackSizeHandler removes a value from the repository.
func (ps *PackService) DeletePackSizeHandler(size int) error {
	return ps.repo.DeletePackSize(size)
}

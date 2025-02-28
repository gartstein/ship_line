package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyRepo struct {
	packSizes []int
}

func (d *dummyRepo) GetPackSizes() ([]int, error) {
	return d.packSizes, nil
}

func (d *dummyRepo) SetPackSizes(sizes []int) error {
	d.packSizes = sizes
	return nil
}

func (d *dummyRepo) DeletePackSize(size int) error {
	return nil
}

func TestCalculatePacks(t *testing.T) {
	// Dummy repository with pack sizes.
	repo := &dummyRepo{packSizes: []int{250, 500, 1000, 2000, 5000}}
	ps := services.NewPackService(repo)

	t.Run("OrderZero", func(t *testing.T) {
		// Test order = 0.
		res, err := ps.CalculatePacks(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, *res.ItemsOrdered)
		assert.Equal(t, 0, *res.TotalItemsUsed)
		assert.Empty(t, *res.PacksUsed)
	})

	t.Run("OrderSmallerThanSmallest", func(t *testing.T) {
		// Test order smaller than the smallest pack (e.g., order = 200).
		res, err := ps.CalculatePacks(200)
		assert.NoError(t, err)
		assert.Equal(t, 200, *res.ItemsOrdered)
		assert.Equal(t, 250, *res.TotalItemsUsed)
		assert.Equal(t, map[string]int{"250": 1}, *res.PacksUsed)
	})

	t.Run("ExactMatch", func(t *testing.T) {
		// Test exact match: order = 500.
		res, err := ps.CalculatePacks(500)
		assert.NoError(t, err)
		assert.Equal(t, 500, *res.TotalItemsUsed)
		assert.Equal(t, map[string]int{"500": 1}, *res.PacksUsed)
	})

	t.Run("MultiplePacks", func(t *testing.T) {
		// Test order requiring multiple packs: order = 501 should yield 750 (one 500 + one 250).
		res, err := ps.CalculatePacks(501)
		assert.NoError(t, err)
		assert.Equal(t, 750, *res.TotalItemsUsed)
		expected := map[string]int{"500": 1, "250": 1}
		assert.Equal(t, expected, *res.PacksUsed)
	})

	t.Run("LargeOrder", func(t *testing.T) {
		// Test large order: order = 1000251.
		res, err := ps.CalculatePacks(1000251)
		assert.NoError(t, err)
		assert.Equal(t, 1000251, *res.ItemsOrdered)
		assert.Equal(t, 1000500, *res.TotalItemsUsed)
		expected := map[string]int{"5000": 200, "500": 1}
		assert.Equal(t, expected, *res.PacksUsed)
	})
}

// Package services provides business logic and service layer functions for the application.
// It encapsulates operations that coordinate data processing and application workflows.
package services

import (
	"fmt"
	"math"
	"ship_line/utils"
	"sort"
	"strconv"

	"ship_line/swagger"
)

// PackRepository defines the interface for retrieving and updating pack sizes.
type PackRepository interface {
	GetPackSizes() ([]int, error)
	SetPackSizes([]int) error
	DeletePackSize(size int) error
}

// PackService provides methods for calculating pack distribution.
type PackService struct {
	repo PackRepository
}

// NewPackService constructs a new PackService.
func NewPackService(repo PackRepository) *PackService {
	return &PackService{repo: repo}
}

// CalculatePacks calculates the pack distribution for a given order.
// It returns a pointer to swagger.CalcResult.
func (ps *PackService) CalculatePacks(order int) (*swagger.CalcResult, error) {
	// TODO: move to config
	const MaxOrder = 1000000000000
	if order > MaxOrder {
		return nil, fmt.Errorf("order %d exceeds maximum allowed value of %d", order, MaxOrder)
	}
	// Special case for zero order.
	if order == 0 {
		return &swagger.CalcResult{
			ItemsOrdered:   utils.Ptr(0),
			TotalItemsUsed: utils.Ptr(0),
			PacksUsed:      &map[string]int{},
		}, nil
	}

	packSizes, err := ps.repo.GetPackSizes()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}
	if len(packSizes) == 0 {
		return nil, fmt.Errorf("no pack sizes configured")
	}

	// Sort packSizes in descending order.
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))
	smallestPack := packSizes[len(packSizes)-1]
	// For orders smaller than the smallest pack, round up.
	if order < smallestPack {
		return &swagger.CalcResult{
			ItemsOrdered:   utils.Ptr(order),
			TotalItemsUsed: utils.Ptr(smallestPack),
			PacksUsed:      &map[string]int{strconv.Itoa(smallestPack): 1},
		}, nil
	}

	// Choose algorithm based on order value.
	// TODO: move to config
	const dpThreshold = 100000
	var result *swagger.CalcResult
	if order <= dpThreshold {
		result, err = calculatePacksDP(order, packSizes)
	} else {
		result, err = calculatePacksGreedy(order, packSizes)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

// calculatePacksDP implements a dynamic programming solution to determine the optimal
// pack distribution for fulfilling an order. The goal is to find the combination of pack sizes
// that results in a total (bestSum) that is at least equal to the order while using the minimum
// number of packs. It returns a CalcResult containing the ordered items, total items used, and
// a breakdown of packs used.
//
// Algorithm Explanation:
//  1. Define maxSum as order plus the smallest pack size. This provides an upper bound for our search.
//  2. Create a dp array where dp[s] represents the minimum number of packs needed to reach the sum s.
//     Initialize dp[0] to 0 (zero packs to reach zero) and all other sums to a very high value (math.MaxInt32).
//  3. Create a combination slice where each element is a map tracking the count of each pack size used to reach that sum.
//  4. Iterate over all sums from 0 to maxSum. For each sum that is reachable (dp[s] != math.MaxInt32), consider adding each pack size.
//     - For each pack size, calculate a new sum ns = s + pack size.
//     - If using one additional pack at sum s results in a lower count than what was previously recorded for ns,
//     update dp[ns] and copy the combination from s, incrementing the count for that pack size.
//  5. After building the dp and combination arrays, search for the first reachable sum (bestSum) that is >= order.
//  6. If no valid sum is found, fallback to returning the smallest pack that is >= order.
//  7. Finally, convert the combination for bestSum to the desired format and return the result.
func calculatePacksDP(order int, packSizes []int) (*swagger.CalcResult, error) {
	smallestPack := packSizes[len(packSizes)-1]
	maxSum := order + smallestPack

	// Initialize dp array with high values (indicating unreachable sums).
	dp := make([]int, maxSum+1)
	// combination[i] stores the mapping of pack sizes used to reach sum i.
	combination := make([]map[int]int, maxSum+1)
	for s := 0; s <= maxSum; s++ {
		dp[s] = math.MaxInt32
		combination[s] = nil
	}
	// Base case: zero packs are needed to achieve a sum of zero.
	dp[0] = 0
	combination[0] = make(map[int]int)

	// Build the dp array.
	for s := 0; s <= maxSum; s++ {
		if dp[s] == math.MaxInt32 {
			continue
		}
		for _, size := range packSizes {
			ns := s + size
			if ns > maxSum {
				continue
			}
			// If using an additional pack leads to a better (lower) pack count, update dp and combination.
			if dp[s]+1 < dp[ns] {
				dp[ns] = dp[s] + 1
				combination[ns] = utils.CopyMap(combination[s])
				combination[ns][size]++
			}
		}
	}

	// Find the smallest sum that is at least the order and is reachable.
	bestSum := -1
	for s := order; s <= maxSum; s++ {
		if dp[s] != math.MaxInt32 {
			bestSum = s
			break
		}
	}
	// If no combination is found, fallback to the smallest pack that is >= order.
	if bestSum == -1 {
		for _, size := range packSizes {
			if size >= order {
				return &swagger.CalcResult{
					ItemsOrdered:   utils.Ptr(order),
					TotalItemsUsed: utils.Ptr(size),
					PacksUsed:      &map[string]int{strconv.Itoa(size): 1},
				}, nil
			}
		}
		return nil, fmt.Errorf("no valid pack combination found for order %d", order)
	}
	// Convert the combination map to the expected result format.
	packsUsed := utils.ConvertMapKeys(combination[bestSum])
	return &swagger.CalcResult{
		ItemsOrdered:   utils.Ptr(order),
		TotalItemsUsed: utils.Ptr(bestSum),
		PacksUsed:      &packsUsed,
	}, nil
}

func calculatePacksGreedy(order int, packSizes []int) (*swagger.CalcResult, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes))) // Sort descending

	// Check for exact match
	for _, size := range packSizes {
		if size == order {
			return &swagger.CalcResult{
				ItemsOrdered:   utils.Ptr(order),
				TotalItemsUsed: utils.Ptr(size),
				PacksUsed:      &map[string]int{strconv.Itoa(size): 1},
			}, nil
		}
	}

	remaining := order
	packsUsed := make(map[int]int)
	totalUsed := 0

	// Greedy selection of packs
	for _, size := range packSizes {
		if remaining <= 0 {
			break
		}
		count := remaining / size
		if count > 0 {
			packsUsed[size] = count
			totalUsed += size * count
			remaining -= size * count
		}
	}

	// Handle remainder by choosing the smallest pack >= remainder
	if remaining > 0 {
		chosenPack := packSizes[len(packSizes)-1]
		for i := len(packSizes) - 1; i >= 0; i-- {
			if packSizes[i] >= remaining {
				chosenPack = packSizes[i]
				break
			}
		}
		packsUsed[chosenPack]++
		totalUsed += chosenPack
	}

	// Post-processing: Optimize by replacing smaller packs with larger ones
	for i := 0; i < len(packSizes); i++ {
		currentSize := packSizes[i]
		smallerPacks := packSizes[i+1:] // Packs smaller than currentSize

		// Calculate total items from smaller packs
		totalFromSmaller := 0
		for _, s := range smallerPacks {
			totalFromSmaller += s * packsUsed[s]
		}

		// Determine how many currentSize packs can replace the smaller ones
		replaceCount := totalFromSmaller / currentSize
		if replaceCount == 0 {
			continue
		}

		// Calculate the number of packs being replaced
		currentPacksCount := 0
		for _, s := range smallerPacks {
			currentPacksCount += packsUsed[s]
		}

		// Replace only if it reduces the number of packs
		if replaceCount < currentPacksCount {
			// Remove smaller packs
			for _, s := range smallerPacks {
				delete(packsUsed, s)
			}
			// Add the replacement packs
			packsUsed[currentSize] += replaceCount
		}
	}

	// Final check: Use a single larger pack if possible
	sortedPacks := append([]int{}, packSizes...)
	sort.Ints(sortedPacks) // Now in ascending order

	var bestSinglePack int
	found := false
	for _, size := range sortedPacks {
		if size >= order && (size <= totalUsed || !found) {
			bestSinglePack = size
			found = true
		}
	}

	if found && bestSinglePack <= totalUsed {
		packsUsed = map[int]int{bestSinglePack: 1}
		totalUsed = bestSinglePack
	}

	// Convert to result format
	result := utils.ConvertMapKeys(packsUsed)
	return &swagger.CalcResult{
		ItemsOrdered:   utils.Ptr(order),
		TotalItemsUsed: utils.Ptr(totalUsed),
		PacksUsed:      &result,
	}, nil
}

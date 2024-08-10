package gonuts

import (
	"fmt"
	"math"
)

// PaginationInfo contains information about the current pagination state
type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	FirstItem    int   `json:"first_item"`
	LastItem     int   `json:"last_item"`
	FirstPage    int   `json:"first_page"`
	LastPage     int   `json:"last_page"`
	NextPage     *int  `json:"next_page"`
	PreviousPage *int  `json:"previous_page"`
}

// NewPaginationInfo creates a new PaginationInfo instance
//
// Parameters:
//   - currentPage: the current page number
//   - perPage: the number of items per page
//   - totalItems: the total number of items in the dataset
//
// Returns:
//   - *PaginationInfo: a new instance of PaginationInfo
//
// Example usage:
//
//	pagination := gonuts.NewPaginationInfo(2, 10, 95)
//	fmt.Printf("Current Page: %d\n", pagination.CurrentPage)
//	fmt.Printf("Total Pages: %d\n", pagination.TotalPages)
//	fmt.Printf("Next Page: %v\n", *pagination.NextPage)
//
//	// Output:
//	// Current Page: 2
//	// Total Pages: 10
//	// Next Page: 3
func NewPaginationInfo(currentPage, perPage int, totalItems int64) *PaginationInfo {
	if currentPage < 1 {
		currentPage = 1
	}
	if perPage < 1 {
		perPage = 10 // Default to 10 items per page
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if currentPage > totalPages {
		currentPage = totalPages
	}

	firstItem := (currentPage-1)*perPage + 1
	lastItem := firstItem + perPage - 1
	if int64(lastItem) > totalItems {
		lastItem = int(totalItems)
	}

	var nextPage, previousPage *int
	if currentPage < totalPages {
		next := currentPage + 1
		nextPage = &next
	}
	if currentPage > 1 {
		prev := currentPage - 1
		previousPage = &prev
	}

	return &PaginationInfo{
		CurrentPage:  currentPage,
		PerPage:      perPage,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		FirstItem:    firstItem,
		LastItem:     lastItem,
		FirstPage:    1,
		LastPage:     totalPages,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}
}

// Offset calculates the offset for database queries
//
// Returns:
//   - int: the offset to use in database queries
func (p *PaginationInfo) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}

// Limit returns the number of items per page
//
// Returns:
//   - int: the number of items per page
func (p *PaginationInfo) Limit() int {
	return p.PerPage
}

// HasNextPage checks if there is a next page
//
// Returns:
//   - bool: true if there is a next page, false otherwise
func (p *PaginationInfo) HasNextPage() bool {
	return p.NextPage != nil
}

// HasPreviousPage checks if there is a previous page
//
// Returns:
//   - bool: true if there is a previous page, false otherwise
func (p *PaginationInfo) HasPreviousPage() bool {
	return p.PreviousPage != nil
}

// PageNumbers returns a slice of page numbers to display
//
// Parameters:
//   - max: the maximum number of page numbers to return
//
// Returns:
//   - []int: a slice of page numbers
//
// This method is useful for generating pagination controls in user interfaces.
// It aims to provide a balanced range of page numbers around the current page.
func (p *PaginationInfo) PageNumbers(max int) []int {
	if max >= p.TotalPages {
		pages := make([]int, p.TotalPages)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	half := max / 2
	start := p.CurrentPage - half
	end := p.CurrentPage + half

	if start < 1 {
		start = 1
		end = max
	}

	if end > p.TotalPages {
		end = p.TotalPages
		start = p.TotalPages - max + 1
		if start < 1 {
			start = 1
		}
	}

	pages := make([]int, end-start+1)
	for i := range pages {
		pages[i] = start + i
	}
	return pages
}

// String returns a string representation of the pagination info
//
// Returns:
//   - string: a string representation of the pagination info
func (p *PaginationInfo) String() string {
	return fmt.Sprintf("Page %d of %d (Total items: %d, Per page: %d)",
		p.CurrentPage, p.TotalPages, p.TotalItems, p.PerPage)
}

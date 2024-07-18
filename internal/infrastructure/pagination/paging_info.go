package pagination

import "fmt"

type PagingInfo[T any] struct {
	PageNumber   int
	ItemsPerPage int
	TotalItems   int64
	Items        []T
}

func (p PagingInfo[T]) HasNextPage() bool {
	return p.PageNumber*p.ItemsPerPage < int(p.TotalItems)
}

func (p PagingInfo[T]) HasPreviousPage() bool {
	return p.PageNumber > 1
}

func (p PagingInfo[T]) DisplayRange() string {
	if p.TotalItems < int64(p.ItemsPerPage) {
		return fmt.Sprintf("%d-%d", 1, p.TotalItems)
	}

	startIndex := ((p.PageNumber - 1) * (p.ItemsPerPage)) + 1

	endIndex := p.PageNumber * p.ItemsPerPage

	if endIndex > int(p.TotalItems) {
		endIndex = int(p.TotalItems)
	}

	return fmt.Sprintf("%d-%d", startIndex, endIndex)
}

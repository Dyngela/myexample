package types

// Usually this file is in a separate package called client that I share between every project to almost

type APIError struct {
	Message string `json:"message"`
}

type GenericResponse struct {
	Message string `json:"message"`
}

type Page[T any] struct {
	TotalRecords int64 `json:"totalRecords"`
	TotalPages   int64 `json:"totalPages"`
	CurrentPage  int   `json:"currentPage"`
	Data         []T   `json:"data"`
}

func CalculateTotalPages(totalRecords int64, pageSize int) int64 {
	return (totalRecords + int64(pageSize) - 1) / int64(pageSize)
}

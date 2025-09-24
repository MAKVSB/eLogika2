package common

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Pagination    *SearchRequestPagination `json:"pagination,omitempty"`
	Sorting       []SearchRequestSorting   `json:"sorting,omitempty"`
	ColumnFilters []SearchRequestFilter    `json:"columnFilters,omitempty"`
}

type SearchRequestPagination struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

type SearchRequestSorting struct {
	Desc bool   `json:"desc"`
	ID   string `json:"id"`
}

type SearchRequestFilter struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
}

func GetSearchParamsOrDefault(c *gin.Context, queryName string) (*SearchRequest, *ErrorResponse) {
	defaultSearchParams := SearchRequest{
		Pagination: &SearchRequestPagination{
			PageIndex: 0,
			PageSize:  25,
		},
	}
	var searchParams SearchRequest

	encoded := c.Query(queryName)
	if encoded == "" {
		return &defaultSearchParams, nil
	}

	jsonBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "Failed to decode search parameters",
			Details: err.Error(),
		}
	}

	if err := json.Unmarshal(jsonBytes, &searchParams); err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "Failed to parse search parameters",
			Details: err.Error(),
		}
	}

	if searchParams.Pagination == nil {
		searchParams.Pagination = defaultSearchParams.Pagination
	}

	return &searchParams, nil
}

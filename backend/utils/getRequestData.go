package utils

import (
	"reflect"

	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func GetRequestDataWithSearch[P any, R any](c *gin.Context, queryName string) (*common.ErrorResponse, *P, *R, *common.SearchRequest) {
	// Parse original data
	err, params, reqData := GetRequestData[P, R](c)
	if err != nil {
		return err, nil, nil, nil
	}

	var searchParams *common.SearchRequest
	searchParams, err = common.GetSearchParamsOrDefault(c, queryName)
	if err != nil {
		return err, nil, nil, nil
	}

	return nil, params, reqData, searchParams
}

func GetRequestData[P any, R any](c *gin.Context) (*common.ErrorResponse, *P, *R) {
	// Allocate params
	var params *P
	if reflect.TypeOf((*P)(nil)).Elem().Kind() != reflect.Interface {
		params = new(P) // allocate!
		if err := c.ShouldBindUri(params); err != nil {
			return &common.ErrorResponse{
				Code:    400,
				Message: "Invalid or missing URI parameters",
			}, nil, nil
		}
	}

	// Allocate request data
	var reqData *R
	if reflect.TypeOf((*R)(nil)).Elem().Kind() != reflect.Interface {
		reqData = new(R) // allocate!
		if err := c.ShouldBindJSON(reqData); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				return &common.ErrorResponse{
					Code:    422,
					Message: "Validation failed",
					Details: ve.Error(),
				}, nil, nil
			}
			return &common.ErrorResponse{
				Code:    422,
				Message: "Validation failed",
				Details: err.Error(),
			}, nil, nil
		}
	}

	return nil, params, reqData
}

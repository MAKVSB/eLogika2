package helpers

import (
	"strconv"
	"strings"

	"elogika.vsb.cz/backend/modules/common"
)

type V1Login struct {
	Username   string
	InstanceID uint
}

func ParseV1Login(s string) (*V1Login, *common.ErrorResponse) {
	parts := strings.Split(s, ";")

	if len(parts) == 1 {
		userLogin := parts[0]

		return &V1Login{
			Username:   userLogin,
			InstanceID: 0,
		}, nil
	} else if len(parts) == 2 {
		userLogin := parts[0]

		instanceId, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return nil, &common.ErrorResponse{
				Code:    400,
				Message: "Invalid user id",
			}
		}

		return &V1Login{
			Username:   userLogin,
			InstanceID: uint(instanceId),
		}, nil
	} else {
		return nil, &common.ErrorResponse{
			Code:    400,
			Message: "Invalid user id",
		}
	}
}

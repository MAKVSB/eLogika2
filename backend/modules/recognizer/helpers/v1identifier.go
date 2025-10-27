package helpers

import (
	"strconv"
	"strings"

	"elogika.vsb.cz/backend/modules/common"
)

type SheetTypeEnum string

const (
	SheetTypeTeacher SheetTypeEnum = "T"
	SheetTypeStudent SheetTypeEnum = "S"
)

type V1Identifier struct {
	CourseID   uint
	TestID     uint
	Type       SheetTypeEnum
	SheetOrder uint
}

func ParseV1Identifier(s string) (*V1Identifier, *common.ErrorResponse) {
	parts := strings.Split(s, ";")
	if len(parts) != 4 || parts[0] != "V1" {
		return nil, &common.ErrorResponse{
			Code:    400,
			Message: "Invalid identifier format",
		}
	}

	courseID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    400,
			Message: "Invalid course id",
		}
	}

	testID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    400,
			Message: "Invalid test id",
		}
	}

	// Expecting exactly 1 character type + number (e.g. "S1" or "T12")
	t := parts[3]
	if len(t) < 2 {
		return nil, &common.ErrorResponse{
			Code:    400,
			Message: "Invalid type or sheet order",
		}
	}

	typePart := t[:1]
	orderPart := t[1:]

	sheetOrder, err := strconv.ParseUint(orderPart, 10, 64)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    400,
			Message: "Invalid sheet order",
		}
	}

	return &V1Identifier{
		CourseID:   uint(courseID),
		TestID:     uint(testID),
		Type:       SheetTypeEnum(typePart),
		SheetOrder: uint(sheetOrder),
	}, nil
}

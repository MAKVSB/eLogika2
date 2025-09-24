package tokens

import (
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
)

type Token interface {
	IsRevoked() bool
	Parse(tokenStr string, allowExpired bool) *common.ErrorResponse
	Invalidate() error
	InvalidateByUser() error
	New(user dtos.LoggedUserDTO) (string, error)
}

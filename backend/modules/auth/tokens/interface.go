package tokens

import (
	"elogika.vsb.cz/backend/modules/common"
)

type Token interface {
	IsRevoked() bool
	Parse(tokenStr string, allowExpired bool) *common.ErrorResponse
	Invalidate() *common.ErrorResponse
	InvalidateByUser() *common.ErrorResponse
}

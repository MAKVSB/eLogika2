package models

import (
	"database/sql"
	"time"

	"elogika.vsb.cz/backend/modules/auth/enums"
)

type JWTTokenType string

const (
	JWTTokenTypeAccess  JWTTokenType = "ACCESS"
	JWTTokenTypeRefresh JWTTokenType = "REFRESH"
	JWTTokenTypeApi     JWTTokenType = "API"
)

// long standing tokens. (RefreshToken, ApiToken)
type AuthToken struct {
	TokenID    string                 `gorm:"primaryKey;size:36"`
	UserID     uint                   ``
	IssuedAt   time.Time              ``
	RevokedAt  sql.NullTime           ``
	RevokedFor enums.RevokedForEnum   ``
	ExpiresAt  time.Time              ``
	TokenType  enums.JWTTokenTypeEnum ``
}

func (AuthToken) TableName() string {
	return "tokens"
}

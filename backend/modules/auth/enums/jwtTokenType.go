package enums

type JWTTokenTypeEnum string

const (
	JWTTokenTypeAccess  JWTTokenTypeEnum = "ACCESS"
	JWTTokenTypeRefresh JWTTokenTypeEnum = "REFRESH"
	JWTTokenTypeApi     JWTTokenTypeEnum = "API"
)

type RevokedForEnum string

const (
	RevokedForUser  RevokedForEnum = "USER"
	RevokedForToken RevokedForEnum = "TOKEN"
)

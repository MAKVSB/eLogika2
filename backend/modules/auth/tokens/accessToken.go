package tokens

import (
	"strings"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/auth/helpers"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AccessToken struct {
	models.AuthToken
	Issuer   string
	UserData dtos.LoggedUserDTO
}

func (t *AccessToken) GetInnerToken() models.AuthToken {
	return t.AuthToken
}

func (t *AccessToken) Get(c *gin.Context, allowExpired bool) *common.ErrorResponse {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return &common.ErrorResponse{
			Message: "Missing token header",
		}
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return &common.ErrorResponse{
			Message: "Invalid token format",
		}
	}

	return t.Parse(strings.TrimPrefix(authHeader, "Bearer "), allowExpired)
}

func (t *AccessToken) IsRevoked() bool {
	return helpers.GetInmemoryRevokeStore().IsRevoked(
		t.TokenID,
		t.UserID,
		t.IssuedAt,
	)
}

func (t *AccessToken) Parse(tokenStr string, allowExpired bool) *common.ErrorResponse {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return initializers.GlobalAppConfig.ACCESS_SECRET, nil
	})

	if err != nil && err.Error() != "token has invalid claims: token is expired" && !allowExpired {
		return &common.ErrorResponse{
			Message: "Failed to parse token",
			Details: err.Error(),
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	expTime, ok := claims["exp"].(float64)
	if !ok {
		return &common.ErrorResponse{
			Message: "Invalid eat claim",
		}
	}
	t.ExpiresAt = time.Unix(int64(expTime), 0)

	iat, ok := claims["iat"].(float64)
	if !ok {
		return &common.ErrorResponse{
			Message: "Invalid iat claim",
		}
	}
	t.IssuedAt = time.Unix(int64(iat), 0)

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return &common.ErrorResponse{
			Message: "Invalid sub claim",
		}
	}
	t.UserID = uint(userIDFloat)

	jti, ok := claims["jti"].(string)
	if !ok {
		return &common.ErrorResponse{
			Message: "Invalid jti claim",
		}
	}
	t.TokenID = jti

	iss, ok := claims["iss"].(string)
	if !ok {
		return &common.ErrorResponse{
			Message: "Invalid iss claim",
		}
	}
	t.Issuer = iss

	typ, ok := claims["type"].(string)
	if !ok {
		return &common.ErrorResponse{
			Message: "Invalid token type",
		}
	}
	t.TokenType = enums.JWTTokenTypeEnum(typ)

	userRaw, ok := claims["user"]
	if !ok {
		return &common.ErrorResponse{
			Message: "User claims not found",
		}
	}

	userMap, ok := userRaw.(map[string]interface{})
	if !ok {
		return &common.ErrorResponse{
			Message: "User claims invalid",
		}
	}

	userData, err := helpers.MapToStruct(userMap)
	if err != nil {
		return &common.ErrorResponse{
			Message: "User claims cannot be parsed",
		}
	}
	t.UserData = *userData
	return nil
}

func (t *AccessToken) Invalidate() error {
	result := initializers.DB.Model(&models.AuthToken{}).
		Where("token_id = ? AND token_type = ? AND revoked_at IS NULL", t.TokenID, t.TokenType).
		Update("revoked_at", time.Now())

	helpers.GetInmemoryRevokeStore().Add(t.AuthToken)

	return result.Error
}

func (t *AccessToken) InvalidateByUser() error {
	result := initializers.DB.Model(&models.AuthToken{}).
		Where("user_id = ? AND token_type = ? AND revoked_at IS NULL", t.UserID, t.TokenType).
		Update("revoked_at", time.Now())

	helpers.GetInmemoryRevokeStore().Add(t.AuthToken)

	return result.Error
}

func (t *AccessToken) New(user dtos.LoggedUserDTO) (string, error) {
	timeFreeze := time.Now()

	t.TokenID = uuid.New().String()
	t.UserID = user.ID
	t.Issuer = "core.api.elogika.vsb.cz"
	t.IssuedAt = timeFreeze
	t.ExpiresAt = timeFreeze.Add(initializers.GlobalAppConfig.ACCESS_LENGTH)
	t.TokenType = enums.JWTTokenTypeAccess

	claims := jwt.MapClaims{
		"iss":  t.Issuer,
		"sub":  t.UserID,
		"exp":  t.ExpiresAt.Unix(),
		"iat":  t.IssuedAt.Unix(),
		"jti":  t.TokenID,
		"type": t.TokenType,
		"user": user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(initializers.GlobalAppConfig.ACCESS_SECRET)
}

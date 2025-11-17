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
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type ApiToken struct {
	models.AuthToken
	Issuer   string
	UserData dtos.LoggedUserDTO
}

func (t *ApiToken) GetInnerToken() models.AuthToken {
	return t.AuthToken
}

func (t *ApiToken) IsRevoked() bool {
	var dbToken models.AuthToken
	if err := initializers.DB.
		Where("token_id = ? AND token_type = ? AND revoked_at IS NULL AND expires_at > ?", t.TokenID, t.TokenType, time.Now()).
		First(&dbToken).Error; err != nil {
		return true
	}

	return false
}

func (t *ApiToken) Parse(tokenStr string, allowExpired bool) *common.ErrorResponse {
	token, err := jwt.Parse(strings.TrimPrefix(tokenStr, "api_"), func(token *jwt.Token) (interface{}, error) {
		return initializers.GlobalAppConfig.API_SECRET, nil
	})
	if err != nil {
		return &common.ErrorResponse{
			Code:    401,
			Message: "Failed to parse api token",
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

func (t *ApiToken) Invalidate() *common.ErrorResponse {
	if err := initializers.DB.
		Model(&models.AuthToken{}).
		Where("token_id = ?", t.TokenID).
		Update("revoked_at", time.Now()).
		Update("revoked_for", enums.RevokedForToken).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to invalidate token",
		}
	}

	return nil
}

func (t *ApiToken) InvalidateByUser() *common.ErrorResponse {
	if err := initializers.DB.
		Model(&models.AuthToken{}).
		Where("user_id = ?, token_type = ?", t.UserID, t.TokenType).
		Update("revoked_at", time.Now()).
		Update("revoked_for", enums.RevokedForUser).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to invalidate token",
		}
	}

	return nil
}

func (t *ApiToken) New(user dtos.LoggedUserDTO, name string, expiresAt time.Time) (string, error) {
	timeFreeze := time.Now()

	t.TokenID = uuid.New().String()
	t.UserID = user.ID
	t.Issuer = "core.api.elogika.vsb.cz"
	t.IssuedAt = timeFreeze
	t.ExpiresAt = expiresAt
	t.TokenType = enums.JWTTokenTypeApi
	t.Name = name

	if err := initializers.DB.Create(&t.AuthToken).Error; err != nil {
		return "", err
	}

	initializers.DB.Create(&t.AuthToken)

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
	tokenString, err := token.SignedString(initializers.GlobalAppConfig.API_SECRET)
	return "api_" + tokenString, err
}

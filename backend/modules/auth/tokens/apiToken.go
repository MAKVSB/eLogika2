package tokens

import (
	"errors"
	"fmt"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		fmt.Println("database error")
		return true
	}

	return false
}

func (t *ApiToken) Parse(tokenStr string, allowExpired bool) *common.ErrorResponse {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return initializers.GlobalAppConfig.API_SECRET, nil
	})
	if err != nil {
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

	return nil
}

func (t *ApiToken) Invalidate() error {
	result := initializers.DB.Model(&models.AuthToken{}).
		Where("token_id = ? AND token_type = ? AND revoked_at IS NULL", t.TokenID, t.TokenType).
		Update("revoked_at", time.Now())

	return result.Error
}

func (t *ApiToken) InvalidateByUser() error {
	result := initializers.DB.Model(&models.AuthToken{}).
		Where("user_id = ? AND token_type = ? AND revoked_at IS NULL", t.UserID, t.TokenType).
		Update("revoked_at", time.Now())

	return result.Error
}

func (t *ApiToken) New(user dtos.LoggedUserDTO) (string, error) {
	timeFreeze := time.Now()

	t.TokenID = uuid.New().String()
	t.UserID = user.ID
	t.Issuer = "core.api.elogika.vsb.cz"
	t.IssuedAt = timeFreeze
	t.ExpiresAt = timeFreeze.Add(initializers.GlobalAppConfig.API_LENGTH)
	t.TokenType = enums.JWTTokenTypeAccess

	initializers.DB.Create(&t.AuthToken)

	claims := jwt.MapClaims{
		"iss":  t.Issuer,
		"sub":  t.UserID,
		"exp":  t.ExpiresAt.Unix(),
		"iat":  t.IssuedAt.Unix(),
		"jti":  t.TokenID,
		"type": t.TokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(initializers.GlobalAppConfig.API_SECRET)
}

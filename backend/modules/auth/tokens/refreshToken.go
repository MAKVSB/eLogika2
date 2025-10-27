package tokens

import (
	"strings"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type RefreshToken struct {
	models.AuthToken
	Issuer string
}

func (t *RefreshToken) GetInnerToken() models.AuthToken {
	return t.AuthToken
}

func (t *RefreshToken) Get(c *gin.Context, allowExpired bool) *common.ErrorResponse {
	refreshTokenStr, err := c.Cookie("refresh_token")
	if err != nil {
		return &common.ErrorResponse{
			Message: "Refresh token missing",
		}
	}

	return t.Parse(strings.TrimPrefix(refreshTokenStr, "rt_"), allowExpired)
}

func (t *RefreshToken) IsRevoked() bool {
	var dbToken models.AuthToken
	if err := initializers.DB.
		Where("token_id = ? AND token_type = ? AND revoked_at IS NULL AND expires_at > ?", t.TokenID, t.TokenType, time.Now()).
		First(&dbToken).Error; err != nil {
		return true
	}

	return false
}

func (t *RefreshToken) Parse(tokenStr string, allowExpired bool) *common.ErrorResponse {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return initializers.GlobalAppConfig.REFRESH_SECRET, nil
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

func (t *RefreshToken) Invalidate() *common.ErrorResponse {
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

func (t *RefreshToken) InvalidateByUser() *common.ErrorResponse {
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

func (t *RefreshToken) New(user dtos.LoggedUserDTO) (string, error) {
	timeFreeze := time.Now()

	t.TokenID = uuid.New().String()
	t.UserID = user.ID
	t.Issuer = "core.api.elogika.vsb.cz"
	t.IssuedAt = timeFreeze
	t.ExpiresAt = timeFreeze.Add(initializers.GlobalAppConfig.REFRESH_LENGTH)
	t.TokenType = enums.JWTTokenTypeRefresh

	if err := initializers.DB.Create(&t.AuthToken).Error; err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iss":  t.Issuer,
		"sub":  t.UserID,
		"exp":  t.ExpiresAt.Unix(),
		"iat":  t.IssuedAt.Unix(),
		"jti":  t.TokenID,
		"type": t.TokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(initializers.GlobalAppConfig.REFRESH_SECRET)
	return "rt_" + tokenString, err
}

package crons

import (
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/enums"
)

func DeleteExpiredExpirations() {
	initializers.DB.
		Where("token_type = ?", enums.JWTTokenTypeAccess).
		Where("expires_at < ?", time.Now().AddDate(0, 0, -1)).
		Delete(&models.AuthToken{})

	err := initializers.DB.Exec("ALTER INDEX ALL ON dbo.tokens REBUILD WITH (FILLFACTOR = 90, ONLINE = ON);").Error
	if err != nil {
		panic(err)
	}
}

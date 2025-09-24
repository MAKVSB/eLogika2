package helpers

import (
	"log"
	"sync"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/enums"
	"gorm.io/gorm"
)

type RevokedTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]time.Time
	users  map[uint]time.Time
}

var (
	store *RevokedTokenStore
	once  sync.Once
)

// GetInmemoryRevokeStore returns the singleton instance
func GetInmemoryRevokeStore() *RevokedTokenStore {
	once.Do(func() {
		store = &RevokedTokenStore{
			tokens: make(map[string]time.Time),
			users:  make(map[uint]time.Time),
		}
	})
	return store
}

func (r *RevokedTokenStore) Add(token models.AuthToken) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if token.RevokedFor == enums.RevokedForUser {
		r.users[token.UserID] = token.RevokedAt.Time
	} else {
		r.tokens[token.TokenID] = token.RevokedAt.Time
	}
}

func (r *RevokedTokenStore) IsRevoked(tokenId string, userID uint, iat time.Time) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, exists := r.tokens[tokenId]
	if exists && iat.Before(value) {
		return true
	}

	value, exists = r.users[userID]
	if exists && iat.Before(value) {
		return true
	}
	return false
}

func (r *RevokedTokenStore) RemoveExpired() {
	tokenTime := initializers.GlobalAppConfig.ACCESS_LENGTH

	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	for token, revokedAt := range r.tokens {
		if (revokedAt.Add(tokenTime)).Before(now) {
			delete(r.tokens, token)
		}
	}
	for user, revokedAt := range r.users {
		if (revokedAt.Add(tokenTime)).Before(now) {
			delete(r.users, user)
		}
	}
}

func StartRevokedTokenSync(db *gorm.DB, interval time.Duration) {
	var lastChecked time.Time
	store := GetInmemoryRevokeStore()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C

			var revoked []models.AuthToken
			if err := db.Where("revoked_at > ? AND token_type = ?", lastChecked, enums.JWTTokenTypeAccess).Find(&revoked).Error; err != nil {
				log.Printf("sync error: %v", err)
				continue
			}

			if len(revoked) > 0 {
				lastChecked = time.Now()
			}

			for _, t := range revoked {
				store.Add(t)
			}

			store.RemoveExpired()
		}
	}()
}

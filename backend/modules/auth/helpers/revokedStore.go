package helpers

import (
	"log"
	"sync"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/enums"
	"gorm.io/gorm"
)

type RevokedTokenRecord struct {
	RevokedAt time.Time
	ExpiresAt time.Time
}

type RevokedTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]RevokedTokenRecord
	users  map[uint]RevokedTokenRecord
}

var (
	store *RevokedTokenStore
	once  sync.Once
)

// GetInmemoryRevokeStore returns the singleton instance
func GetInmemoryRevokeStore() *RevokedTokenStore {
	once.Do(func() {
		store = &RevokedTokenStore{
			tokens: make(map[string]RevokedTokenRecord),
			users:  make(map[uint]RevokedTokenRecord),
		}
	})
	return store
}

func (r *RevokedTokenStore) Add(token models.AuthToken) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if token.RevokedFor == enums.RevokedForUser {
		val, exists := r.users[token.UserID]
		if !exists || token.RevokedAt.Time.After(val.RevokedAt) {
			r.users[token.UserID] = RevokedTokenRecord{
				RevokedAt: token.RevokedAt.Time,
				ExpiresAt: token.ExpiresAt,
			}
		}
	} else {
		val, exists := r.tokens[token.TokenID]
		if !exists || token.RevokedAt.Time.After(val.RevokedAt) {
			r.tokens[token.TokenID] = RevokedTokenRecord{
				RevokedAt: token.RevokedAt.Time,
				ExpiresAt: token.ExpiresAt,
			}
		}
	}
}

func (r *RevokedTokenStore) IsRevoked(tokenId string, userID uint, iat time.Time) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	revokedToken, exists := r.tokens[tokenId]
	if exists && iat.Before(revokedToken.RevokedAt) {
		return true
	}

	revokedToken, exists = r.users[userID]
	if exists && iat.Before(revokedToken.RevokedAt) {
		return true
	}
	return false
}

func (r *RevokedTokenStore) RemoveExpired() {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	for token, revokedToken := range r.tokens {
		if (revokedToken.ExpiresAt).Before(now) {
			delete(r.tokens, token)
		}
	}
	for user, revokedToken := range r.users {
		if (revokedToken.ExpiresAt).Before(now) {
			delete(r.users, user)
		}
	}
}

func StartRevokedTokenSync(db *gorm.DB, interval time.Duration) {
	var lastChecked time.Time
	store := GetInmemoryRevokeStore()

	{
		// On startup load immediately
		var revoked []models.AuthToken
		if err := db.Where("revoked_at > ? AND token_type = ?", lastChecked, enums.JWTTokenTypeAccess).Find(&revoked).Error; err != nil {
			log.Printf("sync error: %v", err)
		}

		if len(revoked) > 0 {
			lastChecked = time.Now()
		}

		for _, t := range revoked {
			store.Add(t)
		}

		store.RemoveExpired()
	}

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

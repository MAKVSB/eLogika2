package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type UserListItemDTO struct {
	ID               uint                       `json:"id"`
	FirstName        string                     `json:"firstName"`
	FamilyName       string                     `json:"familyName"`
	Username         string                     `json:"username"`
	Email            string                     `json:"email"`
	Type             enums.UserTypeEnum         `json:"type"`
	IdentityProvider enums.IdentityProviderEnum `json:"identityProvider"`
}

func (m UserListItemDTO) From(d *models.User) UserListItemDTO {
	dto := UserListItemDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,
		Type:       d.Type,
	}

	return dto
}

package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type UserAPiTokenDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (m UserAPiTokenDTO) From(d *models.AuthToken) UserAPiTokenDTO {
	return UserAPiTokenDTO{
		ID:        d.TokenID,
		Name:      d.Name,
		IssuedAt:  d.IssuedAt,
		ExpiresAt: d.ExpiresAt,
	}
}

type UserDTO struct {
	ID               uint                       `json:"id"`
	Version          uint                       `json:"version"`
	DegreeBefore     string                     `json:"degreeBefore"`
	FirstName        string                     `json:"firstName"`
	FamilyName       string                     `json:"familyName"`
	DegreeAfter      string                     `json:"degreeAfter"`
	Username         string                     `json:"username"`
	Email            string                     `json:"email"`
	Notification     UserNotificationDTO        `json:"notification"`
	Type             enums.UserTypeEnum         `json:"type"`
	IdentityProvider enums.IdentityProviderEnum `json:"identityProvider"`
	Courses          []UserCourseDTO            `json:"courses"`
	ApiTokens        []UserAPiTokenDTO          `json:"apiTokens"`
}

func (m UserDTO) From(d *models.User) UserDTO {
	dto := UserDTO{
		ID:               d.ID,
		DegreeBefore:     d.DegreeBefore,
		FirstName:        d.FirstName,
		FamilyName:       d.FamilyName,
		DegreeAfter:      d.DegreeAfter,
		Username:         d.Username,
		Email:            d.Email,
		Type:             d.Type,
		Notification:     UserNotificationDTO{}.From(&d.Notification),
		IdentityProvider: d.IdentityProvider,
		Version:          d.Version,
		Courses:          make([]UserCourseDTO, len(d.UserCourses)),
		ApiTokens:        make([]UserAPiTokenDTO, len(d.ApiTokens)),
	}

	for i, course := range d.UserCourses {
		dto.Courses[i] = UserCourseDTO{}.From(course)
	}

	for i, token := range d.ApiTokens {
		dto.ApiTokens[i] = UserAPiTokenDTO{}.From(token)
	}

	return dto
}

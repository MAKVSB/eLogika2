package seed

import (
	"fmt"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/helpers"
	"elogika.vsb.cz/backend/modules/common/enums"
)

func CreateUsers() []models.User {
	var users []models.User

	pass1, err := helpers.HashPassword("testing")
	if err != nil {
		fmt.Println("Failed to hash password", err)
	}

	users = append(users, models.User{
		Username:     "MAK0065",
		DegreeBefore: "Bc.",
		FirstName:    "Daniel",
		FamilyName:   "Makovský",
		DegreeAfter:  "",
		Password:     pass1,
		Email:        "mak0065@vsb.cz",
		Notification: models.UserNotification{
			Discord: models.NotificationDiscord{
				UserID: "794923933991305228",
				Level: models.NotificationLevel{
					Results:  true,
					Messages: true,
					Terms:    true,
				},
			},
			Email: models.NotificationEmail{
				Level: models.NotificationLevel{
					Results:  true,
					Messages: true,
					Terms:    true,
				},
			},
			Push: models.NotificationPush{
				Level: models.NotificationLevel{
					Results:  true,
					Messages: true,
					Terms:    true,
				},
				Token: "",
			},
		},
		Version:          1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Type:             enums.UserTypeAdmin,
		IdentityProvider: enums.IdentityProviderInternal,
		// IdentityProviderID: json.RawMessage{},
	})
	users = append(users, models.User{
		Username:   "MAK0065",
		FirstName:  "Daniel",
		FamilyName: "Makovský",
		Password:   pass1,
		Email:      "mak0065.student@vsb.cz",
		Notification: models.UserNotification{
			Discord: models.NotificationDiscord{
				UserID: "794923933991305228",
				Level: models.NotificationLevel{
					Results:  true,
					Messages: true,
					Terms:    true,
				},
			},
			Email: models.NotificationEmail{
				Level: models.NotificationLevel{
					Results:  true,
					Messages: true,
					Terms:    true,
				},
			},
			Push: models.NotificationPush{
				Level: models.NotificationLevel{
					Results:  true,
					Messages: true,
					Terms:    true,
				},
				Token: "",
			},
		},
		Version:          1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Type:             enums.UserTypeNormal,
		IdentityProvider: enums.IdentityProviderInternal,
		// IdentityProviderID: json.RawMessage{},
	})

	for index, user := range users {
		if err := initializers.DB.Create(&user).Error; err != nil {
			fmt.Println("Failed to insert", index, err)
		}
	}

	return users
}

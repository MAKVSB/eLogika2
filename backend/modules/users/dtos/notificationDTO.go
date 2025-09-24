package dtos

import "elogika.vsb.cz/backend/models"

type NotificationLevelDTO struct {
	Results  bool `json:"results"`
	Messages bool `json:"messages"`
	Terms    bool `json:"terms"`
}

func (m NotificationLevelDTO) From(d *models.NotificationLevel) NotificationLevelDTO {
	dto := NotificationLevelDTO{
		Results:  d.Results,
		Messages: d.Messages,
		Terms:    d.Terms,
	}

	return dto
}

type NotificationDiscordDTO struct {
	Level  NotificationLevelDTO `json:"level"`
	UserID string               `json:"userId"`
}

func (m NotificationDiscordDTO) From(d *models.NotificationDiscord) NotificationDiscordDTO {
	dto := NotificationDiscordDTO{
		Level:  NotificationLevelDTO{}.From(&d.Level),
		UserID: d.UserID,
	}

	return dto
}

type NotificationEmailDTO struct {
	Level NotificationLevelDTO `json:"level"`
}

func (m NotificationEmailDTO) From(d *models.NotificationEmail) NotificationEmailDTO {
	dto := NotificationEmailDTO{
		Level: NotificationLevelDTO{}.From(&d.Level),
	}

	return dto
}

type NotificationPushDTO struct {
	Level NotificationLevelDTO `json:"level"`
}

func (m NotificationPushDTO) From(d *models.NotificationPush) NotificationPushDTO {
	dto := NotificationPushDTO{
		Level: NotificationLevelDTO{}.From(&d.Level),
	}

	return dto
}

type UserNotificationDTO struct {
	Discord NotificationDiscordDTO `json:"discord"`
	Email   NotificationEmailDTO   `json:"email"`
	Push    NotificationPushDTO    `json:"push"`
}

func (m UserNotificationDTO) From(d *models.UserNotification) UserNotificationDTO {
	dto := UserNotificationDTO{
		Discord: NotificationDiscordDTO{}.From(&d.Discord),
		Email:   NotificationEmailDTO{}.From(&d.Email),
		Push:    NotificationPushDTO{}.From(&d.Push),
	}

	return dto
}

package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type FileDTO struct {
	ID           uint   `json:"id"`
	OriginalName string `json:"originalName"`
	StoredName   string `json:"storedName"`
	MIMEType     string `json:"mimeType"`
	SizeBytes    int64  `json:"sizeBytes"`
}

func (m FileDTO) From(d *models.File) FileDTO {
	dto := FileDTO{
		ID:           d.ID,
		OriginalName: d.OriginalName,
		StoredName:   d.StoredName,
		MIMEType:     d.MIMEType,
		SizeBytes:    d.SizeBytes,
	}

	return dto
}

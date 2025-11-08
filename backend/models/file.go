package models

import "time"

type File struct {
	CommonModel
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      ``
	OriginalName string    ``
	StoredName   string    ``
	MIMEType     string    ``
	SizeBytes    int64     ``
	UploadedAt   time.Time ``

	Chapters []*File `gorm:"many2many:chapter_files;"`
}

func (File) TableName() string {
	return "files"
}

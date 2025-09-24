package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type TermDTO struct {
	ID           uint      `json:"id"`
	Version      uint      `json:"version"`
	CourseID     uint      `json:"courseId"`
	CourseItemID uint      `json:"courseItemId"`
	Name         string    `json:"name"`
	ActiveFrom   time.Time `json:"activeFrom"`
	ActiveTo     time.Time `json:"activeTo"`
	RequiresSign bool      `json:"requiresSign"`
	SignInFrom   time.Time `json:"signInFrom"`
	SignInTo     time.Time `json:"signInTo"`
	SignOutFrom  time.Time `json:"signOutFrom"`
	SignOutTo    time.Time `json:"signOutTo"`
	// OfflineTo:    n,
	Classroom      string `json:"classroom"`
	StudentsMax    uint   `json:"studentsMax"`
	StudentsJoined uint   `json:"studentsJoined"`
	Tries          uint   `json:"tries"`
}

func (m TermDTO) From(d *models.Term) TermDTO {
	dto := TermDTO{
		ID:      d.ID,
		Version: d.Version,

		CourseID:     d.CourseID,
		CourseItemID: d.CourseItemID,
		Name:         d.Name,
		ActiveFrom:   d.ActiveFrom,
		ActiveTo:     d.ActiveTo,
		RequiresSign: d.RequiresSign,
		SignInFrom:   d.SignInFrom,
		SignInTo:     d.SignInTo,
		SignOutFrom:  d.SignOutFrom,
		SignOutTo:    d.SignOutTo,
		// OfflineTo:    d.OfflineTo,
		Classroom:      d.Classroom,
		StudentsMax:    d.StudentsMax,
		StudentsJoined: uint(len(d.Students)),
		Tries:          d.Tries,
	}

	return dto
}

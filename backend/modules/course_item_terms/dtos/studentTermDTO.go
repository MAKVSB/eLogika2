package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type StudentTermDTO struct {
	ID                  uint      `json:"id"`
	Version             uint      `json:"version"`
	CourseID            uint      `json:"courseId"`
	CourseItemID        uint      `json:"courseItemId"`
	CourseItemName      string    `json:"courseItemName"`
	CourseItemGroupName string    `json:"courseItemGroupName"`
	Name                string    `json:"name"`
	ActiveFrom          time.Time `json:"activeFrom"`
	ActiveTo            time.Time `json:"activeTo"`
	RequiresSign        bool      `json:"requiresSign"`
	SignInFrom          time.Time `json:"signInFrom"`
	SignInTo            time.Time `json:"signInTo"`
	SignOutFrom         time.Time `json:"signOutFrom"`
	SignOutTo           time.Time `json:"signOutTo"`
	// OfflineTo:    n,
	Classroom   string `json:"classroom"`
	StudentsMax uint   `json:"studentsMax"`
	Tries       uint   `json:"tries"`

	StudentsJoined int  `json:"studentsJoined"`
	CanJoin        bool `json:"canJoin"`
	CanLeave       bool `json:"canLeave"`
	WillSignOut    bool `json:"willSignOut"`
	Joined         bool `json:"joined"`
}

func (m StudentTermDTO) From(
	d models.Term,
	joinedCount int,
	joined bool,
	canJoin bool,
	canLeave bool,
	willSignOut bool,
) StudentTermDTO {
	dto := StudentTermDTO{
		ID:      d.ID,
		Version: d.Version,

		CourseID:            d.CourseID,
		CourseItemID:        d.CourseItemID,
		CourseItemName:      d.CourseItem.Name,
		CourseItemGroupName: groupName(d),
		Name:                d.Name,
		ActiveFrom:          d.ActiveFrom,
		ActiveTo:            d.ActiveTo,
		RequiresSign:        d.RequiresSign,
		SignInFrom:          d.SignInFrom,
		SignInTo:            d.SignInTo,
		SignOutFrom:         d.SignOutFrom,
		SignOutTo:           d.SignOutTo,
		// OfflineTo:    d.OfflineTo,
		Classroom:   d.Classroom,
		StudentsMax: d.StudentsMax,
		Tries:       d.Tries,

		StudentsJoined: joinedCount,
		CanJoin:        canJoin,
		CanLeave:       canLeave,
		WillSignOut:    willSignOut,
		Joined:         joined,
	}

	return dto
}

func groupName(t models.Term) string {
	groupString := ""

	courseItem := t.CourseItem.Parent
	for {
		if courseItem != nil {
			groupString = courseItem.Name + groupString

			if courseItem.Parent != nil {
				groupString = "→" + groupString
				courseItem = courseItem.Parent
			} else {
				courseItem = nil
			}
		} else {
			return groupString
		}
	}
}

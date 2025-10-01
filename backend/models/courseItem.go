package models

import (
	"fmt"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type CourseItem struct {
	CommonModel
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      ``
	CreatedById uint           ``
	CreatedBy   *User          ``
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt ``
	Version     uint           ``

	Name              string                      ``
	Type              enums.CourseItemTypeEnum    ``
	PointsMin         uint                        ``
	PointsMax         uint                        ``
	Mandatory         bool                        ``
	StudyForm         enums.StudyFormEnum         ``
	MaxAttempts       uint                        ``
	AllowNegative     bool                        ``
	ManagedBy         enums.CourseUserRoleEnum    ``
	EvaluateByAttempt enums.EvaluateByAttemptEnum `` // TODO TODO TODO add to user interface

	ActivityDetailID *uint ``
	TestDetailID     *uint ``
	GroupDetailID    *uint ``
	ParentID         *uint ``
	CourseID         uint  ``

	ActivityDetail *CourseItemActivity ``
	TestDetail     *CourseItemTest     ``
	GroupDetail    *CourseItemGroup    ``
	Children       []*CourseItem       `gorm:"foreignKey:ParentID"`
	Terms          []*Term             ``
	Parent         *CourseItem         ``

	// Temp values for querying

	Editable bool `gorm:"-"`
}

func (ci *CourseItem) RecursiveLoadParent(tx *gorm.DB) error {
	return ci.recursiveLoadParent(tx, nil)
}

func (ci *CourseItem) recursiveLoadParent(tx *gorm.DB, cacheItems *[]*CourseItem) error {
	if tx == nil {
		tx = initializers.DB
	}

	if cacheItems == nil {
		temp := []*CourseItem{}
		if err := tx.Where("course_id = ?", ci.CourseID).Find(&temp).Error; err != nil {
			return fmt.Errorf("failed to load course items: %w", err)
		}
		cacheItems = &temp
	}

	if ci.ParentID == nil {
		return nil
	}

	for _, ciInner := range *cacheItems {
		if ciInner.ID == *ci.ParentID {
			ci.Parent = ciInner
			return ciInner.recursiveLoadParent(tx, cacheItems)
		}
	}

	return nil
}

func (ci *CourseItem) RecursiveLoadChildren(tx *gorm.DB) error {
	return ci.recursiveLoadChildren(tx, nil)
}

func (ci *CourseItem) recursiveLoadChildren(tx *gorm.DB, cacheItems *[]*CourseItem) error {
	if tx == nil {
		tx = initializers.DB
	}

	if cacheItems == nil {
		temp := []*CourseItem{}
		if err := tx.Where("course_id = ?", ci.CourseID).Find(&temp).Error; err != nil {
			return fmt.Errorf("failed to load course items: %w", err)
		}
		cacheItems = &temp
	}

	ci.Children = make([]*CourseItem, 0)
	for _, ciInner := range *cacheItems {
		if ciInner.ParentID != nil && *ciInner.ParentID == ci.ID {
			if err := ciInner.recursiveLoadChildren(tx, cacheItems); err != nil {
				return err
			}
			ci.Children = append(ci.Children, ciInner)
			ciInner.Parent = ci
		}
	}

	return nil
}

func (CourseItem) TableName() string {
	return "course_items"
}

func (CourseItem) ApplyFilters(query *gorm.DB, filters []common.SearchRequestFilter, model any, extra map[string]interface{}) (*gorm.DB, *common.ErrorResponse) {
	if filters != nil {
		// 1) Handle special cases and build a new slice without them
		var remainingFilters []common.SearchRequestFilter
		for _, filter := range filters {
			if filter.ID == "ParentID" {
				if filter.Value == "NULL" {
					query = query.Where("parent_id IS NULL")
				} else {
					query = query.Where("parent_id = ?", filter.Value)
				}
			} else {
				remainingFilters = append(remainingFilters, filter)
			}
		}

		query, err := CommonModel{}.ApplyFilters(query, remainingFilters, CourseItem{}, nil, "")
		if err != nil {
			return query, err
		}
	}
	return query, nil
}

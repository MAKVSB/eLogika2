package crons

import (
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
)

func ExpireReadyTests() {
	var readyTestInstances []*models.TestInstance
	initializers.DB.
		InnerJoins("Term", initializers.DB.Select("").Where("active_to < ?", time.Now().Add(5*time.Minute))).
		Preload("CourseItem").
		Where("state = ?", enums.TestInstanceStateReady).
		Where("form = ?", enums.TestInstanceFormOnline).
		Find(&readyTestInstances)

	for _, testInstance := range readyTestInstances {
		transaction := initializers.DB.Begin()

		testInstance.State = enums.TestInstanceStateExpired
		rootCoureItem := testInstance.CourseItem.ID
		if testInstance.CourseItem.ParentID != nil {
			rootCoureItem = *testInstance.CourseItem.ParentID
		}

		services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		err := services_course_item.UpdateSelectedResults(transaction, testInstance.CourseItem.CourseID, rootCoureItem, testInstance.ParticipantID)
		if err != nil {
			transaction.Rollback()
			continue
		}

		if err := transaction.Save(&testInstance).Error; err != nil {
			transaction.Rollback()
			continue
		}

		if err := transaction.Commit().Error; err != nil {
			transaction.Rollback()
			continue
		}

	}
}

package crons

import (
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/handlers"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
)

func FinishActiveTests() {
	var activeEndedTestInstances []*models.TestInstance
	initializers.DB.
		Preload("CourseItem").
		Where("state = ?", enums.TestInstanceStateActive).
		Where("ends_at < ?", time.Now().Add(time.Minute)).
		Find(&activeEndedTestInstances)

	for _, testInstance := range activeEndedTestInstances {
		transaction := initializers.DB.Begin()

		testInstance.State = enums.TestInstanceStateFinished
		err := handlers.EvaluateTestInstance(transaction, testInstance.ID, nil)
		if err != nil {
			transaction.Rollback()
			continue
		}

		rootCoureItem := testInstance.CourseItem.ID
		if testInstance.CourseItem.ParentID != nil {
			rootCoureItem = *testInstance.CourseItem.ParentID
		}

		services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		err = services_course_item.UpdateSelectedResults(transaction, testInstance.CourseItem.CourseID, rootCoureItem, testInstance.ParticipantID)
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

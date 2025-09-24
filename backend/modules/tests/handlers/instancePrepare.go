package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstancePrepareRequest struct {
	TermID       uint `json:"termId" uri:"termId" binding:"required"`             // TODO check that user has permission for this term
	CourseItemID uint `json:"courseItemId" uri:"courseItemId" binding:"required"` // TODO check that user has permission for this term
}

type TestInstancePrepareResponse struct {
	InstanceID uint `json:"instanceId"`
}

// @Summary Starts test instance for user
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestGeneratorRequest true "Ability to filter results"
// @Success 200 {object} TestGeneratorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/instance/start [post]
func TestInstancePrepare(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		TestInstancePrepareRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	// check permissions
	coursePermissions := auth.GetClaimCourse(userData.Courses, params.CourseID)
	if coursePermissions == nil || !coursePermissions.IsStudent() {
		c.JSON(403, common.ErrorResponse{
			Message: "Not enough permissions",
		})
		return
	}

	transaction := initializers.DB.Begin()

	// TODO check if student has joined this term
	var term models.Term
	if err := transaction.
		Joins("JOIN user_terms ON user_terms.term_id = terms.id AND user_terms.user_id = ?", userData.ID).
		First(&term, reqData.TermID).Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(404, common.ErrorResponse{
			Message: "Failed to fetch term",
			Details: err.Error(),
		})
		return
	}

	courseItemQuery := `
		WITH course_tree AS (
			-- start with the course_item linked to the term
			SELECT ci.*
			FROM course_items ci
			INNER JOIN terms t ON ci.id = t.course_item_id
			WHERE t.id = ?  -- replace with the term id
			UNION ALL
			-- recursively select children
			SELECT ci.*
			FROM course_items ci
			INNER JOIN course_tree ct ON ci.parent_id = ct.id
		)
		SELECT *
		FROM course_tree WHERE id = ?;
	`

	var courseItem models.CourseItem
	if err := transaction.
		Raw(courseItemQuery, reqData.TermID, reqData.CourseItemID).
		Preload("TestDetail").
		First(&courseItem).Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(404, common.ErrorResponse{
			Message: "Failed to fetch term",
			Details: err.Error(),
		})
		return
	}

	// TODO Check user still has attempts left on term/courseItem/childs/parents/whatever
	// TODO lockování tabulky: 		Table(models.UserTerm{}.TableName()+" WITH (XLOCK, TABLOCK)").

	// Generate test
	generatedTest, err := GenerateSingleTest(
		transaction,
		params.CourseID,
		courseItem.TestDetail.TestTemplateID,
		term.ID,
		&userData,
		userData.Username,
		reqData.CourseItemID,
	)

	if err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	testInstance, err := helpers.CreateInstance(
		transaction,
		generatedTest,
		userData.ID,
		term.ID,
		courseItem.ID,
		enums.TestInstanceFormOnline,
	)
	if err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO check that user has permission
	// TODO check that no test is running

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to commit changes",
		})
		return
	}

	c.JSON(200, TestInstancePrepareResponse{
		InstanceID: testInstance.ID,
	})
}

package handlers

import (
	"fmt"
	"strings"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/recognizer/dtos"
	"elogika.vsb.cz/backend/modules/recognizer/helpers"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecognizerTestInstanceGetResponse struct {
	TestData dtos.TestDTO `json:"testData"`
}

// @Summary Gets the expected format of the answer sheet
// @Tags Recognizer
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param identifier path int true "Test identifier code"
// @Success 200 {object} TestInstanceGetResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/recognizer/test/{identifier} [GET]
func RecognizerTestGet(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			Identifier string `uri:"identifier" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	identifierParts := strings.Split(params.Identifier, ";")
	if strings.ToUpper(identifierParts[0]) == "V1" {
		testIdentifierData, err := helpers.ParseV1Identifier(params.Identifier)
		if err != nil {
			return err
		}

		// Get test data
		var testData *models.Test
		if err := initializers.DB.
			Preload("Questions", func(db *gorm.DB) *gorm.DB {
				tmp := db.
					Joins("Question").
					Preload("Answers").
					Order("\"order\" ASC").
					Limit(18).
					Offset(18 * int(testIdentifierData.SheetOrder))

				switch testIdentifierData.Type {
				case helpers.SheetTypeStudent:
					tmp = tmp.Where("Question.question_format = ?", enums.QuestionFormatTest)
				case helpers.SheetTypeTeacher:
					tmp = tmp.Where("Question.question_format = ?", enums.QuestionFormatOpen)
				default:
					panic(fmt.Sprintf("unexpected helpers.SheetTypeEnum: %#v", testIdentifierData.Type))
				}

				return tmp
			}).
			Find(&testData, testIdentifierData.TestID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
				Details: err.Error(),
			}
		}

		courseItemService := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		_, err = courseItemService.GetCourseItemByID(initializers.DB, testData.CourseID, testData.CourseItemID, userData.ID, userRole, nil, false, nil)
		if err != nil {
			return err
		}

		c.JSON(200, RecognizerTestInstanceGetResponse{
			TestData: dtos.TestDTO{}.From(testData),
		})
	}

	return nil
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/recognizer/helpers"
	testHandlers "elogika.vsb.cz/backend/modules/tests/handlers"
	testHelpers "elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecognizerTestSaveRequestQuestion struct {
	Order   uint   `json:"order" binding:"required"`
	Answers []bool `json:"answers" binding:"required"`
}

type RecognizerTestSaveRequest struct {
	FileName  string                              `json:"fileName" binding:"required"`
	Login     string                              `json:"login" binding:"required"`
	Variant   string                              `json:"variant" binding:"required"`
	Questions []RecognizerTestSaveRequestQuestion `json:"questions" binding:"required"`
}

type RecognizerTestSaveResponse struct {
	Success bool `json:"success"`
}

// @Summary Saves test instance for user by tutor (advanced options)
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestInstanceTutorSaveRequest true "Ability to filter results"
// @Success 200 {object} TestInstanceTutorSaveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/recognizer/test [POST]
func RecognizerTestSave(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// TODO validate from here

	// 1. Get the JSON metadata string from a form field
	metadataStr := c.PostForm("data")
	if metadataStr == "" {
		return &common.ErrorResponse{
			Code:    400,
			Message: "Missing test recognizer data",
		}
	}

	var reqData RecognizerTestSaveRequest
	if err := json.Unmarshal([]byte(metadataStr), &reqData); err != nil {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Validation failed",
			Details: err.Error(),
		}
	}

	identifierParts := strings.Split(reqData.Variant, ";")
	if strings.ToUpper(identifierParts[0]) == "V1" {
		testIdentifierData, err := helpers.ParseV1Identifier(reqData.Variant)
		if err != nil {
			return err
		}

		loginData, err := helpers.ParseV1Login(reqData.Login)
		if err != nil {
			return err
		}

		// Get course item data
		var testData *models.Test
		if err := initializers.DB.
			InnerJoins("CourseItem").
			Find(&testData).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find associated course item",
				Details: err.Error(),
			}
		}

		// Check role validity
		if err := auth.GetClaimCourseRole(userData, testIdentifierData.CourseID, userRole); err != nil {
			return err
		}
		var courseItem *models.CourseItem
		// Check if tutor/garant can view/modify courseItem
		switch userRole {
		case enums.CourseUserRoleAdmin:
			if err := initializers.DB.
				Preload("TestDetail").
				Preload("Parent").
				Find(&courseItem, testData.CourseItemID).Error; err != nil {
				return &common.ErrorResponse{
					Code:    403,
					Message: "Not enough permission for this item",
				}
			}
		case enums.CourseUserRoleGarant:
			if err := initializers.DB.
				Preload("TestDetail").
				Preload("Parent").
				Where("managed_by = ?", enums.CourseUserRoleGarant).
				Find(&courseItem, testData.CourseItemID).Error; err != nil {
				return &common.ErrorResponse{
					Code:    403,
					Message: "Not enough permission for this item",
				}
			}
		case enums.CourseUserRoleTutor:
			if err := initializers.DB.
				Preload("TestDetail").
				Preload("Parent").
				Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userData.ID).
				Find(&courseItem, testData.CourseItemID).Error; err != nil {
				return &common.ErrorResponse{
					Code:    403,
					Message: "Not enough permission for this item",
				}
			}
		default:
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		// Get test data
		instanceId, err := CreateOrVerifyInstance(initializers.DB, testIdentifierData.CourseID, testIdentifierData.TestID, loginData.Username, loginData.InstanceID)
		if err != nil {
			return err
		}

		transaction := initializers.DB.Begin()
		events := make([]*models.TestInstanceEvent, 0)

		instanceData, err := GetQuestionsByInstanceID(transaction, instanceId, testIdentifierData.Type, testIdentifierData.SheetOrder)
		if err != nil {
			transaction.Rollback()
			return err
		}

		if len(instanceData.Questions) != len(reqData.Questions) {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    0,
				Message: "Failed to save recognized data",
				Details: "Question count does not match",
			}
		}

		for rq_i, rq := range reqData.Questions {
			iq := instanceData.Questions[rq_i]

			switch iq.TestQuestion.Question.QuestionFormat {
			case enums.QuestionFormatOpen:
				iq.TextAnswerPercentage = CalculatePercentage(rq.Answers)
				iq.TextAnswerReviewedByID = &userData.ID

				eventData, _ := json.Marshal(map[string]interface{}{
					"QuestionOrder": iq.TestQuestion.Order,
				})

				events = append(events, &models.TestInstanceEvent{
					TestInstanceID: iq.TestInstanceID,
					UserID:         userData.ID,
					OccuredAt:      time.Now(),
					ReceivedAt:     time.Now(),
					EventSource:    enums.TestInstanceEventSourceServer,
					EventType:      enums.TestInstanceEventTypeQuestionUpdate,
					EventData:      eventData,
				})

				if err := transaction.Save(&iq).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to save answers for question",
						Details: err.Error(),
					}
				}
			case enums.QuestionFormatTest:
				if len(rq.Answers) != len(iq.Answers) {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    0,
						Message: "Failed to save recognized data",
						Details: "Answer box count does not match",
					}
				}

				for checked_i, checked := range rq.Answers {
					iq.Answers[checked_i].Selected = checked

					eventData, _ := json.Marshal(map[string]interface{}{
						"QuestionOrder": iq.TestQuestion.Order,
						"AnswerOrder":   rq.Order,
						"AnswerData":    checked,
					})

					events = append(events, &models.TestInstanceEvent{
						TestInstanceID: instanceData.ID,
						UserID:         userData.ID,
						OccuredAt:      time.Now(),
						ReceivedAt:     time.Now(),
						EventSource:    enums.TestInstanceEventSourceServer,
						EventType:      enums.TestInstanceEventTypeQuestionUpdate,
						EventData:      eventData,
					})

					if err := transaction.Save(&iq.Answers[checked_i]).Error; err != nil {
						transaction.Rollback()
						return &common.ErrorResponse{
							Code:    500,
							Message: "Failed to save answers for question",
							Details: err.Error(),
						}
					}
				}

			default:
				panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", iq.TestQuestion.Question.QuestionFormat))
			}
		}

		switch instanceData.State {
		case enums.TestInstanceStateReady:
			instanceData.StartedAt = time.Now()
			instanceData.EndsAt = time.Now()
			instanceData.EndedAt = time.Now()
			instanceData.State = enums.TestInstanceStateFinished
		case enums.TestInstanceStateActive:
			instanceData.EndedAt = time.Now()
			instanceData.State = enums.TestInstanceStateFinished
		}

		if err := transaction.Save(&instanceData).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to save instance",
				Details: err.Error(),
			}
		}

		err = testHandlers.EvaluateTestInstance(transaction, instanceData.ID, &userData, false)
		if err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to evaluate test",
				Details: err,
			}
		}

		rootCoureItem := testData.CourseItem.ID
		if testData.CourseItem.ParentID != nil {
			rootCoureItem = *testData.CourseItem.ParentID
		}

		services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		err = services_course_item.UpdateSelectedResults(transaction, testIdentifierData.CourseID, rootCoureItem, instanceData.ParticipantID)
		if err != nil {
			transaction.Rollback()
			return err
		}

		if len(events) != 0 {
			if err := transaction.Save(&events).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to save events for question",
					Details: err.Error(),
				}
			}
		}

		// Append parser file
		file, header, file_err := c.Request.FormFile("file")
		if file_err == nil {
			defer file.Close()

			ext := strings.ToLower(filepath.Ext(header.Filename))
			if ext == "" || len(ext) > 10 {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    400,
					Message: "Invalid file extension",
				}
			}

			// Generate a random UUID filename
			newFileName, err := utils.GenerateFileName(transaction, ext)
			if err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: err.Error(),
				}
			}

			// Create file
			out, err2 := os.Create(filepath.Join(initializers.GlobalAppConfig.UPLOADS_DESTINATION, newFileName))
			if err2 != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Could not save file",
				}
			}
			defer out.Close()

			// Copy data
			size, err2 := io.Copy(out, file)
			if err2 != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to save file",
				}
			}

			if err := transaction.
				Model(&models.RecognizerFile{}).
				Where("test_instance_id = ?", instanceData.ID).
				Where("unique_ident = ?", reqData.Variant).
				Delete(&models.RecognizerFile{}).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to remove old recognizer file",
					Details: err.Error(),
				}
			}

			recoFile := models.RecognizerFile{
				File: &models.File{
					UserID:       userData.ID,
					OriginalName: header.Filename,
					StoredName:   newFileName,
					MIMEType:     header.Header.Get("Content-Type"),
					SizeBytes:    size,
					UploadedAt:   time.Now(),
				},
				UniqueIdent:    reqData.Variant,
				TestInstanceID: instanceData.ID,
			}

			// Save to database for later linking
			if err := transaction.Create(&recoFile).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "DB save failed",
				}
			}
		}

		if err := transaction.Commit().Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to commit changes",
				Details: err.Error(),
			}
		}

		c.JSON(200, RecognizerTestSaveResponse{
			Success: true,
		})
		return nil
	}

	return &common.ErrorResponse{
		Code:    500,
		Message: "Failed to save answer data",
		Details: "QR codes version not supported",
	}
}

func CreateOrVerifyInstance(dbRef *gorm.DB, courseId uint, testId uint, username string, instanceId uint) (uint, *common.ErrorResponse) {
	var testInstance *models.TestInstance
	if instanceId != 0 {
		// Verify instance exists with said user
		if err := dbRef.
			Select("test_instances.id").
			Where("test_id = ?", testId).
			InnerJoins("Participant", initializers.DB.Select("Participant.id").Where("Participant.username = ?", username)).
			Find(&testInstance, instanceId).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find test insstance",
				Details: err.Error(),
			}
		}
	} else {
		// Get the id of the instance of sspecific user
		if err := dbRef.
			Select("test_instances.id").
			Where("test_id = ?", testId).
			InnerJoins("Participant", initializers.DB.Where("Participant.username = ?", username)).
			Find(&testInstance).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find test insstance",
				Details: err.Error(),
			}
		}
	}

	if testInstance.ID != 0 {
		// Test instance found and verified
		return testInstance.ID, nil
	} else {
		// Instance not found so create new instance

		// Get user id
		var courseUserId uint
		if err := dbRef.
			Model(models.CourseUser{}).
			Where("course_id = ?", courseId).
			InnerJoins("User", initializers.DB.Where("\"User\".\"username\" = ?", username)).
			Pluck("User.id", &courseUserId).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    404,
				Message: "Failed to fetch user data",
				Details: err.Error(),
			}
		}

		if courseUserId == 0 {
			return 0, &common.ErrorResponse{
				Code:    404,
				Message: "User not found in course",
			}
		}

		// Get test data
		var test *models.Test
		if err := dbRef.
			Where("course_id = ?", courseId).
			// Where("course_item_id = ?", params.CourseItemID).
			Preload("Questions").
			Preload("Questions.Answers").
			First(&test, testId).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    404,
				Message: "Failed to fetch term",
				Details: err.Error(),
			}
		}

		testInstance, err := testHelpers.CreateInstance(
			dbRef,
			test,
			courseUserId,
			test.TermID,
			test.CourseItemID,
			enums.TestInstanceFormOffline,
		)
		if err != nil {
			return 0, err
		}

		return testInstance.ID, nil
	}
}

func GetQuestionsByInstanceID(dbRef *gorm.DB, testInstanceID uint, Type helpers.SheetTypeEnum, SheetOrder uint) (*models.TestInstance, *common.ErrorResponse) {
	var instanceData *models.TestInstance
	if err := dbRef.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			tmp := db.
				InnerJoins("TestQuestion").
				InnerJoins("TestQuestion.Question", initializers.DB.Select("id", "QuestionFormat")).
				Preload("Answers", func(db *gorm.DB) *gorm.DB {
					return db.InnerJoins("TestQuestionAnswer").Order("\"TestQuestionAnswer\".\"order\" ASC")
				}).
				Order("\"TestQuestion\".\"order\" ASC").
				Limit(18).
				Offset(18 * int(SheetOrder))

			switch Type {
			case helpers.SheetTypeStudent:
				tmp = tmp.Where("\"TestQuestion__Question\".\"question_format\" = ?", enums.QuestionFormatTest)
			case helpers.SheetTypeTeacher:
				tmp = tmp.Where("\"TestQuestion__Question\".\"question_format\" = ?", enums.QuestionFormatOpen)
			default:
				panic(fmt.Sprintf("unexpected helpers.SheetTypeEnum: %#v", Type))
			}

			return tmp
		}).
		Find(&instanceData, testInstanceID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permission for this item",
			Details: err.Error(),
		}
	}
	return instanceData, nil
}

func CalculatePercentage(checkedBoxes []bool) float64 {
	if len(checkedBoxes) != 11 {
		return 0
	}

	checkedCount := float64(0)
	totalValue := float64(0)
	for valueMultiplier, checked := range checkedBoxes {
		if checked {
			totalValue += float64(valueMultiplier) * 10
			checkedCount++
		}
	}

	if checkedCount == 0 {
		return 0
	}

	return totalValue / checkedCount
}

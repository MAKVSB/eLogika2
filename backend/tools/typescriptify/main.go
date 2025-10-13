package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	authDtos "elogika.vsb.cz/backend/modules/auth/dtos"
	authHandlers "elogika.vsb.cz/backend/modules/auth/handlers"
	"elogika.vsb.cz/backend/modules/course_items/dtos"

	"elogika.vsb.cz/backend/modules/common/enums"

	categoryHandlers "elogika.vsb.cz/backend/modules/categories/handlers"

	chapterHandlers "elogika.vsb.cz/backend/modules/chapters/handlers"

	courseHandlers "elogika.vsb.cz/backend/modules/courses/handlers"

	questionHandlers "elogika.vsb.cz/backend/modules/questions/handlers"

	courseItemsHandlers "elogika.vsb.cz/backend/modules/course_items/handlers"

	templateHandlers "elogika.vsb.cz/backend/modules/templates/handlers"

	userHandlers "elogika.vsb.cz/backend/modules/users/handlers"

	termsHandlers "elogika.vsb.cz/backend/modules/course_item_terms/handlers"

	testHandlers "elogika.vsb.cz/backend/modules/tests/handlers"
	testHelpers "elogika.vsb.cz/backend/modules/tests/helpers"

	classHandlers "elogika.vsb.cz/backend/modules/classes/handlers"

	printHandlers "elogika.vsb.cz/backend/modules/print/handlers"

	activityHandlers "elogika.vsb.cz/backend/modules/activities/handlers"

	fileHandlers "elogika.vsb.cz/backend/modules/files/handlers"

	"github.com/hypersequent/zen"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	frontendPath := "./../frontend"

	converter := typescriptify.New()
	converter.CreateInterface = true
	converter.BackupDir = ""

	converter.AddImport("import type { JSONContent } from '@tiptap/core';")
	converter.AddImport("import type { ZonedDateTime } from '@internationalized/date';")
	converter.AddImport("export type StringDate = string")

	converter.ManageType(time.Time{}, typescriptify.TypeOptions{TSType: "StringDate"})

	converter.Add(questionHandlers.QuestionListResponse{}).
		Add(questionHandlers.QuestionListRequest{}).
		Add(questionHandlers.QuestionInsertRequest{}).
		Add(questionHandlers.QuestionInsertResponse{}).
		Add(questionHandlers.QuestionUpdateRequest{}).
		Add(questionHandlers.QuestionUpdateResponse{}).
		Add(questionHandlers.QuestionToggleActiveResponse{}).
		Add(questionHandlers.QuestionCheckResponse{}).
		Add(questionHandlers.QuestionGetByIdResponse{}).
		Add(authHandlers.LoginRequest{}).
		Add(authHandlers.LoginResponse{}).
		Add(authHandlers.LogoutResponse{}).
		Add(authHandlers.RefreshResponse{}).
		Add(authDtos.LoggedUserDTO{}).
		Add(chapterHandlers.ChapterListRequest{}).
		Add(chapterHandlers.ChapterListResponse{}).
		Add(chapterHandlers.ChapterListResponse{}).
		Add(chapterHandlers.ChapterGetByIdResponse{}).
		Add(chapterHandlers.ChapterInsertRequest{}).
		Add(chapterHandlers.ChapterInsertResponse{}).
		Add(chapterHandlers.ChapterMoveResponse{}).
		Add(chapterHandlers.ChapterUpdateRequest{}).
		Add(chapterHandlers.ChapterUpdateResponse{}).
		Add(courseHandlers.CourseListRequest{}).
		Add(courseHandlers.CourseListResponse{}).
		Add(courseHandlers.CourseGetByIdResponse{}).
		Add(courseHandlers.CourseInsertRequest{}).
		Add(courseHandlers.CourseInsertResponse{}).
		Add(courseHandlers.CourseUpdateRequest{}).
		Add(courseHandlers.CourseUpdateResponse{}).
		Add(courseHandlers.ListCourseUsersResponse{}).
		Add(courseHandlers.UserCourseListResponse{}).
		Add(courseHandlers.AddCourseUserRequest{}).
		Add(courseHandlers.AddCourseUserResponse{}).
		Add(courseHandlers.RemoveCourseUserRequest{}).
		Add(courseHandlers.RemoveCourseUserResponse{}).
		Add(userHandlers.UserListRequest{}).
		Add(userHandlers.UserListResponse{}).
		Add(userHandlers.UserInsertRequest{}).
		Add(userHandlers.UserInsertResponse{}).
		Add(userHandlers.UserUpdateRequest{}).
		Add(userHandlers.UserUpdateResponse{}).
		Add(userHandlers.UserGetByIdResponse{}).
		Add(userHandlers.UserChangePassRequest{}).
		Add(categoryHandlers.CategoryListRequest{}).
		Add(categoryHandlers.CategoryListResponse{}).
		Add(categoryHandlers.CategoryInsertRequest{}).
		Add(categoryHandlers.CategoryInsertResponse{}).
		Add(categoryHandlers.CategoryUpdateRequest{}).
		Add(categoryHandlers.CategoryUpdateResponse{}).
		Add(categoryHandlers.CategoryGetByIdResponse{}).
		Add(templateHandlers.TemplateCreatorResponse{}).
		Add(templateHandlers.TemplateInsertRequest{}).
		Add(templateHandlers.TemplateInsertResponse{}).
		Add(templateHandlers.TemplateUpdateRequest{}).
		Add(templateHandlers.TemplateUpdateResponse{}).
		Add(templateHandlers.TemplateGetByIdResponse{}).
		Add(templateHandlers.TemplateListRequest{}).
		Add(templateHandlers.TemplateListResponse{}).
		Add(courseItemsHandlers.CourseItemInsertRequest{}).
		Add(courseItemsHandlers.CourseItemInsertResponse{}).
		Add(courseItemsHandlers.CourseItemUpdateRequest{}).
		Add(courseItemsHandlers.CourseItemUpdateResponse{}).
		Add(courseItemsHandlers.CourseItemListRequest{}).
		Add(courseItemsHandlers.CourseItemListResponse{}).
		Add(courseItemsHandlers.CourseItemGetByIdResponse{}).
		Add(courseItemsHandlers.StudentCourseItemListResponse{}).
		Add(courseItemsHandlers.CourseItemListResultsResponse{}).
		Add(courseItemsHandlers.CourseItemSelectResultResponse{}).
		Add(termsHandlers.TermsInsertRequest{}).
		Add(termsHandlers.TermsInsertResponse{}).
		Add(termsHandlers.TermsUpdateRequest{}).
		Add(termsHandlers.TermsUpdateResponse{}).
		Add(termsHandlers.TermsListRequest{}).
		Add(termsHandlers.TermsListResponse{}).
		Add(termsHandlers.TermsGetByIdResponse{}).
		Add(termsHandlers.StudentTermsListResponse{}).
		Add(termsHandlers.TermsJoinResponse{}).
		Add(termsHandlers.TermsLeaveResponse{}).
		Add(termsHandlers.TermsListRecursiveResponse{}).
		Add(termsHandlers.ListJoinedStudentsResponse{}).
		Add(testHandlers.ListAvailableTestsResponse{}).
		Add(testHandlers.TestInstancePrepareRequest{}).
		Add(testHandlers.TestInstancePrepareResponse{}).
		Add(testHandlers.TestInstanceStartResponse{}).
		Add(testHandlers.TestInstanceGetResponse{}).
		Add(testHandlers.TestInstanceSaveResponse{}).
		Add(testHandlers.TestListResponse{}).
		Add(testHandlers.TestInstanceListResponse{}).
		Add(testHandlers.TestInstanceGetTelemetryResponse{}).
		Add(testHelpers.TestInstanceQuestion{}).
		Add(testHandlers.TestGeneratorRequest{}).
		Add(testHandlers.TestGeneratorResponse{}).
		Add(testHandlers.TestInstanceCreateRequest{}).
		Add(testHandlers.TestInstanceCreateResponse{}).
		Add(testHandlers.TestInstanceTutorGetResponse{}).
		Add(testHandlers.TestEvaluationRequest{}).
		Add(testHandlers.TestEvaluationResponse{}).
		Add(classHandlers.ClassUpdateRequest{}).
		Add(classHandlers.ClassUpdateResponse{}).
		Add(classHandlers.ClassInsertRequest{}).
		Add(classHandlers.ClassInsertResponse{}).
		Add(classHandlers.ClassListRequest{}).
		Add(classHandlers.ClassListResponse{}).
		Add(classHandlers.ClassGetByIdResponse{}).
		Add(classHandlers.AddStudentRequest{}).
		Add(classHandlers.AddStudentResponse{}).
		Add(classHandlers.ListStudentResponse{}).
		Add(classHandlers.RemoveStudentRequest{}).
		Add(classHandlers.RemoveStudentResponse{}).
		Add(classHandlers.AddTutorRequest{}).
		Add(classHandlers.AddTutorResponse{}).
		Add(classHandlers.ListTutorResponse{}).
		Add(classHandlers.RemoveTutorRequest{}).
		Add(classHandlers.RemoveTutorResponse{}).
		Add(classHandlers.ClassImportStudentsResponse{}).
		Add(printHandlers.PrintTestRequest{}).
		Add(activityHandlers.ListAvailableActivitiesResponse{}).
		Add(activityHandlers.ActivityInstanceGetResponse{}).
		Add(activityHandlers.ActivityInstanceSaveRequest{}).
		Add(activityHandlers.ActivityInstanceSaveResponse{}).
		Add(activityHandlers.ActivityListResponse{}).
		Add(activityHandlers.ActivityInstanceSaveRequest{}).
		Add(activityHandlers.ActivityInstanceSaveResponse{}).
		Add(fileHandlers.FileUploadResponse{})

	// TODO: maybe remove once handlers exists
	converter.Add(dtos.CourseItemDTO{})

	converter.AddEnum(enums.QuestionTypeEnumAll).
		AddEnum(enums.QuestionFormatEnumAll).
		AddEnum(enums.CourseUserRoleEnumAll).
		AddEnum(enums.SemesterEnumAll).
		AddEnum(enums.MoveDirectionEnumAll).
		AddEnum(enums.IdentityProviderEnumAll).
		AddEnum(enums.UserTypeEnumAll).
		AddEnum(enums.StudyFormEnumAll).
		AddEnum(enums.CourseItemTypeEnumAll).
		AddEnum(enums.StepSelectionEnumAll).
		AddEnum(enums.CategoryFilterEnumAll).
		AddEnum(enums.AnswerDistributionEnumAll).
		AddEnum(enums.TestInstanceStateEnumAll).
		AddEnum(enums.QuestionCheckedByFilterEnumAll).
		AddEnum(enums.TestInstanceEventTypeEnumAll).
		AddEnum(enums.ClassTypeEnumAll).
		AddEnum(enums.WeekDayEnumAll).
		AddEnum(enums.WeekParityEnumAll).
		AddEnum(enums.TestInstanceFormEnumAll).
		AddEnum(enums.EvaluateByAttemptEnumAll)

	err := converter.ConvertToFile(frontendPath + "/src/lib/api_types.ts")
	if err != nil {
		panic(err.Error())
	}

	customTagHandlers := map[string]zen.CustomFn{
		"optional": func(c *zen.Converter, t reflect.Type, validate string, indent int) string {
			return ".optional()"
		},
	}
	customTypeHandlers := map[string]zen.CustomFn{
		"encoding/json.RawMessage": func(c *zen.Converter, t reflect.Type, v string, indent int) string {
			return "z.any()"
		},
	}

	c := zen.NewConverterWithOpts(
		zen.WithCustomTypes(customTypeHandlers),
		zen.WithCustomTags(customTagHandlers),
	)
	c.AddType(questionHandlers.QuestionInsertRequest{})
	c.AddType(questionHandlers.QuestionUpdateRequest{})

	c.AddType(authHandlers.LoginRequest{})

	c.AddType(chapterHandlers.ChapterInsertRequest{})
	c.AddType(chapterHandlers.ChapterUpdateRequest{})

	c.AddType(categoryHandlers.CategoryInsertRequest{})
	c.AddType(categoryHandlers.CategoryUpdateRequest{})

	c.AddType(courseHandlers.CourseInsertRequest{})
	c.AddType(courseHandlers.CourseUpdateRequest{})

	c.AddType(courseItemsHandlers.CourseItemInsertRequest{})
	c.AddType(courseItemsHandlers.CourseItemUpdateRequest{})

	c.AddType(userHandlers.UserInsertRequest{})
	c.AddType(userHandlers.UserUpdateRequest{})

	c.AddType(categoryHandlers.CategoryListResponse{})

	c.AddType(authHandlers.LoginRequest{})

	c.AddType(templateHandlers.TemplateInsertRequest{})
	c.AddType(templateHandlers.TemplateUpdateRequest{})

	c.AddType(termsHandlers.TermsInsertRequest{})
	c.AddType(termsHandlers.TermsUpdateRequest{})
	c.AddType(termsHandlers.TermsListRequest{})

	c.AddType(classHandlers.ClassUpdateRequest{})
	c.AddType(classHandlers.ClassInsertRequest{})
	c.AddType(classHandlers.ClassListRequest{})

	c.AddType(activityHandlers.ActivityInstanceSaveRequest{})

	print := "import z from \"zod/v4\"; \n"
	print += "import { en, cs } from \"zod/v4/locales\"; \n"
	print += "import { getLocale } from '$lib/paraglide/runtime'; \n"
	print += "let locale = getLocale(); \n"
	print += "if (locale === \"cs\") {; \n"
	print += "	z.config(cs()); \n"
	print += "} else {; \n"
	print += "	z.config(en()); \n"
	print += "}; \n\n"

	print += c.Export()

	err = os.WriteFile(frontendPath+"/src/lib/schemas.ts", []byte(print), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

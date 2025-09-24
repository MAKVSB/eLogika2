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

	converter.Add(questionHandlers.QuestionListResponse{})
	converter.Add(questionHandlers.QuestionListRequest{})
	converter.Add(questionHandlers.QuestionInsertRequest{})
	converter.Add(questionHandlers.QuestionInsertResponse{})
	converter.Add(questionHandlers.QuestionUpdateRequest{})
	converter.Add(questionHandlers.QuestionUpdateResponse{})
	converter.Add(questionHandlers.QuestionToggleActiveResponse{})

	converter.Add(questionHandlers.QuestionCheckResponse{})
	converter.Add(questionHandlers.QuestionGetByIdResponse{})

	converter.Add(authHandlers.LoginRequest{})
	converter.Add(authHandlers.LoginResponse{})
	converter.Add(authHandlers.LogoutResponse{})
	converter.Add(authHandlers.RefreshResponse{})
	converter.Add(authDtos.LoggedUserDTO{})

	converter.Add(chapterHandlers.ChapterListRequest{})
	converter.Add(chapterHandlers.ChapterListResponse{})
	converter.Add(chapterHandlers.ChapterListResponse{})
	converter.Add(chapterHandlers.ChapterGetByIdResponse{})
	converter.Add(chapterHandlers.ChapterInsertRequest{})
	converter.Add(chapterHandlers.ChapterInsertResponse{})
	converter.Add(chapterHandlers.ChapterMoveResponse{})
	converter.Add(chapterHandlers.ChapterUpdateRequest{})
	converter.Add(chapterHandlers.ChapterUpdateResponse{})

	converter.Add(courseHandlers.CourseListRequest{})
	converter.Add(courseHandlers.CourseListResponse{})
	converter.Add(courseHandlers.CourseGetByIdResponse{})
	converter.Add(courseHandlers.CourseInsertRequest{})
	converter.Add(courseHandlers.CourseInsertResponse{})
	converter.Add(courseHandlers.CourseUpdateRequest{})
	converter.Add(courseHandlers.CourseUpdateResponse{})
	converter.Add(courseHandlers.ListCourseUsersResponse{})
	converter.Add(courseHandlers.UserCourseListResponse{})
	converter.Add(courseHandlers.AddCourseUserRequest{})
	converter.Add(courseHandlers.AddCourseUserResponse{})
	converter.Add(courseHandlers.RemoveCourseUserRequest{})
	converter.Add(courseHandlers.RemoveCourseUserResponse{})

	converter.Add(userHandlers.UserListRequest{})
	converter.Add(userHandlers.UserListResponse{})
	converter.Add(userHandlers.UserInsertRequest{})
	converter.Add(userHandlers.UserInsertResponse{})
	converter.Add(userHandlers.UserUpdateRequest{})
	converter.Add(userHandlers.UserUpdateResponse{})
	converter.Add(userHandlers.UserGetByIdResponse{})

	converter.Add(categoryHandlers.CategoryListRequest{})
	converter.Add(categoryHandlers.CategoryListResponse{})
	converter.Add(categoryHandlers.CategoryInsertRequest{})
	converter.Add(categoryHandlers.CategoryInsertResponse{})
	converter.Add(categoryHandlers.CategoryUpdateRequest{})
	converter.Add(categoryHandlers.CategoryUpdateResponse{})
	converter.Add(categoryHandlers.CategoryGetByIdResponse{})

	converter.Add(templateHandlers.TemplateCreatorResponse{})
	converter.Add(templateHandlers.TemplateInsertRequest{})
	converter.Add(templateHandlers.TemplateInsertResponse{})
	converter.Add(templateHandlers.TemplateUpdateRequest{})
	converter.Add(templateHandlers.TemplateUpdateResponse{})
	converter.Add(templateHandlers.TemplateGetByIdResponse{})
	converter.Add(templateHandlers.TemplateListRequest{})
	converter.Add(templateHandlers.TemplateListResponse{})

	converter.Add(courseItemsHandlers.CourseItemInsertRequest{})
	converter.Add(courseItemsHandlers.CourseItemInsertResponse{})
	converter.Add(courseItemsHandlers.CourseItemUpdateRequest{})
	converter.Add(courseItemsHandlers.CourseItemUpdateResponse{})
	converter.Add(courseItemsHandlers.CourseItemListRequest{})
	converter.Add(courseItemsHandlers.CourseItemListResponse{})
	converter.Add(courseItemsHandlers.CourseItemGetByIdResponse{})
	converter.Add(courseItemsHandlers.StudentCourseItemListResponse{})

	converter.Add(termsHandlers.TermsInsertRequest{})
	converter.Add(termsHandlers.TermsInsertResponse{})
	converter.Add(termsHandlers.TermsUpdateRequest{})
	converter.Add(termsHandlers.TermsUpdateResponse{})
	converter.Add(termsHandlers.TermsListRequest{})
	converter.Add(termsHandlers.TermsListResponse{})
	converter.Add(termsHandlers.TermsGetByIdResponse{})
	converter.Add(termsHandlers.StudentTermsListResponse{})
	converter.Add(termsHandlers.TermsJoinResponse{})
	converter.Add(termsHandlers.TermsLeaveResponse{})
	converter.Add(termsHandlers.TermsListRecursiveResponse{})
	converter.Add(termsHandlers.ListJoinedStudentsResponse{})

	converter.Add(testHandlers.ListAvailableTestsResponse{})
	converter.Add(testHandlers.TestInstancePrepareRequest{})
	converter.Add(testHandlers.TestInstancePrepareResponse{})
	converter.Add(testHandlers.TestInstanceStartResponse{})
	converter.Add(testHandlers.TestInstanceGetResponse{})
	converter.Add(testHandlers.TestInstanceSaveResponse{})
	converter.Add(testHandlers.TestListResponse{})
	converter.Add(testHandlers.TestInstanceListResponse{})
	converter.Add(testHandlers.TestInstanceGetTelemetryResponse{})
	converter.Add(testHelpers.TestInstanceQuestion{})
	converter.Add(testHandlers.TestGeneratorRequest{})
	converter.Add(testHandlers.TestGeneratorResponse{})
	converter.Add(testHandlers.TestInstanceCreateRequest{})
	converter.Add(testHandlers.TestInstanceCreateResponse{})
	converter.Add(testHandlers.TestInstanceTutorGetResponse{})

	converter.Add(classHandlers.ClassUpdateRequest{})
	converter.Add(classHandlers.ClassUpdateResponse{})
	converter.Add(classHandlers.ClassInsertRequest{})
	converter.Add(classHandlers.ClassInsertResponse{})
	converter.Add(classHandlers.ClassListRequest{})
	converter.Add(classHandlers.ClassListResponse{})
	converter.Add(classHandlers.ClassGetByIdResponse{})
	converter.Add(classHandlers.AddStudentRequest{})
	converter.Add(classHandlers.AddStudentResponse{})
	converter.Add(classHandlers.ListStudentResponse{})
	converter.Add(classHandlers.RemoveStudentRequest{})
	converter.Add(classHandlers.RemoveStudentResponse{})
	converter.Add(classHandlers.AddTutorRequest{})
	converter.Add(classHandlers.AddTutorResponse{})
	converter.Add(classHandlers.ListTutorResponse{})
	converter.Add(classHandlers.RemoveTutorRequest{})
	converter.Add(classHandlers.RemoveTutorResponse{})
	converter.Add(classHandlers.ClassImportStudentsResponse{})

	converter.Add(printHandlers.PrintTestRequest{})

	converter.Add(activityHandlers.ListAvailableActivitiesResponse{})
	converter.Add(activityHandlers.ActivityInstanceGetResponse{})
	converter.Add(activityHandlers.ActivityInstanceSaveRequest{})
	converter.Add(activityHandlers.ActivityInstanceSaveResponse{})
	converter.Add(activityHandlers.ActivityListResponse{})
	converter.Add(activityHandlers.ActivityInstanceSaveRequest{})
	converter.Add(activityHandlers.ActivityInstanceSaveResponse{})

	converter.Add(fileHandlers.FileUploadResponse{})

	// TODO: maybe remove once handlers exists
	converter.Add(dtos.CourseItemDTO{})

	converter.AddEnum(enums.QuestionTypeEnumAll)
	converter.AddEnum(enums.QuestionFormatEnumAll)
	converter.AddEnum(enums.CourseUserRoleEnumAll)
	converter.AddEnum(enums.SemesterEnumAll)
	converter.AddEnum(enums.MoveDirectionEnumAll)
	converter.AddEnum(enums.IdentityProviderEnumAll)
	converter.AddEnum(enums.UserTypeEnumAll)
	converter.AddEnum(enums.StudyFormEnumAll)
	converter.AddEnum(enums.CourseItemTypeEnumAll)
	converter.AddEnum(enums.StepSelectionEnumAll)
	converter.AddEnum(enums.CategoryFilterEnumAll)
	converter.AddEnum(enums.AnswerDistributionEnumAll)
	converter.AddEnum(enums.TestInstanceStateEnumAll)
	converter.AddEnum(enums.QuestionCheckedByFilterEnumAll)
	converter.AddEnum(enums.TestInstanceEventTypeEnumAll)
	converter.AddEnum(enums.ClassTypeEnumAll)
	converter.AddEnum(enums.WeekDayEnumAll)
	converter.AddEnum(enums.WeekParityEnumAll)
	converter.AddEnum(enums.TestInstanceFormEnumAll)
	converter.AddEnum(enums.EvaluateByAttemptEnumAll)

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

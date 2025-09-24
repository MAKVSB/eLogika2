package helpers

import (
	"fmt"
	"slices"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

func QuestionMeetsStepsRequirements(question models.Question, requiredSteps []models.Step, stepsMode enums.StepSelectionEnum) bool {
	switch stepsMode {

	case enums.StepSelectionNC: // Only questions that use all selected steps and no others
		if len(requiredSteps) != len(question.CourseLink.Steps) { // Precheck of the same count
			return false
		}

		for _, cs := range requiredSteps { // Verifies that they are indeed the ones selected
			if !slices.ContainsFunc(question.CourseLink.Steps, func(qs models.Step) bool {
				return qs.ID == cs.ID
			}) {
				return false
			}
		}
		return true
	case enums.StepSelectionND: // Pouze otázky, které používají alespoň jeden z vybraných kroků a žádný z nevybraných
		for _, cs := range question.CourseLink.Steps { // All question steps are within requested
			if !slices.ContainsFunc(requiredSteps, func(s models.Step) bool {
				return s.ID == cs.ID
			}) {
				return false
			}
		}

		for _, cs := range requiredSteps { // At least one of requested is present
			if slices.ContainsFunc(question.CourseLink.Steps, func(s models.Step) bool {
				return s.ID == cs.ID
			}) {
				return true
			}
		}
		return false
	case enums.StepSelectionSC: // Only questions that use all selected steps and possibly others
		for _, cs := range requiredSteps {
			if !slices.ContainsFunc(question.CourseLink.Steps, func(s models.Step) bool {
				return s.ID == cs.ID
			}) {
				return false
			}
		}
		return true
	case enums.StepSelectionSD: // Only questions that use at least one of the selected steps and possibly others
		for _, cs := range requiredSteps { // Ověřuje, že jsou to opravdu ty vybrané
			if slices.ContainsFunc(question.CourseLink.Steps, func(s models.Step) bool {
				return s.ID == cs.ID
			}) {
				return true
			}
		}
		return false
	default:
		// TODO some error reporting
		fmt.Println("unexpected enums.StepSelectionEnum: " + stepsMode)
		return false
	}
}

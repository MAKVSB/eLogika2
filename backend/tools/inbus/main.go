package main

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/utils"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	inbusClient := GetInbusClient()

	semester, err := inbusClient.GetSemesterFromDate("2025-09-09")
	if err != nil {
		panic(err)
	}
	utils.DebugPrintJSON(semester)

	subject, err := inbusClient.GetSubjectVersionFromcode("460-2051/03")
	if err != nil {
		panic(err)
	}
	utils.DebugPrintJSON(subject)

	concreteActivities, err := inbusClient.GetConcreteActivities((*subject)[0].SubjectVersionId, &semester.SemesterId)
	if err != nil {
		panic(err)
	}
	utils.DebugPrintJSON(concreteActivities)

	concreteActivityStudents, err := inbusClient.GetConcreteActivityStudents((*concreteActivities)[1].ConcreteActivityId)
	if err != nil {
		panic(err)
	}
	utils.DebugPrintJSON(concreteActivityStudents)

}

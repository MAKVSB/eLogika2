package main

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type InbusSemester struct {
	// AcademicYearId     uint   `json:"academicYearId"`
	// AcademicYearTitle  string `json:"academicYearTitle"`
	BeginDate string `json:"beginDate"`
	// EducationWeeks     uint   `json:"educationWeeks"`
	EndDate string `json:"endDate"`
	// ExaminationWeeks   uint   `json:"examinationWeeks"`
	SemesterId uint `json:"semesterId"`
	// SemesterTypeAbbrev string `json:"semesterTypeAbbrev"`
	// SemesterTypeId     uint   `json:"semesterTypeId"`
	// SemesterTypeTitle  string `json:"semesterTypeTitle"`
}

type SubjectVersion struct {
	// Credits   float64 `json:"credits"`
	// Elearning string  `json:"elearning"`
	Guarantee struct {
		// 	FullName string `json:"fullName"`
		// 	Login    string `json:"login"`
		PersonId uint `json:"personId"`
	} `json:"guarantee"`
	// OtherRequirements          string `json:"otherRequirements"`
	// Schema                     string `json:"schema"`
	SubjectVersionCompleteCode string `json:"subjectVersionCompleteCode"`
	SubjectVersionId           uint   `json:"subjectVersionId"`
	// InstallYear struct {
	// 	AcademicYearId uint   `json:"academicYearId"`
	// 	Title          string `json:"title"`
	// } `json:"installYear"`
	// Language struct {
	// 	Code       string `json:"code"`
	// 	LanguageId uint   `json:"languageId"`
	// 	Title      string `json:"title"`
	// } `json:"language"`
	Subject struct {
		Abbrev string `json:"abbrev"`
		Title  string `json:"title"`
		// AdvisedLiterature string `json:"advisedLiterature"`
		// Aims              string `json:"aims"`
		// Annotation        string `json:"annotation"`
		// Code              string `json:"code"`
		// Department        struct {
		// 	Abbrev    string `json:"abbrev"`
		// 	OrgUnitId uint   `json:"orgUnitId"`
		// 	Title     string `json:"title"`
		// } `json:"department"`
		// Literature string `json:"literature"`
		// 	SubjectId uint   `json:"subjectId"`
	}
}

type StudyRelation struct {
	// StudyRelationId      uint   `json:"studyRelationId"`
	PersonId uint   `json:"personId"`
	Login    string `json:"login"`
	// FullName             string `json:"fullName"`
	Email string `json:"email"`
	// StudyCode            string `json:"studyCode"`
	// BeginDate            string `json:"beginDate"`
	// EndDate              string `json:"endDate"`
	// FacultyId            uint   `json:"facultyId"`
	// FacultyAbbrev        string `json:"facultyAbbrev"`
	// FacultyTitle         string `json:"facultyTitle"`
	// StudyTypeId          uint   `json:"studyTypeId"`
	// StudyTypeCode        string `json:"studyTypeCode"`
	// StudyTypeTitle       string `json:"studyTypeTitle"`
	// StudyFormId          uint   `json:"studyFormId"`
	// StudyFormCode        string `json:"studyFormCode"`
	// StudyFormTitle       string `json:"studyFormTitle"`
	// StudyLanguageId      uint   `json:"studyLanguageId"`
	StudyLanguageCode string `json:"studyLanguageCode"`
	// StudyLanguageTitle   string `json:"studyLanguageTitle"`
	// ClassYearId          uint   `json:"classYearId"`
	// ClassYearYear        uint   `json:"classYearYear"`
	// TutorialCentreId     uint   `json:"tutorialCentreId"`
	// TutorialCentreAbbrev string `json:"tutorialCentreAbbrev"`
	// TutorialCentreTitle  string `json:"tutorialCentreTitle"`
	// StudyProgrammeId     uint   `json:"studyProgrammeId"`
	// StudyProgrammeCode   string `json:"studyProgrammeCode"`
	// StudyProgrammeAbbrev string `json:"studyProgrammeAbbrev"`
	// StudyProgrammeTitle  string `json:"studyProgrammeTitle"`
	FirstName    string `json:"firstName"`
	SecondName   string `json:"secondName"`
	DegreeBefore string `json:"degreeBefore"`
}

type ConcreteActivity struct {
	ConcreteActivityId uint `json:"concreteActivityId"`
	// Template                   string `json:"template"`
	// Order                      uint   `json:"order"`
	// SubjectVersionId           uint   `json:"subjectVersionId"`
	// SubjectVersionCompleteCode string `json:"subjectVersionCompleteCode"`
	// SubjectId                  uint   `json:"subjectId"`
	// SubjectAbbrev string `json:"subjectAbbrev"`
	// SubjectTitle  string `json:"subjectTitle"`
	// EducationTypeId            uint   `json:"educationTypeId"`
	EducationTypeAbbrev string `json:"educationTypeAbbrev"`
	EducationTypeTitle  string `json:"educationTypeTitle"`
	// SemesterId            uint   `json:"semesterId"`
	// SemesterTypeId        uint   `json:"semesterTypeId"`
	// SemesterTypeAbbrev    string `json:"semesterTypeAbbrev"`
	// SemesterTypeTitle     string `json:"semesterTypeTitle"`
	// AcademicYearId        uint   `json:"academicYearId"`
	// AcademicYearTitle     string `json:"academicYearTitle"`
	// TutorialCentreId      uint   `json:"tutorialCentreId"`
	// TutorialCentreAbbrev  string `json:"tutorialCentreAbbrev"`
	// TutorialCentreTitle   string `json:"tutorialCentreTitle"`
	// EducationWeekId       uint   `json:"educationWeekId"`
	EducationWeekTitle string `json:"educationWeekTitle"`
	// BeginScheduleWindowId uint   `json:"beginScheduleWindowId"`
	// ActivityDuration uint   `json:"activityDuration"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
	WeekDayId uint   `json:"weekDayId"`
	// WeekDayAbbrev    string `json:"weekDayAbbrev"`
	// WeekDayTitle     string `json:"weekDayTitle"`
	WekActivities []struct {
		WeekActivityId uint   `json:"weekActivityId"`
		WeekNumber     uint   `json:"weekNumber"`
		Date           string `json:"date"`
	} `json:"weekActivities"`
	RoomFullcodes string `json:"roomFullcodes"`
	// TeacherLogins     string `json:"teacherLogins"`
	// TeacherShortNames string `json:"teacherShortNames"`
	// TeacherFullNames  string `json:"teacherFullNames"`
	// StudyGroupCodes   string `json:"studyGroupCodes"`
	// RoomIds           []uint `json:"roomIds"`
	TeacherIds    []uint `json:"teacherIds"`
	StudyGroupIds []uint `json:"studyGroupIds"`
}

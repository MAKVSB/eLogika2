package helpers

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
	"github.com/skip2/go-qrcode"
)

const maxQuestionsPerSheet = 18
const pageSpacing = float64(30)
const headingHeight = float64(105)
const maxHeadingLineChars = 20

type SheetData struct {
	Format         enums.QuestionFormatEnum
	SheetOrder     uint
	MaxAnswerCount int
	Questions      []models.TestQuestion
}

type AnswerSheetPrinter struct {
	WorkDir    string
	OutputDir  string
	OutputName string
}

func (asp AnswerSheetPrinter) GenerateAnswerSheets(testData *models.Test, testInstance *models.TestInstance) string {
	// Divide questions into answer sheets
	sheets := SplitToSheets(testData)

	// Prepare pdf document
	pdf := fpdf.New("P", "pt", "A4", "")
	pdf.SetFontLocation(asp.WorkDir)
	pdf.AddUTF8Font("Arial", "", filepath.Join("assets", "fonts", "ARIAL.TTF"))
	pdf.SetFont("Arial", "", 11)
	pdf.SetLineWidth(2.5)

	// Genearte sheets
	for _, sheet := range sheets {
		pdf.AddPage()
		pw, ph, _ := pdf.PageSize(0)

		if sheet != nil {

		}

		//Draw heading
		{
			// Draw triangle (filled)
			points := []fpdf.PointType{
				{X: pageSpacing, Y: pageSpacing},      // top point
				{X: pageSpacing, Y: pageSpacing + 50}, // bottom-left (right angle here)
				{X: pageSpacing + 50, Y: pageSpacing}, // bottom-right
			}
			pdf.Polygon(points, "D")                                             // D=draw, F=fill
			pdf.Line(pageSpacing-1.25, pageSpacing, pageSpacing+25, pageSpacing) // vertical

			// Title strings
			pdf.SetXY(70, 50)

			headingLines := make([]string, 0)
			if len(testData.Course.Name) > maxHeadingLineChars-9 {
				headingLines = append(headingLines, "Předmět: "+testData.Course.Name[:maxHeadingLineChars-9-3]+"...")
			} else {
				headingLines = append(headingLines, "Předmět: "+testData.Course.Name)
			}

			headingLines = append(headingLines, "Datum: "+testData.Term.ActiveFrom.Format("02.01.2006"))
			headingLines = append(headingLines, "Čas: "+testData.Term.ActiveFrom.Format("15:04")+" - "+testData.Term.ActiveTo.Format("15:04"))
			// TODO Here maybe put allowed length of the test as the end time
			headingLines = append(headingLines, "Test: "+testData.CourseItem.Name)
			headingLines = append(headingLines, "Variant: "+testData.Group)

			pdf.MultiCell(220, 14, strings.Join(headingLines, "\n"), "", "L", false)

			// Draw test instance ID-QR
			DrawTestIdentifierQR(pdf, pw-pageSpacing, pageSpacing, 90, strconv.Itoa(int(testData.CourseID))+"-"+strconv.Itoa(int(testData.ID)))

			// Draw participant QR
			DrawParticipantQR(pdf, pw-pageSpacing-100, pageSpacing, 90, testInstance)
		}
		//big rect
		pdf.Rect(pageSpacing, pageSpacing+headingHeight, pw-pageSpacing*2, ph-pageSpacing*2-headingHeight, "")

		DrawAnswers(pdf, pageSpacing, pageSpacing+headingHeight, pw-pageSpacing*2, ph-headingHeight-pageSpacing*2, sheet)
		pdf.AddPage() // Make sure the other side of the page is empty for both-sided printing
	}

	outputPath := filepath.Join(asp.OutputDir, asp.OutputName)
	// Save the PDF
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		panic(err)
	}
	return outputPath
}

func SplitToSheets(testData *models.Test) []*SheetData {
	var sheets []*SheetData

	var lastTeacherSheet *SheetData
	lastTeacherSheetOrder := uint(0)
	var lastStudentSheet *SheetData
	lastStudentSheetOrder := uint(0)

	for _, q := range testData.Questions {
		switch q.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			if lastTeacherSheet == nil {
				newSheet := &SheetData{
					Format:         enums.QuestionFormatOpen,
					SheetOrder:     lastTeacherSheetOrder,
					Questions:      []models.TestQuestion{},
					MaxAnswerCount: 11,
				}
				lastTeacherSheetOrder++
				sheets = append(sheets, newSheet)
				lastTeacherSheet = newSheet
			}

			lastTeacherSheet.Questions = append(lastTeacherSheet.Questions, q)

			if len(lastTeacherSheet.Questions) >= maxQuestionsPerSheet {
				lastTeacherSheet = nil
			}
		case enums.QuestionFormatTest:
			if lastStudentSheet == nil {
				newSheet := &SheetData{
					Format:     enums.QuestionFormatTest,
					SheetOrder: lastStudentSheetOrder,
					Questions:  []models.TestQuestion{},
				}
				lastStudentSheetOrder++
				sheets = append(sheets, newSheet)
				lastStudentSheet = newSheet
			}

			lastStudentSheet.Questions = append(lastStudentSheet.Questions, q)
			if len(q.Answers) > lastStudentSheet.MaxAnswerCount {
				lastStudentSheet.MaxAnswerCount = len(q.Answers)
			}

			if len(lastStudentSheet.Questions) >= maxQuestionsPerSheet {
				lastStudentSheet = nil
			}
		default:
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", q.Question.QuestionFormat))
		}
	}

	return sheets
}

func DrawAnswers(pdf *fpdf.Fpdf, posX float64, posY float64, w float64, h float64, sheet *SheetData) {
	posX += 25

	squareSize := float64(16)
	offsetYAdd := float64(30)
	offsetYAddHeader := float64(20)

	offsetY := float64(20)

	for i, q := range sheet.Questions {
		if i == 0 {
			DrawHeader(pdf, posX, posY+offsetY, squareSize, sheet)
			offsetY += offsetYAddHeader
		}

		questionString := strconv.Itoa(int(q.Order))
		switch q.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			DrawAnswerRow(pdf, posX, posY+offsetY, squareSize, questionString, 11)
		case enums.QuestionFormatTest:
			DrawAnswerRow(pdf, posX, posY+offsetY, squareSize, questionString, len(q.Answers))
		default:
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", q.Question.QuestionFormat))
		}

		if i == maxQuestionsPerSheet-1 {
			offsetY += offsetYAddHeader + 5
			DrawHeader(pdf, posX, posY+offsetY, squareSize, sheet)
		} else if i == (maxQuestionsPerSheet/2)-1 {
			offsetY += offsetYAddHeader + 10
			DrawHeader(pdf, posX, posY+offsetY, squareSize, sheet)
			offsetY += offsetYAddHeader + 2
		} else {
			offsetY += offsetYAdd
		}
	}
}

func DrawHeader(pdf *fpdf.Fpdf, posX float64, posY float64, squareSize float64, sheet *SheetData) {
	switch sheet.Format {
	case enums.QuestionFormatOpen:
		DrawHeaderPercentage(pdf, posX, posY, squareSize)
	case enums.QuestionFormatTest:
		DrawHeaderTest(pdf, posX, posY, squareSize, sheet.MaxAnswerCount)
	default:
		panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", sheet.Format))
	}
}

func DrawHeaderPercentage(pdf *fpdf.Fpdf, posX float64, posY float64, squareSize float64) {
	localSquareSize := squareSize * 0.97
	for i := 0; i < 11; i++ {
		pdf.SetXY(posX+float64(i+1)*localSquareSize, posY)
		pdf.Cell(10, 10, strconv.Itoa(i*10))
	}

	repairAnswersOffset := float64(45)
	posX += 11*squareSize + repairAnswersOffset

	for i := range 11 {
		pdf.SetXY(posX+float64(i+1)*localSquareSize, posY)
		pdf.Cell(10, 10, strconv.Itoa(i*10))
	}

	repairColumnOffset := float64(60)
	posX += 11*squareSize + repairColumnOffset

	pdf.SetXY(posX+10, posY)
	pdf.Cell(10, 10, "op")
}

func DrawHeaderTest(pdf *fpdf.Fpdf, posX float64, posY float64, squareSize float64, answerCount int) {
	localSquareSize := squareSize * 1.01
	for i := 0; i < answerCount; i++ {
		pdf.SetXY(posX+float64(i+1)*localSquareSize, posY)
		pdf.Cell(10, 10, GetTestAnswerLabel(i))
	}

	repairAnswersOffset := float64(45)
	posX += 11*squareSize + repairAnswersOffset

	for i := 0; i < answerCount; i++ {
		pdf.SetXY(posX+float64(i+1)*localSquareSize, posY)
		pdf.Cell(10, 10, GetTestAnswerLabel(i))
	}

	repairColumnOffset := float64(60)
	posX += 11*squareSize + repairColumnOffset

	pdf.SetXY(posX+10, posY)
	pdf.Cell(10, 10, "op")
}

func DrawAnswerRow(pdf *fpdf.Fpdf, posX float64, posY float64, squareSize float64, questionHeader string, answerCount int) {
	pdf.SetXY(posX-8, posY+3)
	pdf.Cell(10, 10, questionHeader+".")

	for i := float64(1); i < float64(answerCount)+1; i++ {
		pdf.Rect(posX+i*squareSize, posY, squareSize, squareSize, "")
	}

	repairAnswersOffset := float64(45)
	posX += 11*squareSize + repairAnswersOffset

	pdf.SetXY(posX-8, posY+3)
	pdf.Cell(10, 10, questionHeader+".")

	for i := float64(1); i < float64(answerCount)+1; i++ {
		pdf.Rect(posX+i*squareSize, posY, squareSize, squareSize, "")
	}

	repairColumnOffset := float64(60)
	posX += 11*squareSize + repairColumnOffset

	pdf.Rect(posX, posY, 2*squareSize, squareSize, "")
}

func DrawTestIdentifierQR(pdf *fpdf.Fpdf, toprightX float64, toprightY float64, size float64, content string) {
	identifier := "testInstanceQR"

	qr, err := qrcode.Encode(content, qrcode.High, 512)
	if err != nil {
		log.Fatal(err)
	}

	imgOpts := fpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}
	pdf.RegisterImageOptionsReader(identifier, imgOpts, bytes.NewReader(qr))

	topleftX := toprightX - size

	pdf.ImageOptions(identifier, topleftX, toprightY, size, size, false, imgOpts, 0, "")

	// Draw rectangle around it
	pdf.Rect(topleftX, toprightY, size, size, "")
}

func DrawParticipantQR(pdf *fpdf.Fpdf, toprightX float64, toprightY float64, size float64, testInstance *models.TestInstance) {
	identifier := "participantQR"
	width := float64(180)

	topleftX := toprightX - width
	topleftXQR := toprightX - size

	if testInstance != nil {
		qr, err := qrcode.Encode(strconv.Itoa(int(testInstance.ID)), qrcode.High, 512)
		if err != nil {
			log.Fatal(err)
		}

		imgOpts := fpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}
		pdf.RegisterImageOptionsReader(identifier, imgOpts, bytes.NewReader(qr))
		pdf.ImageOptions(identifier, topleftXQR, toprightY, size, size, false, imgOpts, 0, "")

		participantLines := []string{}
		participantLines = append(participantLines, testInstance.Participant.DegreeBefore+" "+testInstance.Participant.FirstName)
		participantLines = append(participantLines, testInstance.Participant.FamilyName+" "+testInstance.Participant.DegreeAfter)
		participantLines = append(participantLines, testInstance.Participant.Username)

		pdf.SetXY(topleftX+10, toprightY+10)
		pdf.MultiCell(220, 14, strings.Join(participantLines, "\n"), "", "L", false)
	}

	// Draw rectangle around it
	pdf.Rect(topleftX, toprightY, width, size, "")
}

func GetTestAnswerLabel(n int) string {
	label := ""
	for n >= 0 {
		label = string(rune('a'+(n%26))) + label
		n = n/26 - 1
	}
	return label
}

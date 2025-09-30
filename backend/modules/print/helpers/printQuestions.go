package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"elogika.vsb.cz/backend/models"
)

type QuestionPrinter struct {
	WorkDir   string
	AssetDir  string
	Outputdir string
}

func (qp QuestionPrinter) PrintQuestions(questions []*models.Question) string {
	latexCode := `
	\documentclass{article}
	\usepackage[a4paper, margin=15mm]{geometry}
	\usepackage[utf8]{inputenc}
	\usepackage{amsmath, amssymb}
	\usepackage{ulem}
	\usepackage{graphicx}
	\usepackage{tabularx}
	\usepackage{enumitem}
	\usepackage{hyperref}
	\usepackage{multirow}
	\usepackage{color,soul}
	\newcolumntype{P}[1]{>{\vspace{6pt}\arraybackslash}m{#1}<{\vspace{6pt}}}
	
	% Define counters
	\newcounter{question}

	% Format for block (like a chapter)
	\newcommand{\block}[1]{%
		\bigskip\par\noindent{\LARGE\textbf{#1}}\medskip
	}

	% Format for question (question, continuous numbering)
	\newcommand{\question}{%
		\refstepcounter{question}%
		\bigskip\par\noindent{\Large\textbf{\thequestion. otázka}}\par\medskip
	}

	\begin{document}
	`

	var lastChapterID uint
	var lastCategoryID uint

	for _, question := range questions {
		var qContent map[string]interface{}
		if err := json.Unmarshal(question.Content, &qContent); err != nil {
			panic(err)
		}

		latexCode += `\begin{samepage}`

		if question.CourseLink.ChapterID != lastChapterID {
			latexCode += `\block{Kapitola: ` + question.CourseLink.Chapter.Name + " " + strconv.Itoa(int(question.CourseLink.ChapterID)) + `}`
			lastChapterID = question.CourseLink.ChapterID
		}

		if question.CourseLink.CategoryID != nil {
			if *question.CourseLink.CategoryID != lastCategoryID {
				if question.CourseLink.CategoryID != nil && question.CourseLink.Category != nil {
					latexCode += `\block{Kategorie: ` + question.CourseLink.Category.Name + " " + strconv.Itoa(int(*question.CourseLink.CategoryID)) + `}`
				}
				lastCategoryID = *question.CourseLink.CategoryID
			}
		} else {
			if lastCategoryID != 0 {
				latexCode += `\block{Kategorie: Bez kategorie}`
			}
			lastCategoryID = 0
		}

		latexCode += `\question `
		latexCode += qp.ConvertNodeToLaTeX(qContent)

		if len(question.Answers) == 0 {

			latexCode += `\textcolor{red}{No answers found}`
			latexCode += `\vspace{5cm plus 10cm}`
		} else {
			latexCode += `\begin{enumerate}[label=\alph*)]`
			for _, answer := range question.Answers {
				var aContent map[string]interface{}
				if err := json.Unmarshal(answer.Answer.Content, &aContent); err != nil {
					panic(err)
				}
				if answer.Answer.Correct {
					latexCode += `\item \textcolor{green}{\rule{1em}{1em}} \hspace{0.2em}`
				} else {
					latexCode += `\item \textcolor{red}{\rule{1em}{1em}} \hspace{0.2em}`
				}
				latexCode += qp.ConvertNodeToLaTeX(aContent)
			}
			latexCode += `\end{enumerate}`
		}

		latexCode += `\end{samepage}`
	}

	latexCode += `\end{document}`

	// Write LaTeX code to file
	texFile := qp.Outputdir + "/output.tex"
	var err error
	err = os.WriteFile(texFile, []byte(latexCode), 0644)
	if err != nil {
		panic(err)
	}

	// Run pdflatex
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory="+qp.Outputdir+"/", texFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	return filepath.Join(qp.Outputdir, "output.pdf")
}

func (qp QuestionPrinter) ConvertNodeToLaTeX(node map[string]interface{}) string {
	typ := node["type"].(string)
	content, _ := node["content"].([]interface{})

	switch typ {
	case "heading":
		level := int(node["attrs"].(map[string]interface{})["level"].(float64))
		prefix := "\\LARGE"
		if level == 2 {
			prefix = "\\Large"
		}
		if level == 3 {
			prefix = "\\large"
		}
		inner := qp.concatContent(content)
		return fmt.Sprintf("%s{%s}\n\n", prefix, inner)

	case "paragraph":
		return qp.concatContent(content) + "\n\n"

	case "text":
		text := node["text"].(string)
		marks, _ := node["marks"].([]interface{})
		for _, m := range marks {
			mark := m.(map[string]interface{})["type"].(string)
			switch mark {
			case "bold":
				text = "\\textbf{" + text + "}"
			case "italic":
				text = "\\textit{" + text + "}"
			case "underline":
				text = "\\underline{" + text + "}"
			case "strike":
				text = "\\sout{" + text + "}"
			case "highlight":
				text = "\\hl{\\mbox{" + text + "}}"
			case "link":
				attrs := m.(map[string]interface{})["attrs"]
				href := attrs.(map[string]interface{})["href"].(string)
				text = "\\href{" + href + "}{\\mbox{" + text + "}}"
			case "subscript":
				text = "\\textsubscript{" + text + "}"
			case "superscript":
				text = "\\textsuperscript{" + text + "}"
			default:
				panic(mark)
			}
		}
		return text

	case "inlineMath":
		latex := node["attrs"].(map[string]interface{})["latex"].(string)
		return "$" + latex + "$"

	case "blockMath":
		latex := node["attrs"].(map[string]interface{})["latex"].(string)
		return `\begin{equation*}
		` + latex + `
		\end{equation*}` + "\n\n\n"

	case "bulletList":
		items := ""
		for _, c := range content {
			items += "\\item " + qp.ConvertNodeToLaTeX(c.(map[string]interface{})) + ""
		}
		return "\\begin{itemize}[nosep,align=left,leftmargin=*]\n" + items + "\\end{itemize}\n"

	case "orderedList":
		items := ""
		for _, c := range content {
			items += "\\item " + qp.ConvertNodeToLaTeX(c.(map[string]interface{})) + ""
		}
		return "\\begin{enumerate}[nosep,align=left,leftmargin=*]\n" + items + "\\end{enumerate}\n"

	case "custom-image":
		src := node["attrs"].(map[string]interface{})["src"].(string)
		width := node["attrs"].(map[string]interface{})["width"].(float64)
		mode := node["attrs"].(map[string]interface{})["mode"].(string)
		if mode == "storage" {
			linkFile(qp.WorkDir+"/uploads/"+src, qp.AssetDir+"/"+src)
			src = strings.ReplaceAll(qp.AssetDir+"/"+src, "\\", "/")
		} else {
			src2, err := downloadImage(src, qp.AssetDir)
			if err != nil {
				panic(err)
				// TODO throw correct error
			} else {
				src = src2
			}
		}

		width = min(width/2, 500)

		return fmt.Sprintf("\\vspace{10pt}  \\includegraphics[width="+strconv.Itoa(int(width))+"pt, keepaspectratio]{%s}\n\n  \\vspace{10pt}", src)

	case "table":
		type TableCell struct {
			Rowspan int
			Colspan int
			Content []interface{} // could be paragraph, list, image, or raw Tiptap JSON
		}

		tableCols := getTableColsCount(node)
		tableRows := len(node["content"].([]interface{}))

		reservationMatrix := make([][]bool, tableRows)
		for i := range reservationMatrix {
			reservationMatrix[i] = make([]bool, tableCols)
		}

		clineMatrix := make([][]bool, tableRows)
		for i := range reservationMatrix {
			clineMatrix[i] = make([]bool, tableCols)
		}

		widthMatrix := make([]int, tableCols)

		matrix := make([][]*TableCell, tableRows)
		for i := range matrix {
			matrix[i] = make([]*TableCell, tableCols)
		}

		for row_index, row_d := range node["content"].([]interface{}) {
			row_data := row_d.(map[string]interface{})

			correct_cell_index := 0

			for _, row_cell_d := range row_data["content"].([]interface{}) {
				row_cell_data := row_cell_d.(map[string]interface{})
				colspan := row_cell_data["attrs"].(map[string]interface{})["colspan"].(float64)
				rowspan := row_cell_data["attrs"].(map[string]interface{})["rowspan"].(float64)

				tableCell := TableCell{
					Rowspan: int(rowspan),
					Colspan: int(colspan),
					Content: row_cell_data["content"].([]interface{}),
				}

				for reservationMatrix[row_index][correct_cell_index] {
					correct_cell_index += 1
				}

				matrix[row_index][correct_cell_index] = &tableCell

				for i := 0; i < int(rowspan); i++ {
					for j := 0; j < int(colspan); j++ {
						reservationMatrix[row_index+i][correct_cell_index+j] = true
						if i == int(rowspan-1) {
							clineMatrix[row_index+i][correct_cell_index+j] = true
						}
					}
					matrix[row_index+i][correct_cell_index] = &TableCell{
						Rowspan: 1,
						Colspan: int(colspan),
					}
				}

				if tableCell.Colspan == 1 {
					colWidth, ok := row_cell_data["attrs"].(map[string]interface{})["colwidth"].([]interface{})
					if ok {
						widthMatrix[correct_cell_index] = int(colWidth[0].(float64))
					}
				}

				if tableCell.Rowspan == 1 {
					clineMatrix[row_index][correct_cell_index] = true
				}
				matrix[row_index][correct_cell_index] = &tableCell

				correct_cell_index += int(colspan)
			}
		}

		totalWidth := 0
		totalWidthIsStatic := true

		for _, w := range widthMatrix {
			if w == 0 {
				totalWidthIsStatic = false
			} else {
				totalWidth += (w / 2) + 13
			}
		}

		ratio := 500 / float64(totalWidth)
		if ratio < 1 {
			totalWidth = int(float64(totalWidth) * ratio)
		}

		colDef := "|"

		for _, w := range widthMatrix {
			if w == 0 {
				colDef += "X|"
			} else {
				colDef += `P{` + strconv.Itoa(int((float64(w)/2)*ratio)) + `pt}|`
			}
		}

		text_rows := make([]string, 0)

		for c_i, matrix_row := range matrix {
			row_cells_text := make([]string, 0)
			row_cells_clines := ""
			for c_j, matrix_cell := range matrix_row {
				if matrix_cell != nil {
					cell_text := "" + qp.concatContent(matrix_cell.Content) + ""
					if matrix_cell.Rowspan != 1 {
						cell_text = `\multirow{` + strconv.Itoa(int(matrix_cell.Rowspan)) + `}{*}{\parbox{100cm}{` + cell_text + `}}`
					}
					if matrix_cell.Colspan != 1 {
						cell_text = `\multicolumn{` + strconv.Itoa(int(matrix_cell.Colspan)) + `}{X|}{` + cell_text + `}`
					}
					row_cells_text = append(row_cells_text, cell_text)
				}

				if clineMatrix[c_i][c_j] {
					row_cells_clines += " \\cline{" + strconv.Itoa(c_j+1) + "-" + strconv.Itoa(c_j+1) + "} "
				}
			}

			text_rows = append(text_rows, strings.Join(row_cells_text, " & ")+"\\\\"+row_cells_clines)

		}

		tableWidth := `\linewidth`
		if totalWidthIsStatic {
			tableWidth = strconv.Itoa(totalWidth) + "pt"
		}

		return `
				\begin{tabularx}{` + tableWidth + `}{` + colDef + `}
					\cline{1-` + strconv.Itoa(tableCols) + `}
					` + strings.Join(text_rows, "") + `
				\end{tabularx}
		` + "\n\n"
	case "listItem":
		return qp.concatContent(content)
	case "tableHeader":
		return qp.concatContent(content)
	case "doc":
		return qp.concatContent(content)
	default:
		return qp.concatContent(content)
	}
}

func (qp QuestionPrinter) concatContent(content []interface{}) string {
	result := ""
	for _, c := range content {
		result += qp.ConvertNodeToLaTeX(c.(map[string]interface{}))
	}
	return result
}

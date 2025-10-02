package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
	"github.com/google/uuid"
)

type TestPrinter struct {
	WorkDir   string
	AssetDir  string
	Outputdir string
}

func (tp TestPrinter) GenerateTestContent(testData *models.Test) string {
	latexCode := `
	\documentclass{article}
	\usepackage[a4paper, margin=15mm]{geometry}
	\usepackage[utf8]{inputenc}
	\usepackage{amsmath, amssymb}
	\usepackage{ulem}
	\usepackage{graphicx}
	\usepackage{ifoddpage}
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

	lastBlockID := uint(0)

	for _, question := range testData.Questions {
		var qContent map[string]interface{}
		if err := json.Unmarshal(question.Question.Content, &qContent); err != nil {
			panic(err)
		}

		latexCode += `\begin{samepage}`

		if question.BlockID != lastBlockID {
			lastBlockID = question.BlockID
			for _, block := range testData.Blocks {
				if block.ID == question.BlockID {
					if block.ShowName {
						latexCode += `\block{` + block.Title + `}`
					}
					break
				}
			}
		}

		latexCode += `\question{} `
		latexCode += tp.ConvertNodeToLaTeX(qContent)

		switch question.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			latexCode += `\vspace{5cm plus 10cm}`
		case enums.QuestionFormatTest:
			latexCode += `\begin{enumerate}[label=\alph*)]`

			for _, answer := range question.Answers {
				var aContent map[string]interface{}
				if err := json.Unmarshal(answer.Answer.Content, &aContent); err != nil {
					panic(err)
				}
				latexCode += `\item ` + tp.ConvertNodeToLaTeX(aContent)
			}
			latexCode += `\end{enumerate}`
		default:
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", question.Question.QuestionFormat))
		}

		latexCode += `\end{samepage}`
	}

	latexCode += `
	\clearpage
	\checkoddpage
	\ifoddpage
	% do nothing, page count is odd
	\else
	% insert a blank page
	\thispagestyle{empty}
	\null
	\newpage
	\fi
	`

	latexCode += `\end{document}`

	// Write LaTeX code to file

	uniqueName := "output" + uuid.NewString()

	texFile := tp.Outputdir + "/" + uniqueName + ".tex"
	var err error
	err = os.WriteFile(texFile, []byte(latexCode), 0644)
	if err != nil {
		panic(err)
	}

	// Run pdflatex
	cmd := exec.Command("pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory="+tp.Outputdir+"/", texFile)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.Stdin = nil
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	return filepath.Join(tp.Outputdir, uniqueName+".pdf")
}

func (tp TestPrinter) ConvertNodeToLaTeX(node map[string]interface{}) string {
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
		inner := tp.concatContent(content)
		return fmt.Sprintf("%s{%s}\n\n", prefix, inner)

	case "paragraph":
		return tp.concatContent(content) + "\n\n"

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
			items += "\\item " + tp.ConvertNodeToLaTeX(c.(map[string]interface{})) + ""
		}
		return "\\begin{itemize}[nosep,align=left,leftmargin=*]\n" + items + "\\end{itemize}\n"

	case "orderedList":
		items := ""
		for _, c := range content {
			items += "\\item " + tp.ConvertNodeToLaTeX(c.(map[string]interface{})) + ""
		}
		return "\\begin{enumerate}[nosep,align=left,leftmargin=*]\n" + items + "\\end{enumerate}\n"

	case "custom-image":
		src := node["attrs"].(map[string]interface{})["src"].(string)
		width := node["attrs"].(map[string]interface{})["width"].(float64)
		mode := node["attrs"].(map[string]interface{})["mode"].(string)
		if mode == "storage" {
			linkFile(tp.WorkDir+"/uploads/"+src, tp.AssetDir+"/"+src)
			src = strings.ReplaceAll(tp.AssetDir+"/"+src, "\\", "/")
		} else {
			src2, err := downloadImage(src, tp.AssetDir)
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
					cell_text := "" + tp.concatContent(matrix_cell.Content) + ""
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
		return tp.concatContent(content)
	case "tableHeader":
		return tp.concatContent(content)
	case "doc":
		return tp.concatContent(content)
	default:
		return tp.concatContent(content)
	}
}

func (tp TestPrinter) concatContent(content []interface{}) string {
	result := ""
	for _, c := range content {
		result += tp.ConvertNodeToLaTeX(c.(map[string]interface{}))
	}
	return result
}

func downloadImage(url, folder string) (string, error) {
	// Compute MD5 hash of the URL to use as filename
	hash := md5.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + filepath.Ext(url) // keep extension if any
	filepath := filepath.Join(folder, filename)

	// Check if file already exists
	if _, err := os.Stat(filepath); err == nil {
		// File already downloaded
		return filepath, nil
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download: %s", resp.Status)
	}

	// Create local file
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write data to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func linkFile(source string, destination string) {
	// Check if file already exists
	if _, err := os.Stat(destination); err == nil {
		return
	}

	err := os.Link(source, destination)
	if err != nil {
		panic(err)
	}
	return
}

func getTableColsCount(node map[string]interface{}) int {
	cols := float64(0)

	firstRowData := node["content"].([]interface{})[0].(map[string]interface{})
	for _, firstRowCell := range firstRowData["content"].([]interface{}) {
		attrs := firstRowCell.(map[string]interface{})["attrs"]
		cols += attrs.(map[string]interface{})["colspan"].(float64)
	}

	return int(cols)
}

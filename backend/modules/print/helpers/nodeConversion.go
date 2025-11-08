package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/utils"
)

type NodeConvertor struct {
	WorkDir  string
	AssetDir string
}

func (qp NodeConvertor) concatContent(content []*models.TipTapContent) (string, error) {
	result := ""
	for _, c := range content {
		if res, err := qp.ConvertNodeToLaTeX(c); err != nil {
			return "", err
		} else {
			result += res
		}
	}
	return result, nil
}

func (tp NodeConvertor) ConvertNodeToLaTeX(node *models.TipTapContent) (string, error) {
	switch node.Type {
	case "heading":
		level := int(node.Attrs["level"].(float64))
		prefix := "\\LARGE"
		if level == 2 {
			prefix = "\\Large"
		}
		if level == 3 {
			prefix = "\\large"
		}
		inner, err := tp.concatContent(node.Content)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s{%s}\n\n", prefix, inner), nil

	case "paragraph":
		if res, err := tp.concatContent(node.Content); err != nil {
			return "", err
		} else {
			return res + "\n\n", nil
		}

	case "text":
		text := node.Text
		marks := node.Marks
		for _, m := range marks {
			mark := m.Type
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
				attrs := m.Attrs
				href := attrs["href"].(string)
				text = "\\href{" + href + "}{\\mbox{" + text + "}}"
			case "subscript":
				text = "\\textsubscript{" + text + "}"
			case "superscript":
				text = "\\textsuperscript{" + text + "}"
			}
		}
		return text, nil

	case "inlineMath":
		latex := node.Attrs["latex"].(string)
		return "$" + latex + "$", nil

	case "blockMath":
		latex := node.Attrs["latex"].(string)
		return `\begin{equation*}` + latex + `\end{equation*}\n\n\n`, nil

	case "bulletList":
		items := ""
		for _, c := range node.Content {
			if res, err := tp.ConvertNodeToLaTeX(c); err != nil {
				return "", err
			} else {
				items += "\\item " + res
			}
		}
		return "\\begin{itemize}[nosep,align=left,leftmargin=*]\n" + items + "\\end{itemize}\n", nil

	case "orderedList":
		items := ""
		for _, c := range node.Content {
			if res, err := tp.ConvertNodeToLaTeX(c); err != nil {
				return "", err
			} else {
				items += "\\item " + res
			}
		}
		return "\\begin{enumerate}[nosep,align=left,leftmargin=*]\n" + items + "\\end{enumerate}\n", nil

	case "custom-image":
		src := node.Attrs["src"].(string)
		width := node.Attrs["width"].(float64)
		mode := node.Attrs["mode"].(string)
		if mode == "storage" {
			err := utils.LinkFile(tp.WorkDir+"/uploads/"+src, tp.AssetDir+"/"+src)
			if err != nil {
				return "", err
			}
			src = strings.ReplaceAll(tp.AssetDir+"/"+src, "\\", "/")
		} else {
			src2, err := utils.DownloadFile(src, tp.AssetDir)
			if err != nil {
				return "", err
			} else {
				src = src2
			}
		}

		width = min(width/2, 500)

		return fmt.Sprintf("\\vspace{10pt}  \\includegraphics[width=%dpt, keepaspectratio]{%s}\n\n  \\vspace{10pt}", width, src), nil

	case "table":
		type TableCell struct {
			Rowspan int
			Colspan int
			Content []*models.TipTapContent
		}

		tableCols := getTableColsCount(node)
		tableRows := len(node.Content)

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

		for row_index, row_data := range node.Content {
			correct_cell_index := 0

			for _, row_cell_data := range row_data.Content {
				colspan := row_cell_data.Attrs["colspan"].(float64)
				rowspan := row_cell_data.Attrs["rowspan"].(float64)

				tableCell := TableCell{
					Rowspan: int(rowspan),
					Colspan: int(colspan),
					Content: row_cell_data.Content,
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
					colWidth, ok := row_cell_data.Attrs["colwidth"].([]interface{})
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
					res, err := tp.concatContent(matrix_cell.Content)
					if err != nil {
						return "", err
					}

					cell_text := res
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
		\n\n`, nil
	case "listItem":
		return tp.concatContent(node.Content)
	case "tableHeader":
		return tp.concatContent(node.Content)
	case "doc":
		return tp.concatContent(node.Content)
	default:
		return tp.concatContent(node.Content)
	}
}

func getTableColsCount(node *models.TipTapContent) int {
	cols := float64(0)

	firstRowData := node.Content[0]
	for _, firstRowCell := range firstRowData.Content {
		cols += firstRowCell.Attrs["colspan"].(float64)
	}

	return int(cols)
}

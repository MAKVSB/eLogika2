package helpers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func GenerateTestContent(
	testData *models.Test,
	workDir string,
	assetDir string,
	outputDir string,
) (string, int, error) {
	nc := NodeConvertor{
		WorkDir:  workDir,
		AssetDir: assetDir,
	}

	latexCode := `
	\documentclass{article}
	\usepackage[a4paper, margin=15mm]{geometry}
	\usepackage[utf8]{inputenc}
	\usepackage{amsmath, amssymb}
	\usepackage{ulem}
	\usepackage{graphicx}
	\usepackage{tabularx}
	\usepackage{enumitem}
	\usepackage{booktabs}
	\usepackage{hyperref}
	\usepackage{pifont}
	\usepackage{tikz}
	\usetikzlibrary{arrows,positioning,chains,fit,matrix,shapes,calc,external,petri, arrows.meta}
	\usepackage{multirow}
	\usepackage{tcolorbox}
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
		if res, err := nc.ConvertNodeToLaTeX(question.Question.Content); err != nil {
			return "", 0, err
		} else {
			latexCode += res
		}

		switch question.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			latexCode += `\vspace{5cm plus 10cm}`
		case enums.QuestionFormatTest:
			latexCode += `\begin{enumerate}[label=\alph*)]`

			for _, answer := range question.Answers {
				if res, err := nc.ConvertNodeToLaTeX(answer.Answer.Content); err != nil {
					return "", 0, err
				} else {
					latexCode += `\item ` + res
				}
			}
			latexCode += `\end{enumerate}`
		default:
			return "", 0, fmt.Errorf("unexpected enums.QuestionFormatEnum: %#v", question.Question.QuestionFormat)
		}

		latexCode += `\end{samepage}`
	}

	latexCode += `
	\clearpage
	`

	latexCode += `\end{document}`

	// Write LaTeX code to file

	uniqueName := "output" + uuid.NewString()

	texFile := outputDir + "/" + uniqueName + ".tex"
	var err error
	err = os.WriteFile(texFile, []byte(latexCode), 0644)
	if err != nil {
		return "", 0, err
	}

	// Run pdflatex
	cmd := exec.Command("pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory="+outputDir+"/", texFile)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.Stdin = nil
	err = cmd.Run()
	if err != nil {
		return "", 0, err
	}

	returnPath := filepath.Join(outputDir, uniqueName+".pdf")

	ctx, err := api.ReadContextFile(returnPath)
	if err != nil {
		return "", 0, errors.New("cannot read number of pages")
	}

	return returnPath, ctx.PageCount, nil
}

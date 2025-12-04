package helpers

import (
	"os"
	"os/exec"
	"path/filepath"

	"elogika.vsb.cz/backend/models"
)

func classificationString(link *models.CourseQuestion) string {
	if link.Chapter != nil {
		if link.Category != nil {
			return link.Chapter.Name + ", " + link.Category.Name
		} else {
			return link.Chapter.Name
		}
	} else {
		return ""
	}
}

func PrintQuestions(
	questions []*models.Question,
	workDir string,
	assetDir string,
	outputDir string,
) (string, error) {
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
	\newcommand{\question}[2]{%
		\refstepcounter{question}%
		\bigskip\par\noindent{\Large\textbf{\thequestion. Otázka} (#2) #1}\par\medskip
	}

	\begin{document}
	`

	for _, question := range questions {
		latexCode += `\begin{tcolorbox}[colback=white, colframe=white, boxrule=0pt, left=0pt, right=0pt, top=0pt, bottom=0pt]`

		latexCode += `\question{` + question.Title + `}` + `{` + classificationString(question.CourseLink) + `}` + "\n"
		if res, err := nc.ConvertNodeToLaTeX(question.Content); err != nil {
			return "", err
		} else {
			latexCode += res
		}

		if len(question.Answers) == 0 {
			latexCode += `\textcolor{red}{No answers found}`
			latexCode += `\vspace{5cm plus 10cm}`
		} else {
			latexCode += `\begin{enumerate}`
			for _, answer := range question.Answers {
				if answer.Answer.Correct {
					latexCode += `\item \ding{51} \hspace{0.2em}`
				} else {
					latexCode += `\item \ding{55} \hspace{0.2em}`
				}
				if res, err := nc.ConvertNodeToLaTeX(answer.Answer.Content); err != nil {
					return "", err
				} else {
					latexCode += res
				}
			}
			latexCode += `\end{enumerate}`
		}

		latexCode += `\end{tcolorbox}`
	}

	latexCode += `\end{document}`

	// Write LaTeX code to file
	texFile := outputDir + "/output.tex"
	var err error
	err = os.WriteFile(texFile, []byte(latexCode), 0644)
	if err != nil {
		return "", err
	}

	// Run pdflatex
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory="+outputDir+"/", texFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return filepath.Join(outputDir, "output.pdf"), nil
}

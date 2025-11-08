package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"elogika.vsb.cz/backend/models"
)

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
	\usepackage{hyperref}
	\usepackage{tikz}
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
		if res, err := nc.ConvertNodeToLaTeX(question.Content); err != nil {
			return "", err
		} else {
			latexCode += res
		}

		if len(question.Answers) == 0 {

			latexCode += `\textcolor{red}{No answers found}`
			latexCode += `\vspace{5cm plus 10cm}`
		} else {
			latexCode += `\begin{enumerate}[label=\alph*)]`
			for _, answer := range question.Answers {
				if answer.Answer.Correct {
					latexCode += `\item \textcolor{green}{\rule{1em}{1em}} \hspace{0.2em}`
				} else {
					latexCode += `\item \textcolor{red}{\rule{1em}{1em}} \hspace{0.2em}`
				}
				if res, err := nc.ConvertNodeToLaTeX(answer.Answer.Content); err != nil {
					return "", err
				} else {
					latexCode += res
				}
			}
			latexCode += `\end{enumerate}`
		}

		latexCode += `\end{samepage}`
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

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"elogika.vsb.cz/backend/initializers"
	"gorm.io/gorm"
)

func ConvertCodeToImage(dbRef *gorm.DB, code string) (string, int64, error) {
	const latexTemplate = `
	\documentclass{standalone}
	\usepackage[utf8]{inputenc}
	\usepackage{amsmath, amssymb}
	\usepackage{ulem}
	\usepackage{graphicx}
	\usepackage{tabularx}
	\usepackage{enumitem}
	\usepackage{hyperref}
	\usepackage{tikz}
	\usetikzlibrary{arrows,positioning,chains,fit,matrix,shapes,calc,external,petri, arrows.meta}
	\usepackage{multirow}
	\usepackage{color,soul}
	\newcolumntype{P}[1]{>{\vspace{6pt}\arraybackslash}m{#1}<{\vspace{6pt}}}

	\begin{document}
	%s
	\end{document}
	`

	latexCode := fmt.Sprintf(latexTemplate, code)

	// Create temp folder
	workDir, err := os.Getwd()
	if err != nil {
		return "", 0, err
	}

	tmpDir, err := CreateTmpFolder(workDir)
	if err != nil {
		return "", 0, err
	}

	// Write LaTeX code to file
	texFile := tmpDir + "/output.tex"
	err = os.WriteFile(texFile, []byte(latexCode), 0644)
	if err != nil {
		return "", 0, err
	}

	// Run pdflatex
	cmd := exec.Command("latex", "-interaction=batchmode", "-output-directory="+tmpDir+"/", texFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", 0, err
	}

	// Convert to SVG
	cmd = exec.Command("dvisvgm", "--no-fonts", tmpDir+"/output.dvi", "-o", tmpDir+"/output.svg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", 0, err
	}

	//Rename and copy to uploads folder
	newFileName, err := GenerateFileName(dbRef, ".svg")
	if err != nil {
		return "", 0, err
	}

	err = CopyFile(tmpDir+"/output.svg", filepath.Join(initializers.GlobalAppConfig.UPLOADS_DESTINATION, newFileName))
	if err != nil {
		return "", 0, err
	}

	fileInfo, err := os.Stat(filepath.Join(initializers.GlobalAppConfig.UPLOADS_DESTINATION, newFileName))
	if err != nil {
		return "", 0, err
	}

	return newFileName, fileInfo.Size(), nil
}

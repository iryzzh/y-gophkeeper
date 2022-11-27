package tui

import (
	"github.com/AlecAivazis/survey/v2"
)

type ViewStruct struct {
	Title string
	ID    int
}

func View(vs []ViewStruct) (int, error) {
	titles := make([]string, len(vs))
	for i, m := range vs {
		titles[i] = m.Title
	}
	var qs = &survey.Select{
		Message: "Entry:",
		Options: titles,
	}

	answerIndex := 0
	err := survey.AskOne(qs, &answerIndex)
	if err != nil {
		return 0, err
	}

	return vs[answerIndex].ID, nil
}

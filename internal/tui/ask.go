package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iryzzh/gophkeeper/internal/models"
)

func AskString(message string, suggest func(string) []string) (string, error) {
	prompt := &survey.Input{
		Message: message,
		Suggest: suggest,
	}

	return askOne(prompt)
}

func AskSelect(message string, options []string, def string) (string, error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Default: def,
	}

	return askOne(prompt)
}

func askOne(prompt survey.Prompt) (string, error) {
	output := ""
	err := survey.AskOne(prompt, &output, survey.WithValidator(survey.Required))

	return output, err
}

func AskFile(message string, validate ...bool) (string, error) {
	f := func(toComplete string) []string {
		files, _ := filepath.Glob(toComplete + "*")
		return files
	}

	q := []*survey.Question{
		{
			Name: "file",
			Prompt: &survey.Input{
				Message: message,
				Suggest: f,
				Help:    "Any file",
			},
			Validate: func(ans interface{}) error {
				if err := survey.Required(ans); err != nil {
					return err
				}
				if len(validate) > 0 {
					if !validate[0] {
						return nil
					}
				}
				f, _ := ans.(string)
				_, err := os.ReadFile(f)

				return err
			},
			Transform: func(ans interface{}) (newAns interface{}) {
				abs, err := filepath.Abs(ans.(string))
				if err == nil {
					return abs
				}

				return ans
			},
		},
	}

	file := ""
	err := survey.Ask(q, &file)

	return file, err
}

func AskPassword() (string, error) {
	password := ""
	prompt := &survey.Password{
		Message: "password:",
	}
	err := survey.AskOne(prompt, &password)

	return password, err
}

func AskCard() (models.Card, error) {
	qs := []*survey.Question{
		{
			Name:      "type",
			Prompt:    &survey.Input{Message: "Card type:"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:   "number",
			Prompt: &survey.Input{Message: "Number:"},
			Validate: func(ans interface{}) error {
				v, _ := ans.(string)
				if len(v) < 15 || len(v) > 19 {
					return fmt.Errorf("invalid card number length")
				}
				return isValidLuhn(ans)
			},
		},
		{
			Name:     "expiryMonth",
			Prompt:   &survey.Input{Message: "Expiry Month:"},
			Validate: isValidMonthYear,
		},
		{
			Name:     "expiryYear",
			Prompt:   &survey.Input{Message: "Expiry Year:"},
			Validate: isValidMonthYear,
		},
		{
			Name:     "ccv",
			Prompt:   &survey.Input{Message: "CCV:"},
			Validate: isValidCCV,
		},
	}

	var card models.Card
	if err := survey.Ask(qs, &card); err != nil {
		return models.Card{}, err
	}

	return card, nil
}

type EntryAnswer struct {
	Name      string `survey:"name"`
	Value     string `survey:"value"`
	ValueType string `survey:"type"`
}

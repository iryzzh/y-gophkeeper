//nolint:gomnd
package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iryzzh/y-gophkeeper/internal/models"
)

func AskInit(m *models.Init) error {
	var qs = []*survey.Question{
		{
			Name:     "remote",
			Prompt:   &survey.Input{Message: "Remote:"},
			Validate: survey.Required,
		},
		{
			Name:     "user",
			Prompt:   &survey.Input{Message: "User:"},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "password:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Remote   string `survey:"remote"`
		User     string `survey:"user"`
		Password string `survey:"password"`
	}{}

	if err := survey.Ask(qs, &answers, survey.WithIcons(func(set *survey.IconSet) {
		set.Question.Text = ">>"
	})); err != nil {
		return err
	}

	m.Remote = answers.Remote
	m.User.Login = answers.User
	m.User.Password = answers.Password

	return nil
}

func AskEntry(m *models.Entry) error {
	var name, entryType, entryValue string
	var err error

	name, err = askOne(&survey.Input{Message: "name:"})
	if err != nil {
		return err
	}

	entryType, err = askOne(&survey.Select{
		Message: "entry type:",
		Options: []string{
			models.EntryTypeText,
			models.EntryTypeFile,
			models.EntryTypeImage,
		},
		Default: models.EntryTypeText,
	})
	if err != nil {
		return err
	}

	if entryType == models.EntryTypeText {
		entryValue, err = askOne(&survey.Input{Message: "value:"})
		if err != nil {
			return err
		}
	} else {
		entryValue, err = AskFile("file:")
		if err != nil {
			return err
		}
	}

	m.Value = entryValue
	m.Name = name
	m.EntryType = entryType

	return nil
}

func AskConfirm(message string, answer *bool) error {
	prompt := &survey.Confirm{
		Message: message,
	}

	return survey.AskOne(prompt, answer)
}

func askOne(prompt survey.Prompt) (string, error) {
	output := ""
	err := survey.AskOne(prompt, &output, survey.WithValidator(survey.Required),
		survey.WithIcons(func(set *survey.IconSet) {
			set.Question.Text = ">>"
		}))

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

func AskCard(m *models.Card) (string, error) {
	qs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
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
				return models.IsValidLuhn(ans)
			},
		},
		{
			Name:   "month",
			Prompt: &survey.Input{Message: "Month:"},
			Validate: func(ans interface{}) error {
				v, _ := ans.(string)
				if len(v) != 2 {
					return fmt.Errorf("invalid month length")
				}

				return nil
			},
		},
		{
			Name:   "year",
			Prompt: &survey.Input{Message: "Year:"},
			Validate: func(ans interface{}) error {
				v, _ := ans.(string)
				if len(v) != 2 {
					return fmt.Errorf("invalid year length")
				}

				return nil
			},
		},
		{
			Name:   "cvv",
			Prompt: &survey.Input{Message: "CVV:"},
			Validate: func(ans interface{}) error {
				v, _ := ans.(string)
				if len(v) < 2 || len(v) > 4 {
					return fmt.Errorf("invalid cvv length")
				}

				return nil
			},
		},
	}

	answers := struct {
		Name   string `survey:"name"`
		Type   string `survey:"type"`
		Number string `survey:"number"`
		Month  string `survey:"month"`
		Year   string `survey:"year"`
		CVV    string `survey:"cvv"`
	}{}

	if err := survey.Ask(qs, &answers, survey.WithIcons(func(set *survey.IconSet) {
		set.Question.Text = ">>"
	})); err != nil {
		return "", err
	}

	m.Type = answers.Type
	m.Number = answers.Number
	m.Month = answers.Month
	m.Year = answers.Year
	m.CVV = answers.CVV

	return answers.Name, nil
}

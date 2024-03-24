package i18n

import (
	"embed"
	"encoding/json"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Languages map[string]Messages

type Messages map[string]string

var CurrentLanguage = language.English

//go:embed translations.json
var translations embed.FS

func LoadTranslations() (Languages, error) {
	data, err := translations.ReadFile("translations.json")
	if err != nil {
		return Languages{}, err
	}

	var languages Languages
	err = json.Unmarshal(data, &languages)
	if err != nil {
		return Languages{}, err
	}
	return languages, nil
}

func SetupTranslations(translations Languages) error {
	for langTag, messages := range translations {
		tag, err := language.Parse(langTag)
		if err != nil {
			return err
		}

		for key, translation := range messages {
			err := message.SetString(tag, key, translation)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Translate(key string) string {
	p := message.NewPrinter(CurrentLanguage)
	return p.Sprintf(key)
}

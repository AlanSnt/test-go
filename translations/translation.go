package translations

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizer *i18n.Localizer

func GetFilePath() string {
	filePath := os.Getenv("LOCALES_PATH")

	if filePath != "" {
		return filePath
	}

	filePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !strings.Contains(filePath, "/translations") {
		filePath = filePath + "/translations/locales"
	} else if !strings.Contains(filePath, "/locales") {
		filePath = filePath + "/locales"
	}

	// Find in path, if utils or middleware exists remove it
	for _, path := range []string{"utils", "middleware"} {
		if strings.Contains(filePath, path) {
			filePath = strings.Replace(filePath, path, "", 1)
		}
	}

	return filePath
}

func CreateBundle() *i18n.Bundle {
	filePath := GetFilePath()
	bundle := i18n.NewBundle(language.MustParse("fr-FR"))

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(filePath + "/active.fr-FR.json")
	bundle.MustLoadMessageFile(filePath + "/active.en-US.json")

	return bundle
}

func NewLocalizer(lang string, accept string) {
	bundle := CreateBundle()
	localizer = i18n.NewLocalizer(bundle, lang, accept)
}

func GetLocalizer() *i18n.Localizer {
	return localizer
}

func GetTranslation(MessageID string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: MessageID,
	})
}

func GetTranslationWithArgs(MessageID string, args map[string]interface{}) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    MessageID,
		TemplateData: args,
	})
}

func init() {
	// Default localizer
	NewLocalizer("fr-FR", "fr-FR")
}

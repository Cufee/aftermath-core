package localization

import (
	"embed"
	"encoding/json"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type LocalePrinter func(string) string

//go:embed resources
var resources embed.FS

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	dir, err := resources.ReadDir("resources")
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			bundle.LoadMessageFileFS(resources, "resources/"+file.Name())
		}
	}
}

func getLocalizer(code string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, code)
}

func GetPrinter(locale language.Tag) LocalePrinter {
	return func(s string) string {
		localized, err := getLocalizer(locale.String()).Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: s,
			},
		})
		if err != nil {
			log.Warn().Err(err).Msg("failed to localize string")
			return "? " + s
		}
		return localized
	}
}

package localization

import (
	"github.com/cufee/aftermath-core/internal/core/localization/resources"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type LocalePrinter func(string) string

func GetPrinter(locale language.Tag) LocalePrinter {
	return func(s string) string {
		localized, err := resources.GetLocalizer(locale.String()).Localize(&i18n.LocalizeConfig{
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

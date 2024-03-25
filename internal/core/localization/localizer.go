package localization

import (
	"github.com/cufee/aftermath-core/internal/core/localization/resources"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type LocalePrinter func(string) string

type SupportedLanguage struct {
	WargamingCode string
	Tag           language.Tag
}

/*
Languages supported by Wargaming
*/
var (
	LanguageEN SupportedLanguage = SupportedLanguage{WargamingCode: "en", Tag: language.English} // English
	LanguageRU SupportedLanguage = SupportedLanguage{WargamingCode: "ru", Tag: language.Russian} // Russian
	LanguagePL SupportedLanguage = SupportedLanguage{WargamingCode: "pl", Tag: language.Polish}  // Polish
	// LanguageDE   SupportedLanguage = SupportedLanguage{WargamingCode: "de", Tag: language.German}                // German
	// LanguageFR   SupportedLanguage = SupportedLanguage{WargamingCode: "fr", Tag: language.French}                // French
	LanguageES SupportedLanguage = SupportedLanguage{WargamingCode: "es", Tag: language.Spanish} // Spanish
	// LanguageTR   SupportedLanguage = SupportedLanguage{WargamingCode: "tr", Tag: language.Turkish}               // Turkish
	// LanguageCS   SupportedLanguage = SupportedLanguage{WargamingCode: "cs", Tag: language.Czech}                 // Czech // Thai
	// LanguageTH   SupportedLanguage = SupportedLanguage{WargamingCode: "th", Tag: language.Thai}                  // Thai
	// LanguageKO   SupportedLanguage = SupportedLanguage{WargamingCode: "ko", Tag: language.Korean}                // Korean
	// LanguageVI   SupportedLanguage = SupportedLanguage{WargamingCode: "vi", Tag: language.Vietnamese}            // Vietnamese
	// LanguageZhCH SupportedLanguage = SupportedLanguage{WargamingCode: "zh-cn", Tag: language.SimplifiedChinese}  // Simplified Chinese
	// LanguageZhTW SupportedLanguage = SupportedLanguage{WargamingCode: "zh-tw", Tag: language.TraditionalChinese} // Traditional Chinese
)

func GetPrinter(locale SupportedLanguage) LocalePrinter {
	return func(s string) string {
		localized, err := resources.GetLocalizer(locale.Tag.String()).Localize(&i18n.LocalizeConfig{
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

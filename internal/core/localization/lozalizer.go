package localization

import (
	"embed"
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
	// LanguagePL   SupportedLanguage = SupportedLanguage{WargamingCode: "pl", Tag: language.Polish}                // Polish
	// LanguageDE   SupportedLanguage = SupportedLanguage{WargamingCode: "de", Tag: language.German}                // German
	// LanguageFR   SupportedLanguage = SupportedLanguage{WargamingCode: "fr", Tag: language.French}                // French
	// LanguageES   SupportedLanguage = SupportedLanguage{WargamingCode: "es", Tag: language.Spanish}               // Spanish
	// LanguageTR   SupportedLanguage = SupportedLanguage{WargamingCode: "tr", Tag: language.Turkish}               // Turkish
	// LanguageCS   SupportedLanguage = SupportedLanguage{WargamingCode: "cs", Tag: language.Czech}                 // Czech // Thai
	// LanguageTH   SupportedLanguage = SupportedLanguage{WargamingCode: "th", Tag: language.Thai}                  // Thai
	// LanguageKO   SupportedLanguage = SupportedLanguage{WargamingCode: "ko", Tag: language.Korean}                // Korean
	// LanguageVI   SupportedLanguage = SupportedLanguage{WargamingCode: "vi", Tag: language.Vietnamese}            // Vietnamese
	// LanguageZhCH SupportedLanguage = SupportedLanguage{WargamingCode: "zh-cn", Tag: language.SimplifiedChinese}  // Simplified Chinese
	// LanguageZhTW SupportedLanguage = SupportedLanguage{WargamingCode: "zh-tw", Tag: language.TraditionalChinese} // Traditional Chinese
)

//go:embed resources/*.json
var resources embed.FS

var localizer *i18n.Localizer
var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	dir, err := resources.ReadDir("resources")
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if !file.IsDir() {
			bundle.LoadMessageFileFS(resources, "resources/"+file.Name())
		}
	}
	localizer = i18n.NewLocalizer(bundle,
		LanguageEN.Tag.String(),
		// LanguageRU.Tag.String(),
		// LanguagePL.Tag.String(),
		// LanguageDE.Tag.String(),
		// LanguageFR.Tag.String(),
		// LanguageES.Tag.String(),
		// LanguageTR.Tag.String(),
		// LanguageCS.Tag.String(),
		// LanguageTH.Tag.String(),
		// LanguageKO.Tag.String(),
		// LanguageVI.Tag.String(),
		// LanguageZhCH.Tag.String(),
		// LanguageZhTW.Tag.String(),
	)
}

func GetPrinter(locale SupportedLanguage) LocalePrinter {
	return func(s string) string {
		localized, err := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: s,
			},
		})
		if err != nil {
			log.Warn("failed to localize string: %s\n", s)
			return "? " + s
		}
		return localized
	}
}

package resources

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed *.json
var resources embed.FS

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	dir, err := resources.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if !file.IsDir() {
			bundle.LoadMessageFileFS(resources, file.Name())
		}
	}
}

func GetLocalizer(code string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, code)
}

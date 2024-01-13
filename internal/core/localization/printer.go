package localization

func GetPrinter(locale SupportedLanguage) func(string) string {
	return func(key string) string {
		return key
	}
}

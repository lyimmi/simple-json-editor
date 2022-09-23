package lang

var text map[string]map[string]string

func init() {
	text = map[string]map[string]string{
		"en": {
			"file":          "File",
			"file.new":      "New",
			"file.open":     "Open",
			"file.save":     "Save",
			"file.quit":     "Quit",
			"language":      "Language",
			"language.hu":   "Magyar",
			"language.en":   "English",
			"view":          "View",
			"view.darkmode": "Dark mode",
		},
		"hu": {
			"file":          "Fájl",
			"file.new":      "Új",
			"file.open":     "Megnyitás",
			"file.save":     "Mentés",
			"file.quit":     "Bezár",
			"language":      "Nyelv",
			"language.hu":   "Magyar",
			"language.en":   "English",
			"view":          "Nézet",
			"view.darkmode": "Sötét mód",
		},
	}
}

func Text(lang, textCode string) string {
	if txt, ok := text[lang][textCode]; ok {
		return txt
	}
	return ""
}

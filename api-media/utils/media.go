package utils

import "strings"

func NameMatchLocale(user_locale *Locale, lang string) bool {
	name_locale := strings.ToLower(lang)
	return lang == "__" ||
		name_locale == user_locale.Lang ||
		name_locale == user_locale.Region ||
		((user_locale.Lang == "gb" || user_locale.Region == "gb") && name_locale == "gb")
}

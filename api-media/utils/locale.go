package utils

import (
	"context"
	"log"
	"regexp"
	"strings"

	"google.golang.org/grpc/metadata"
)

type Locale struct {
	Lang, Region string
}

const (
	COOKIE_locale = "locale"
)

func NameMatchLocale(user_locale *Locale, lang string) bool {
	name_locale := strings.ToLower(lang)
	return lang == "__" ||
		name_locale == user_locale.Lang ||
		name_locale == user_locale.Region ||
		((user_locale.Lang == "gb" || user_locale.Region == "gb") && name_locale == "gb")
}

// var Locales []string = []string{"en", "fr"}
var LocaleMatcher = regexp.MustCompile("(?i)^(fr|en)(-\\w+)?$")

func GetLocale(ctx context.Context) Locale {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Locale{Lang: "en"}
	}

	cookie_headers := md.Get("grpcgateway-cookie")
	if len(cookie_headers) != 1 {
		return Locale{Lang: "en"}
	}
	raw_cookies := cookie_headers[0]
	cookies := strings.Split(raw_cookies, ";")
	if len(cookie_headers) == 0 {
		return Locale{Lang: "en"}
	}

	for _, raw_cookie := range cookies {
		cookie := strings.TrimSpace(raw_cookie)
		raw_cookie := strings.SplitN(cookie, "=", 2)
		if len(raw_cookie) != 2 {
			continue
		}

		key := raw_cookie[0]
		value := strings.TrimSpace(raw_cookie[1])

		if key == COOKIE_locale {
			matches := LocaleMatcher.FindStringSubmatch(value)
			if len(matches) >= 2 {
				locale := Locale{Lang: strings.ToLower(matches[1])}
				if len(matches) == 3 {
					locale.Region = matches[2]
				}
				return locale
			}
			log.Println("Invalid locale cookie '" + value + "' default to en")

			return Locale{Lang: "en"}
		}
	}

	return Locale{Lang: "en"}
}

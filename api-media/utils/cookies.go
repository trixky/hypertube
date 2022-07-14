package utils

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/trixky/hypertube/api-media/sanitizer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Locale struct {
	Lang, Region string
}

const (
	COOKIE_locale = "locale"
)

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
	log.Println(cookies)

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

const (
	COOKIE_token = "token"
)

var (
	COOKIE_ERROR_token_missing   = status.Errorf(codes.PermissionDenied, "token cookie missing")
	COOKIE_ERROR_token_corrupted = status.Errorf(codes.InvalidArgument, "token cookie corrupted")
)

func ExtractSanitizedTokenFromGrpcGatewayCookies(token string, ctx context.Context) (string, error) {
	if len(token) != 0 {
		if err := sanitizer.SanitizeToken(token); err != nil {
			return "", COOKIE_ERROR_token_corrupted
		}
		return token, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", COOKIE_ERROR_token_missing
	}

	cookie_headers := md.Get("grpcgateway-cookie")

	if len(cookie_headers) != 1 {
		return "", COOKIE_ERROR_token_missing
	}

	cookies := cookie_headers[0]

	splitted_cookies := strings.Split(cookies, ";")

	if len(cookie_headers) == 0 {
		return "", COOKIE_ERROR_token_missing
	}

	for _, splitted_cookie := range splitted_cookies {
		cookie := strings.TrimSpace(splitted_cookie)

		splitted_cookie := strings.SplitN(cookie, "=", 2)

		if len(splitted_cookie) != 2 {
			return "", COOKIE_ERROR_token_missing
		}

		cookie_key := splitted_cookie[0]
		cookie_value := splitted_cookie[1]

		if cookie_key == COOKIE_token {
			if err := sanitizer.SanitizeToken(cookie_value); err != nil {
				return "", COOKIE_ERROR_token_corrupted
			}
			return cookie_value, nil
		}
	}

	return "", COOKIE_ERROR_token_missing
}

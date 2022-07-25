package utils

import (
	"context"
	"strings"

	"github.com/trixky/hypertube/.shared/sanitizer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	COOKIE_token     = "token"
	COOKIE_user_info = "userInfo"
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

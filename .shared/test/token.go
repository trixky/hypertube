package test

// InvalidateToken corrupts a token
func CorrupteToken(token string) string {
	return token + "@"
}

// InvalidateToken invalidates a token
func InvalidateToken(token string) string {
	if len(token) > 2 {
		switch token[0] {
		case '0':
			token = "1" + token[1:]
		default:
			token = "0" + token[1:]
		}
	}

	return token
}

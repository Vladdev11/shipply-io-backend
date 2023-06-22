package util

func IsValidPassword(password string) bool {
	// Password must be at least 8 characters, 1 uppercase, 1 lowercase, 1 number, 1 special character
	var (
		hasUpperCase   = false
		hasLowerCase   = false
		hasNumber      = false
		hasSpecialChar = false
	)

	if len(password) < 8 {
		return false
	}

	for _, ch := range password {
		switch {
		case 'A' <= ch && ch <= 'Z':
			hasUpperCase = true
		case 'a' <= ch && ch <= 'z':
			hasLowerCase = true
		case '0' <= ch && ch <= '9':
			hasNumber = true
		case ch == '!' || ch == '@' || ch == '#' || ch == '$' || ch == '%' || ch == '^' || ch == '&':
			hasSpecialChar = true
		}
	}

	return hasUpperCase && hasLowerCase && hasNumber && hasSpecialChar
}

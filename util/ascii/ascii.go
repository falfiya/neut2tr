package unused

func IsControlCharacter(c byte) bool {
	return c < ' ' || c == 127
}

func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsLetter(c byte) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func IsAlphanumeric(c byte) bool {
	return IsDigit(c) || IsLetter(c)
}

func IsAlphanumericOr_(c byte) bool {
	return IsAlphanumeric(c) || c == '_'
}

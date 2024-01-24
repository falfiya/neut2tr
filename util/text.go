package util

func IndentOnly(s string) string {
	bytes := []byte(s)
	for i, char := range bytes {
		switch char {
		case '	', ' ':
		default:
			bytes[i] = ' '
		}
	}
	return string(bytes)
}

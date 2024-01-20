package util

func IndentOnly(s []byte) []byte {
	for i, char := range s {
		switch char {
		case '	', ' ':
		default:
			s[i] = ' '
		}
	}
	return s
}

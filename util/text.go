package util
import "strings"

func Reindent(s string, times int) string {
	indent := strings.Repeat("   ", times)
	lines := strings.Split(s, "\n")
	out := ""
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		out += indent + line + "\n"
	}
	return out
}

func IsControlCharacter(c byte) bool {
	return c <= ' ' || c == 127
}

// preserves tabs and newlines,
// ignores control characters
// replaces every other character with space
func WhitespaceOnly(s string) string {
	var sb strings.Builder
	for _, c := range s {
		if c == '\t' {
			sb.WriteByte('\t')
			continue
		}
		if c == '\n' {
			sb.WriteByte('\n')
			continue
		}
		if c == ' ' {
			sb.WriteByte(' ')
			continue
		}
		if IsControlCharacter(byte(c)) {
			continue
		}

		sb.WriteByte(' ')
	}
	return sb.String()
}

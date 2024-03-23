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

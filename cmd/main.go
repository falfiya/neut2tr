package main

import (
	"fmt"
	"neut2tr/neut"
	"neut2tr/racket"
	"neut2tr/rewrite"
	"neut2tr/util"
	"os"
	"strings"
)

func main() {
	n := os.Args[0]
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp(n)
		os.Exit(1)
	}
	if len(args) > 3 {
		printHelp(n)
		os.Exit(1)
	}
	switch args[0] {
	case "/?", "help", "-help", "--help":
		printHelp(n)
		os.Exit(0)
	}

	if args[0][0] == '-' {
		fmt.Fprintf(os.Stderr, "Unknown flag '%s'\n", args[0])
		printHelp(n)
		os.Exit(1)
	}

	bytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading '%s': %v\n", args[0], err)
		os.Exit(1)
	}

	text := string(bytes)
	if len(args) == 1 {
		os.Stdout.WriteString(translate(text))
	} else {
		var outfile string
		if len(args) == 3 {
			switch args[1] {
			case "-o", "--out":
			default:
				fmt.Fprintf(os.Stderr, "Unknown flag '%s'\n", args[1])
				printHelp(n)
				os.Exit(1)
			}
			outfile = args[2]
		} else {
			outfile = args[1]
		}
		os.WriteFile(outfile, []byte(translate(text)), 0777)
	}
}

func printHelp(n string) {
	fmt.Print("" +
		"neut2tr: Northeastern University Type Comments to typed/racket\n" +
		"Usage:\n"+
		"   " + n + " /?\n" +
		"   " + n + " help\n" +
		"   " + n + " -help\n" +
		"   " + n + " --help\n" +
		"      Prints this message\n" +
		"   " + n + " filename.rkt\n" +
		"      Converts NEU Type Comments to typed/racket and prints the result to stdout\n" +
		"   " + n + " filename.rkt output.rkt\n" +
		"   " + n + " filename.rkt -o output.rkt\n" +
		"   " + n + " filename.rkt --out output.rkt\n" +
		"      Converts NEU Type Comments to typed/racket and writes to output.rkt\n")
}

// prints errors to stderr
func translate(s string) string {
	s = strings.TrimPrefix(s, "#lang racket\n")
	comments, racketError := racket.ExtractComments(s)
	if racketError != nil {
		fmt.Fprintf(os.Stderr, "Racket Parsing Error: %s\nBailing out...", racketError.Msg)
		os.Exit(1)
	}
	offset := 0
	var sb strings.Builder
	sb.WriteString("#lang typed/racket\n")
	for _, comment := range comments {
		sb.WriteString(s[offset:comment.End()])
		offset = comment.End()
		meaningfulNodes, tokenizerError := neut.Parse(comment.Text)
		if tokenizerError != nil {
			fmt.Fprintf(os.Stderr, "Racket Tokenizer Error: %s\n", tokenizerError.Msg)
			continue
		}
		indentation := util.WhitespaceOnly(s[comment.LineAt:comment.Offset])
		for _, syntaxNode := range meaningfulNodes {
			sb.WriteString(indentation)
			sb.WriteString(rewrite.Rewrite(syntaxNode))
			sb.WriteByte('\n')
		}
	}
	sb.WriteString(s[offset:])
	return sb.String()
}


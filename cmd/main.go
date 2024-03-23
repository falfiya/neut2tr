package main

import (
	"fmt"
	"log"
	"neut2tr/neut"
	"neut2tr/neut2tr"
	"neut2tr/racket"
	"os"
)

func main() {
	testCommentFinder()
	// testParser()
}

func testCommentFinder() {
	dat, err := os.ReadFile("examples/homework_10.rkt")
	file := string(dat)
	if err != nil {
		log.Fatal(err)
	}
	comments, err2 := racket.ExtractComments(file)
	if err2 != nil {
		log.Fatal(err)
	}
	for _, c := range comments {
		testParser(c.Text)
	}
}

func testParser(s string) {
	meaningful, err := neut.Parse(s)
	if err != nil {
		fmt.Print(err.Msg)
	}
	for _, m := range meaningful {
		sel := m.SelF()
		fmt.Printf("%s\n", s[sel.Offset:sel.End()])
		fmt.Printf("   %s\n", neut2tr.Rewrite(m))
	}
}


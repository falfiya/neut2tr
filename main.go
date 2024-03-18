package main

import (
	"fmt"
	"log"
	"neut2tr/neut"
	"neut2tr/racket"
	"os"
)

func main() {
	testCommentFinder()
	// testParser()
}

func testCommentFinder() {
	dat, err := os.ReadFile("examples/examples.rkt")
	file := string(dat)
	if err != nil {
		log.Fatal(err)
	}
	comments, err2 := racket.ExtractComments(file)
	if err2 != nil {
		log.Fatal(err)
	}
	for i, c := range comments {
		fmt.Printf("$$$$$$$$$$$$$$$$$ %3d $$$$$$$$$$$$$$$$$\n%s\n", i, c.Text)
		testParser(c.Text)
	}
}

func testParser(s string) {
	meaningful, err := neut.Parse(s)
	if err != nil {
		fmt.Print(err.Msg)
	}
	for _, m := range meaningful {
		fmt.Print(m.Print())
	}
}


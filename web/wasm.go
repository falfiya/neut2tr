package main

import "syscall/js"

func main() {
	js.Global().Set("handleClick", js.FuncOf(handleClick))
}

func handleClick(this js.Value, inputs []js.Value) interface{} {
	println("Button clicked!")
	return nil
}

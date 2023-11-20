// Basic Go example compiled to WASM.
package main

import (
	"fmt"
	"os"
)

// Show how to import functions from the host using go:wasmimport
// This is portable between go and tinygo.

//go:wasmimport env log_i32
func logInt(i int32)

//go:wasmimport env log_string
func logString(s string)

func main() {
	logInt(42)
	logString("testtttt")
	fmt.Println("goenv environment:")

	for _, e := range os.Environ() {
		fmt.Println(" ", e)
	}
}

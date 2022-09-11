// It's internal tests and experiments
// only for learning :-)
package main

import (
	_ "embed"
	"fmt"
)

//go:generate ls -l
var src string

func main() {
	fmt.Printf("file content:\n\"%v\"\n", src)
}

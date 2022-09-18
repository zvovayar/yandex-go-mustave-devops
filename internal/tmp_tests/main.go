// It's internal tests and experiments
// only for learning :-)
package main

import (
	_ "embed"
	"fmt"
)

//go:generate ls -l
//go:embed ls.txt
var src string

func main() {
	fmt.Printf("file content:\n\"%v\"\n", src)
}

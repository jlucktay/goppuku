package main

import (
	"fmt"
	"os"

	"go.jlucktay.dev/goppuku/cmd"
)

func main() {
	if err := cmd.Run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)

		return
	}
}

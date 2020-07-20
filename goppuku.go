package main

import (
	"fmt"
	"os"

	"go.jlucktay.dev/goppuku/cmd"
)

func main() {
	if err := cmd.Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
}

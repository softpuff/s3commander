package main

import (
	"fmt"
	"os"

	"github.com/softpuff/s3commander/cmd"
)

func main() {
	if err := cmd.S3CommanderCMD.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

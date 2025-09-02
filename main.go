package main

import (
	"os"

	"github.com/github/gh-runtime-cli/cmd"
)

func main() {
	code := cmd.Execute()
	os.Exit(int(code))
}

package main

import (
	"flag"
	"os"
)

var (
	LinkwardenToken = os.Getenv("LINWEDEN_TOKEN")
)

func main() {
	flags := flag.NewFlagSet("migrate-json-export", flag.ExitOnError)
	flags.String()
	flag.Parse()
}

package main

import (
	"flag"
	"fmt"
	"github.com/mightyguava/jl"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var formatFlag = flag.String("format", "compact", `Formatter for logs. The options are "compact" and "logfmt"`)
	flag.Parse()
	var printer jl.EntryPrinter
	switch *formatFlag {
	case "logfmt":
		printer = jl.NewLogfmtPrinter(os.Stdout)
	case "compact":
		printer = jl.NewCompactPrinter(os.Stdout)
	}
	return jl.NewParser(os.Stdin, printer).Consume()
}

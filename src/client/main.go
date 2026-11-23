package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/apimgr/timezones/src/data"
	"github.com/apimgr/timezones/src/timezones"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Print matching timezone records as JSON")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: timezones-cli [--json] <name, abbreviation, city, or IANA identifier>\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	query := strings.TrimSpace(flag.Arg(0))
	if query == "" || len(query) > 128 {
		fmt.Fprintln(os.Stderr, "query must contain between 1 and 128 characters")
		os.Exit(2)
	}

	content, err := data.ReadFile("timezones.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read embedded timezone data: %v\n", err)
		os.Exit(1)
	}
	service, err := timezones.NewService(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load embedded timezone data: %v\n", err)
		os.Exit(1)
	}

	matches := service.Search(query)
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "no timezone matched %q\n", query)
		os.Exit(1)
	}
	if *jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(matches); err != nil {
			fmt.Fprintf(os.Stderr, "write JSON output: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, zone := range matches {
		fmt.Printf("%s (%s) UTC%+.2f — %s\n", zone.Value, zone.Abbr, zone.Offset, zone.Text)
		for _, identifier := range zone.UTC {
			fmt.Printf("  %s\n", identifier)
		}
	}
}

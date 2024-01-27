package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ParfenovVS/tagtolibv"
)

type tagToLibResult struct {
	Success bool
	Tag     string
	LibVer  string
	Reason  string
}

func NewTagToLibResult(success bool, tag string, libVer string, reason string) (*tagToLibResult, error) {
	ttlr := &tagToLibResult{
		Success: success,
		Tag:     tag,
		LibVer:  libVer,
		Reason:  reason,
	}

	return ttlr, nil
}

var supportedFormats = []string{
	"json", "md",
}

func main() {
	wd := flag.String("path", "", "Path to repository")
	lib := flag.String("lib", "", "Library name to be found")
	lim := flag.Int("limit", 1, "Number of tags to check (from latest to oldest)")
	format := flag.String("format", "json", "Output format (json / md)")

	flag.Parse()

	if len(*lib) == 0 {
		log.Fatal("Library name must be defined using flag --lib <name>")
	}

	if *lim <= 0 {
		log.Fatal("Limit cannot be <= 0")
	}

	formatFound := false
	for _, f := range supportedFormats {
		if *format == f {
			formatFound = true
			break
		}
	}
	if !formatFound {
		log.Fatal("Supported formats: json, md")
	}

	if len(*wd) != 0 {
		os.Chdir(*wd)
	}

	tags, err := tagtolibv.GetTags(*lim)
	if err != nil {
		log.Fatalf("Cannot get tags: %s", err.Error())
	}

	currentBranch, err := tagtolibv.GetCurrentBranch()
	if err != nil {
		log.Fatalf("Cannot get current branch: %s", err.Error())
	}
	defer tagtolibv.GitCheckout(currentBranch)

	result := []tagToLibResult{}
	for _, t := range tags {
		var ver string
		err := tagtolibv.GitCheckout(t)
		if err == nil {
			ver, err = tagtolibv.GetLibVersion(t, *lib)
		}
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		ttlr, _ := NewTagToLibResult(err == nil, t, ver, errMsg)
		result = append(result, *ttlr)
	}

	switch *format {
	case "json":
		j, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatal(string(j))
		}

		fmt.Println(string(j))
	case "md":
		fmt.Printf("| Tag | %s version |\n", *lib)
		fmt.Println("| --- | --- |")
		for _, r := range result {
			if r.Success {
				fmt.Printf("| %s | %s |\n", r.Tag, r.LibVer)
			} else {
				fmt.Printf("| %s | **fail:** %s |\n", r.Tag, r.Reason)
			}
		}
	}

}

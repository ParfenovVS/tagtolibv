package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ParfenovVS/tagtolibv"
)

func main() {
	wd := flag.String("p", "", "Path to repository")
	lib := flag.String("l", "", "Library name to be found")

	flag.Parse()

	if len(*wd) != 0 {
		os.Chdir(*wd)
	}

	currentBranch, err := tagtolibv.GetCurrentBranch()
	if err != nil {
		log.Fatalf("Cannot get current branch: %s", err.Error())
	}
	defer tagtolibv.GitCheckout(currentBranch)

	tags, err := tagtolibv.GetTags()
	if err != nil {
		log.Fatalf("Cannot get tags: %s", err.Error())
	}

	result := make(map[string]string)
	for _, t := range tags {
		err := tagtolibv.GitCheckout(t)
		if err != nil {
			log.Fatalf("Cannot checkout %s: %s", t, err.Error())
		}

		ver, err := tagtolibv.GetLibVersion(t, *lib)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("cannot parse lib version for tag %s: %s", t, err.Error()))
		} else {
			result[t] = ver
		}

	}

	j, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(string(j))
	}

	fmt.Println(string(j))
}

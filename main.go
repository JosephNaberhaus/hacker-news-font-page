package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/JosephNaberhaus/hacker-news-font-page/internal/dataset"
	"github.com/JosephNaberhaus/hacker-news-font-page/internal/queryer"
)

var (
	output    = flag.String("output", "titles.csv", "Output file. If this file already exists then the existing data will be re-used.")
	startYear = flag.Int("startYear", 2020, "The year to start querying data from.")
)

func main() {
	start := time.Date(*startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().AddDate(0, 0, -1)

	qry := &queryer.Queryer{
		WaitMillisecondsMin: 20000,
		WaitMillisecondsMax: 30000,
	}

	fmt.Printf("Reading dataset...\n")

	ds, err := dataset.New(*output)
	if err != nil {
		fmt.Printf("Failed to read dataset: %s\n", err.Error())
		os.Exit(10)
	}

	fmt.Printf("Adding missing entries...\n")

	err = ds.AddMissingEntries(start, end, qry)
	if err != nil {
		fmt.Printf("Failed to add missing dataset entries: %s\n", err.Error())
		// Even though we encountered an error we can still go ahead and save what we have.
	}

	fmt.Printf("Saving dataset...\n")

	err = ds.Save()
	if err != nil {
		fmt.Printf("Failed to save dataset: %s\n", err.Error())
		os.Exit(30)
	}

	fmt.Printf("Done\n")
}

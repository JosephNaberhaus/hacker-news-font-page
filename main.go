package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/JosephNaberhaus/hacker-news-font-page/internal/dataset"
	"github.com/JosephNaberhaus/hacker-news-font-page/internal/queryer"
)

var (
	output    = flag.String("output", "titles.csv", "Output file. If this file already exists then the existing data will be re-used.")
	startYear = flag.Int("startYear", 2010, "The year to start querying data from.")
	duration  = flag.Int("duration", 60, "How many seconds to run the program before stopping.")
)

func main() {
	flag.Parse()

	start := time.Date(*startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().AddDate(0, 0, -1)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*duration)*time.Second)
	defer cancel()

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

	err = ds.AddMissingEntries(ctx, start, end, qry)
	if err != nil {
		if errors.Is(err, context.Cause(ctx)) {
			fmt.Print("Context was cancelled\n")
		} else {
			fmt.Printf("Failed to add missing dataset entries: %s\n", err.Error())
			// Even though we encountered an error we can still go ahead and save what we have.
		}
	}

	fmt.Printf("Saving dataset...\n")

	err = ds.Save()
	if err != nil {
		fmt.Printf("Failed to save dataset: %s\n", err.Error())
		os.Exit(30)
	}

	fmt.Printf("Done\n")
}

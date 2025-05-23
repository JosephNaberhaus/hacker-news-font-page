package dataset

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"time"
)

type Dataset struct {
	filename string
	entries  map[entryDate][30]string
}

func New(filename string) (Dataset, error) {
	if !exists(filename) {
		return Dataset{
			filename: filename,
			entries:  make(map[entryDate][30]string),
		}, nil
	}

	datasetFile, err := os.Open(filename)
	if err != nil {
		return Dataset{}, fmt.Errorf("error opening dataset file: %w", err)
	}
	defer datasetFile.Close()

	reader := csv.NewReader(datasetFile)
	reader.FieldsPerRecord = 31

	records, err := reader.ReadAll()
	if err != nil {
		return Dataset{}, fmt.Errorf("error reading dataset file: %w", err)
	}
	if len(records) == 0 {
		return Dataset{}, fmt.Errorf("dataset is empty")
	}
	if !slices.Equal(records[0], createHeader()) {
		return Dataset{}, fmt.Errorf("unexpected header row in dataset")
	}

	entries := make(map[entryDate][30]string, len(records)-1)
	for i, record := range records[1:] {
		date, err := parseDate(record[0])
		if err != nil {
			return Dataset{}, fmt.Errorf("error parsing date at row %d: %w", i+1, err)
		}

		var titles [30]string
		for i := range 30 {
			titles[i] = record[i+1]
		}

		entries[date] = titles
	}

	return Dataset{
		filename: filename,
		entries:  entries,
	}, nil
}

func (d Dataset) AddMissingEntries(ctx context.Context, startTime, endTime time.Time, q Queryer) error {
	start := newEntryDateFromTime(startTime)
	end := newEntryDateFromTime(endTime)

	cur := start
	for cur.Compare(end) <= 0 {
		if _, ok := d.entries[cur]; !ok {
			titles, err := q.GetTitles(ctx, cur.Year, cur.Month, cur.Day)
			if err != nil {
				return err
			}

			d.entries[cur] = titles
		}

		cur = cur.nextDay()
	}

	return nil
}

func (d Dataset) Save() error {
	datasetFile, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error opening dataset file when saving: %w", err)
	}
	defer datasetFile.Close()

	writer := csv.NewWriter(datasetFile)
	err = writer.Write(createHeader())
	if err != nil {
		return fmt.Errorf("error writing header of dataset file: %w", err)
	}

	dates := slices.Collect(maps.Keys(d.entries))
	slices.SortFunc(dates, entryDate.Compare)

	for _, date := range dates {
		titles := d.entries[date]

		records := make([]string, 0, 31)
		records = append(records, date.String())
		for i := range 30 {
			records = append(records, titles[i])
		}

		err = writer.Write(records)
		if err != nil {
			return fmt.Errorf("error entry of dataset file %s: %w", date.String(), err)
		}
	}

	writer.Flush()

	return nil
}

func createHeader() []string {
	header := make([]string, 0, 31)
	header = append(header, "Date")
	for i := range 30 {
		header = append(header, fmt.Sprintf("Title %d", i+1))
	}

	return header
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

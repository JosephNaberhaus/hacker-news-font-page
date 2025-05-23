package dataset

import (
	"cmp"
	"fmt"
	"time"
)

type entryDate struct {
	Year, Month, Day int
}

func newEntryDateFromTime(t time.Time) entryDate {
	// Always switch to UTC so that the local time zone doesn't affect the behavior of the program.
	t = t.In(time.UTC)

	return entryDate{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

func parseDate(dateStr string) (entryDate, error) {
	t, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return entryDate{}, fmt.Errorf("error parsing date %s: %w", dateStr, err)
	}

	return newEntryDateFromTime(t), nil
}

func (e entryDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", e.Year, e.Month, e.Day)
}

func (e entryDate) nextDay() entryDate {
	t := time.Date(e.Year, time.Month(e.Month), e.Day, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, 1)

	return newEntryDateFromTime(t)
}

func (e entryDate) Compare(other entryDate) int {
	if e.Year != other.Year {
		return cmp.Compare(e.Year, other.Year)
	}

	if e.Month != other.Month {
		return cmp.Compare(e.Month, other.Month)
	}

	return cmp.Compare(e.Day, other.Day)
}

package dataset

import "context"

type Queryer interface {
	GetTitles(ctx context.Context, year, month, day int) ([30]string, error)
}

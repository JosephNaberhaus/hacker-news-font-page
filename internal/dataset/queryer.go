package dataset

type Queryer interface {
	GetTitles(year, month, day int) ([30]string, error)
}

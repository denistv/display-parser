package domain

type ModelEntity struct {
	URL    string `db:"url"`
	Brand  string `db:"brand"`
	Series string `db:"series"`
	Name   string `db:"name"`
	Year   int64  `db:"year"`
	Size   int64  `db:"size"`
	PPI    int64  `db:"ppi"`
}

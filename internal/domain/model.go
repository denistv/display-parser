package domain

type ModelEntity struct {
	URL   string `db:"url"`
	Brand string `db:"brand"`
	Name  string `db:"name"`
	PPI   int64  `db:"ppi"`
}

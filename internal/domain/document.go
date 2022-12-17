package domain

type DocumentEntity struct {
	URL  string `db:"url"`
	Body string `db:"body"`
}

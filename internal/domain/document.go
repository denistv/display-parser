package domain

type ModelDocument struct {
	URL  string `db:"url"`
	Body string `db:"body"`
}

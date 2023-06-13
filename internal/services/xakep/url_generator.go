package xakep

import "context"

func NewURLGenerator(from, to int) *URLGenerator {
	return &URLGenerator{
		from: from,
		to:   to,
	}
}

type URLGenerator struct {
	from int
	to   int
}

func (u *URLGenerator) Run(ctx context.Context) {

}

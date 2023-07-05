package domain

import (
	"errors"
	"regexp"
)

var entityIDRegexp = regexp.MustCompile("^(?P<schema>http|https)://(?P<domain>[A-z0-9.]+)/(?P<lang>[a-z]+)/model/(?P<id>[a-z0-9]+)$")

// EntityID Идентификатор сущности модели во внешней системе. Извлекаем его из URL
func EntityID(pageURL string) (string, error) {
	matches := entityIDRegexp.FindStringSubmatch(pageURL)
	if len(matches) == 0 {
		return "", errors.New("incorrect URL: no matches")
	}

	i := entityIDRegexp.SubexpIndex("id")

	if i == -1 {
		return "", errors.New("cannot get ID from page URL")
	}

	return matches[i], nil
}

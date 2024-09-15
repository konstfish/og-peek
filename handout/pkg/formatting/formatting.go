package formatting

import (
	"net/url"
	"regexp"
	"strings"
)

func DomainFromUrl(urlStr string) (string, error) {
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	return urlParsed.Host, nil
}

func UrlToSlug(urlStr string) (string, error) {
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := reg.ReplaceAllString(urlParsed.Path, "-")

	// removes leading `-`
	slug = strings.Trim(slug, "-")

	// safety for domains ending with just the tld, so `-` also serves as index
	if len(slug) == 0 {
		slug = "-"
	}

	return strings.ToLower(slug), nil
}

package formatting

import (
	"regexp"
	"strings"
)

func UrlToSlug(url string) string {
	url = strings.TrimPrefix(url, "https://")

	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	if idx := strings.Index(url, "#"); idx != -1 {
		url = url[:idx]
	}

	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := reg.ReplaceAllString(url, "-")

	slug = strings.Trim(slug, "-")
	return strings.ToLower(slug)
}

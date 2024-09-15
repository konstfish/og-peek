package main

import (
	"log"
	"net/url"
)

func main() {
	urlStr := "https://test.konst.fish/asdfsl#sldkfjs?asdf=asdf"
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	path := urlParsed.Path
	log.Println(path)
	log.Println(urlParsed.Host)
}

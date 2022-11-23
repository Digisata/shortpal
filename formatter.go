package main

import (
	"fmt"
	"os"
)

type UrlFormatter struct {
	ID       uint   `json:"id"`
	Url      string `json:"url"`
	ShortUrl string `json:"short_url"`
}

func FormatUrl(url Url) UrlFormatter {
	urlFormatter := UrlFormatter{
		ID:       url.ID,
		Url:      url.Url,
		ShortUrl: fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), url.ShortUrl[0].ShortUrl),
	}

	return urlFormatter
}

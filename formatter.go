package main

import (
	"fmt"
	"os"
)

type UrlFormatter struct {
	ID        uint   `json:"id"`
	OriginUrl string `json:"origin_url"`
	ShortUrl  string `json:"short_url"`
}

func FormatUrl(url Url) UrlFormatter {
	urlFormatter := UrlFormatter{
		ID:        url.ID,
		OriginUrl: url.OriginUrl,
		ShortUrl:  fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), url.ShortUrl),
	}

	return urlFormatter
}

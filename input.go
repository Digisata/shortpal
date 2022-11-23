package main

type ShortenUrlInput struct {
	Url string `json:"url" binding:"required"`
}

type GetUrlDetailInput struct {
	ShortUrl string `uri:"SHORTURL" binding:"required"`
}

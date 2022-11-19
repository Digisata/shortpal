package main

type ShortenUrlInput struct {
	OriginUrl string `json:"origin_url" binding:"required"`
}

type GetUrlDetailInput struct {
	ShortUrl string `uri:"SHORTURL" binding:"required"`
}

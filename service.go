package main

import (
	"github.com/Digisata/shortpal/helper"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s Service) ShortenUrl(input ShortenUrlInput) (Url, error) {
	url := Url{
		OriginUrl: input.OriginUrl,
	}

	shortUrl, err := helper.GenerateShortUrl(url.OriginUrl)
	if err != nil {
		return Url{}, err
	}

	url.ShortUrl = shortUrl

	newUrl, err := s.repository.UpdateOrCreate(url)
	if err != nil {
		return newUrl, err
	}

	return newUrl, nil
}

func (s Service) GetUrlByShortUrl(shortUrl string) (Url, error) {
	url, err := s.repository.FindByShortUrl(shortUrl)
	if err != nil {
		return url, err
	}

	return url, nil
}

package main

import (
	"github.com/Digisata/shortpal/helper"
)

type Service interface {
	ShortenUrl(input ShortenUrlInput) (Url, error)
	GetUrlByShortUrl(shortUrl string) (Url, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ShortenUrl(input ShortenUrlInput) (Url, error) {
	url := Url{
		OriginUrl: input.OriginUrl,
	}

	shortUrl, err := helper.GenerateShortUrl(url.OriginUrl)
	if err != nil {
		return Url{}, err
	}

	url.ShortUrl = shortUrl

	newUrl, err := s.repository.Save(url)
	if err != nil {
		return newUrl, err
	}

	return newUrl, nil
}

func (s *service) GetUrlByShortUrl(shortUrl string) (Url, error) {
	url, err := s.repository.FindByShortUrl(shortUrl)
	if err != nil {
		return url, err
	}

	return url, nil
}

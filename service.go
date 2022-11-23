package main

import (
	"context"

	"github.com/Digisata/shortpal/helper"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s Service) ShortenUrl(input ShortenUrlInput, ctx context.Context) (Url, error) {
	url, count, err := s.repository.FindUrl(input.Url)
	if err != nil {
		return url, err
	}

	if count == 0 {
		url, err = s.repository.CreateUrl(Url{Url: input.Url})
		if err != nil {
			return url, err
		}
	}

	shortLink, err := helper.GenerateShortLink(url.Url)
	if err != nil {
		return url, err
	}

	shortUrl := ShortUrl{
		IdUrl:    url.ID,
		ShortUrl: shortLink,
	}

	insertedShortUrl, err := s.repository.CreateShortUrl(shortUrl, ctx)
	if err != nil {
		return url, err
	}

	url.ShortUrl = append([]ShortUrl{}, insertedShortUrl)

	return url, nil
}

func (s Service) Redirect(shortLink string, ctx context.Context) (Url, error) {
	shortUrl, err := s.repository.FindShortUrl(shortLink, ctx)
	if err != nil {
		return Url{}, err
	}

	url, err := s.repository.FindUrlByID(shortUrl.IdUrl, ctx)
	if err != nil {
		return url, err
	}

	return url, nil
}

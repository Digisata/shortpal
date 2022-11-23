package main

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) (*Repository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}
	if rdb == nil {
		return nil, errors.New("rdb cannot be nil")
	}

	return &Repository{db: db, rdb: rdb}, nil
}

func (r Repository) CreateUrl(url Url) (Url, error) {
	err := r.db.Create(&url).Error
	if err != nil {
		return url, err
	}

	return url, nil
}

func (r Repository) CreateShortUrl(shortUrl ShortUrl, ctx context.Context) (ShortUrl, error) {
	err := r.db.Create(&shortUrl).Error
	if err != nil {
		return shortUrl, err
	}

	dataJson, err := json.Marshal(shortUrl)
	if err != nil {
		return shortUrl, err
	}

	err = r.rdb.Set(ctx, shortUrl.ShortUrl, string(dataJson), 10800*time.Second).Err()
	if err != nil {
		return shortUrl, err
	}

	return shortUrl, nil
}

func (r Repository) FindUrlByID(id uint, ctx context.Context) (Url, error) {
	var url Url

	val, err := r.rdb.Get(ctx, strconv.Itoa(int(id))).Result()
	if err == redis.Nil {
		err := r.db.First(&url, id).Error
		if err != nil {
			return url, err
		}

		dataJson, err := json.Marshal(url)
		if err != nil {
			return url, err
		}

		err = r.rdb.Set(ctx, strconv.Itoa(int(id)), string(dataJson), 10800*time.Second).Err()
		if err != nil {
			return url, err
		}

		return url, nil
	} else if err != nil {
		return url, err
	} else {
		err := json.Unmarshal([]byte(val), &url)
		if err != nil {
			return url, err
		}

		return url, nil
	}
}

func (r Repository) FindUrl(link string) (Url, int64, error) {
	var url Url

	result := r.db.Where("url = ?", link).First(&url)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return url, 0, result.Error
	}

	return url, result.RowsAffected, nil
}

func (r Repository) FindShortUrl(shortLink string, ctx context.Context) (ShortUrl, error) {
	var shortUrl ShortUrl

	val, err := r.rdb.Get(ctx, shortLink).Result()
	if err == redis.Nil {
		err = r.db.Where("short_url = ?", shortLink).First(&shortUrl).Error
		if err != nil {
			return shortUrl, err
		}

		dataJson, err := json.Marshal(shortUrl)
		if err != nil {
			return shortUrl, err
		}

		err = r.rdb.Set(ctx, shortLink, string(dataJson), 10800*time.Second).Err()
		if err != nil {
			return shortUrl, err
		}

		return shortUrl, nil
	} else if err != nil {
		return shortUrl, err
	} else {
		err := json.Unmarshal([]byte(val), &shortUrl)
		if err != nil {
			return shortUrl, err
		}

		return shortUrl, nil
	}
}

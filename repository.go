package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Repository interface {
	UpdateOrCreate(url Url) (Url, error)
	FindByShortUrl(shortUrl string) (Url, error)
}

type repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) *repository {
	return &repository{db: db, rdb: rdb}
}

func (r *repository) UpdateOrCreate(url Url) (Url, error) {
	if r.db.Model(&url).Where("short_url = ?", url.ShortUrl).Updates(&url).RowsAffected == 0 {
		err := r.db.Create(&url).Error
		if err != nil {
			return url, err
		}
	}

	return url, nil
}

func (r *repository) FindByShortUrl(shortUrl string) (Url, error) {
	var url Url
	var ctx = context.Background()

	val, err := r.rdb.Get(ctx, shortUrl).Result()
	if err == redis.Nil {
		err = r.db.Where("short_url = ?", shortUrl).First(&url).Error
		if err != nil {
			return url, err
		}

		dataJson, err := json.Marshal(url)
		if err != nil {
			return url, err
		}

		err = r.rdb.Set(ctx, shortUrl, string(dataJson), 10800*time.Second).Err()
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

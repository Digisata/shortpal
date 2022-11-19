package main

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	OriginUrl string
	ShortUrl  string `gorm:"uniqueIndex"`
}

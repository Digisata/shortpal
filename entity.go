package main

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	Url      string
	ShortUrl []ShortUrl `gorm:"foreignKey:IdUrl"`
}

type ShortUrl struct {
	gorm.Model
	IdUrl    uint
	ShortUrl string
}

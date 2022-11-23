package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	/* err := godotenv.Load()
	if err != nil {
		panic("failed to load env")
	} */

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName, dbPort)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
		},
	)
	db, errCon := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if errCon != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Url{}, &ShortUrl{})

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RDB_HOST"),
		Password: os.Getenv("RDB_PASSWORD"),
		DB:       0,
	})

	repository, err := NewRepository(db, rdb)
	if err != nil {
		panic(err)
	}
	service := NewService(repository)
	handler := NewHandler(service)

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST"}
	router.Use(cors.New(config))

	router.GET("/:SHORTURL", handler.Redirect)

	api := router.Group("/api/v1")

	api.POST("/url", handler.ShortenUrl)

	router.Run()
}

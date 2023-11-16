package main

import (
	// "fmt"

	"todo/todolist"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	gorm1 "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize Dependencies
	// Service Port, Database, Logger, Cache etc.
	router := gin.Default()

	dsn := "postgres://" + "postgres" + ":" + "12345678" + "@" + "localhost" + ":" + "5432" + "/" + "todo"
	db, err := gorm1.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	//redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("REDIS_CLIENT", rdb)
		// Other sets...
	})


	// Boostrap services
	registerSvc := &todolist.HandlerService{}
	registerSvc.Bootstrap(router)

	port := "3000"
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

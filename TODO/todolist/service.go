package todolist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	gorm1 "github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	_ "github.com/rs/zerolog/log"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api/v1/")
	qrg.GET("todo", hs.GetAllList)
	qrg.POST("todo", hs.PostTodo)
	qrg.POST("temptodo", hs.TempTodo)
	qrg.GET("temptodo", hs.GetAllTempList)
	qrg.DELETE("todo/:id", hs.DeleteTodo)
	qrg.POST("addtion", hs.AddTwoNumber)
	qrg.POST("division", hs.Division)

}

func (hs *HandlerService) PostTodo(c *gin.Context) {

	// db := c.MustGet("DB").(*gorm.DB)
	// _ = db
	// fmt.Println("0-0-0-0-0-0-0-0-")
	dsn := "postgres://" + "postgres" + ":" + "12345678" + "@" + "localhost" + ":" + "5432" + "/" + "todo"
	// log.Info().Msg(dsn)
	db, err := gorm1.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	var todo TodoList
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	fmt.Println("todotodotodo", todo)
	todo.CreatedAt = time.Now()
	if err := db.Debug().Table("todo_list").Create(&todo).Error; err != nil {
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Actor Created Successfully."})
		return
	}

}

func (hs *HandlerService) GetAllList(c *gin.Context) {
	dsn := "postgres://" + "postgres" + ":" + "12345678" + "@" + "localhost" + ":" + "5432" + "/" + "todo"
	// log.Info().Msg(dsn)
	db, err := gorm1.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	// db := c.MustGet("DB").(*gorm.DB)

	var todo []TodoList
	if err := db.Debug().Table("todo_list").Find(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "erroresponse")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func (hs *HandlerService) DeleteTodo(c *gin.Context) {
	// db := c.MustGet("DB").(*gorm.DB)

	dsn := "postgres://" + "postgres" + ":" + "12345678" + "@" + "localhost" + ":" + "5432" + "/" + "todo"
	// log.Info().Msg(dsn)
	db, err := gorm1.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	todoId := c.Param("id")
	var todo TodoList
	if err := db.Debug().Table("todo_list").Where("id=?  ", todoId).Find(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, "Record not found")
		return
	}
	if err := db.Debug().Table("todo_list").Where("id=? ", todoId).Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Todo record Deleted Successfully.", "Status": http.StatusOK})
		return
	}
}

func (hs *HandlerService) TempTodo(c *gin.Context) {

	var todo TodoList
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	todo.CreatedAt = time.Now()
	todo.Id = uuid.NewString()
	redisData := make(map[string]interface{})
	redisData["data"] = todo
	m, _ := json.Marshal(redisData)
	rediserr := PostRedisDataWithKey(c, todo.Id, m)
	if rediserr != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Expired Todo record Created Successfully.", "Status": http.StatusOK})
		return
	}

}

func redisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}

func PostRedisDataWithKey(c *gin.Context, key string, data []byte) error {
	rdb := redisClient()
	// ctx := context.Background()
	ttl := 2
	var ttlDuration time.Duration
	if ttl != 0 {
		ttlDuration = time.Duration(ttl) * time.Minute
	} else {
		ttlDuration = 2 * time.Minute
	}
	err := rdb.Set(key, data, ttlDuration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (hs *HandlerService) GetAllTempList(c *gin.Context) {
	rdb := redisClient()
	keys, err := rdb.Keys("*").Result()
	if err != nil {
		return
	}
	var todos []map[string]interface{}
	for _, key := range keys {
		data, err := rdb.Get(key).Result()
		if err != nil {
			return
		}
		var redisData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &redisData); err != nil {
			return
		}
		fmt.Println("redisDataredisData", redisData)
		// var todo map[string]interface{}
		if data, ok := redisData["data"].(map[string]interface{}); ok {
			_ = data
			todos = append(todos, redisData["data"].(map[string]interface{}))
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func (hs *HandlerService) AddTwoNumber(c *gin.Context) {

	// db := c.MustGet("DB").(*gorm.DB)
	// _ = db
	// fmt.Println("0-0-0-0-0-0-0-0-")
	
	var twoNumbers Numbers
	if err := c.ShouldBindJSON(&twoNumbers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	sum := twoNumbers.Number1 + twoNumbers.Number2
		
	c.JSON(http.StatusOK, gin.H{"Sum of Two Numbers": sum})


}

func (hs *HandlerService) Division(c *gin.Context) {

	var twoNumbers Division
	if err := c.ShouldBindJSON(&twoNumbers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	fmt.Println("twoNumbers.Dividend" , twoNumbers.Dividend)
	fmt.Println("twoNumbers.Dividend" , twoNumbers.Divisor)
	sum := twoNumbers.Dividend / twoNumbers.Divisor
		
	c.JSON(http.StatusOK, gin.H{"value after Division": sum})

}
package todolist

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func TestDivision(t *testing.T) {
    router := gin.Default()

    router.POST("/api/v1/division", func(c *gin.Context) {
        var request struct {
            Dividend int `form:"dividend" binding:"required"`
            Divisor  int `form:"divisor" binding:"required"`
        }

        if err := c.ShouldBind(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }


        result := request.Dividend / request.Divisor
        c.JSON(http.StatusOK, gin.H{"result": result})
    })

    requestBody := []byte(`{"dividend": 22, "divisor": 3}`)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/division", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"result":7}`, w.Body.String())
}


func TestAdditionAPI(t *testing.T) {
    router := gin.Default()

    router.POST("/api/v1/addition", func(c *gin.Context) {
        var request struct {
            Number1 int `form:"number1" binding:"required"`
            Number2 int `form:"number2" binding:"required"`
        }

        if err := c.ShouldBind(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        result := request.Number1 + request.Number2
        c.JSON(http.StatusOK, gin.H{"result": result})
    })

    requestBody := []byte(`{"number1": 4, "number2": 3}`)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/addition", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"result":7}`, w.Body.String())
}
package todolist

import "time"

type TodoList struct {
	Id        string    `gorm:"primaryKey"`
	TodoText  string    `json:"todo_text"`
	CreatedAt time.Time `json:"created_at"`
}


type Numbers struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

type Division struct {
	Dividend int `json:"dividend"`
	Divisor int `json:"divisor"`
}
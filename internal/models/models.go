package models

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Age        string `json:"age"`
	Is_student bool   `json:"is_student"`
	Data       map[string]any
}

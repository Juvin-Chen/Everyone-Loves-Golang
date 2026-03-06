package model

type User struct {
	ID   int    `json:"id"`
	Name string `json:"username"`
	Age  int    `json:"age"`
}

package model

type Post struct {
	Id           int
	Title        string
	Views        int
	CreatedAt    string
	UpdatedAt    string
	CategoryId   int
	CategoryName string
	TagIds       []int
	Description  string
	Content      string
	TagNames     []string
	Status       int
}

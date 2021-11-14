package model

import "time"

type Post struct {
	Id           int
	Title        string
	Views        int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CategoryId   int
	CategoryName string
	TagIds       []int
	Description  string
	Content      string
	TagNames     []string
	Status       int
}

package model

type Tag struct {
	Id int
	Name string
	Count int
	CreatedAt string
	UpdatedAt string
}

type TagPost struct {
	Id int
	TagId int
	PostId int
}

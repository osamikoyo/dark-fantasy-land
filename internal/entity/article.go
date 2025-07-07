package entity

import "time"

type Article struct {
	Title     string    `bson:"title"`
	Topics    []string  `bson:"topics"`
	Timestamp time.Time `bson:"timestamp"`
	Content   string    `bson:"content"`
	Author    string    `bson:"author"`
}

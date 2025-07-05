package entity

import "time"

type Article struct {
	Topic     string    `bson:"topic"`
	Timestamp time.Time `bson:"timestamp"`
	Content   string    `bson:"content"`
	Author    string    `bson:"author"`
}

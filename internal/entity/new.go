package entity

import "time"

type New struct {
	Title     string    `bson:"title"`
	Topic     string    `bson:"topic"`
	Author    string    `bson:"author"`
	Censor    uint8     `bson:"censor"`
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"timestamp"`
}

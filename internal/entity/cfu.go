package entity

import "time"

type ContentFromUser struct {
	Content   string    `bson:"content"`
	MediaName string    `bson:"media_name"`
	Author    string    `bson:"author"`
	Timestamp time.Time `bson:"timestamp"`
}

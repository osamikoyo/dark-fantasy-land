package entity

import "time"

type Article struct {
	Title     string    `bson:"title" redis:"title" mapstructure:"title"`
	Topics    []string  `bson:"topics" redis:"topics" mapstructure:"topics"`
	Timestamp time.Time `bson:"timestamp" redis:"timestamp" mapstructure:"timestamp"`
	Content   string    `bson:"content" redis:"content" mapstructure:"content"`
	Author    string    `bson:"author" redis:"author" mapstructure:"author"`
}

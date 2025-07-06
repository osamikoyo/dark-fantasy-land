package entity

import "time"

type Mem struct{
	ImageName string `bson:"image_name"`
	Topics []string `bson:"topics"`
	Author string `bson:"author"`
	Timestamp time.Time `bson:"timestamp"`
	Description string `bson:"description"`
}
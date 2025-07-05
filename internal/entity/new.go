package entity

type New struct {
	Title   string `bson:"title"`
	Topic   string `bson:"topic"`
	Censor  uint8  `bson:"censor"`
	Content string `bson:"content"`
}

package entity

type Wallpaper struct {
	ImageName  string `bson:"image_name"`
	Topic      string `bson:"topic"`
	Resolution string `bson:"resolution"`
}

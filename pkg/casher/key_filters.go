package casher

import "fmt"

func newArticleKey(author, title string) string {
	return fmt.Sprintf("article:%s:%s", author, title)
}

func newWallpaperKey(imageName, author string) string {
	return fmt.Sprintf("wallpaper:%s:%s", imageName, author)
}

func newNewKey(title, author string) string {
	return fmt.Sprintf("new:%s:%s", author, title)
}

func newMemKey(imageName, author string) string {
	return fmt.Sprintf("mem:%s:%s", imageName, author)
}

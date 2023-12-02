package lib

import "strings"


func RemoveUnnecessaryTokens(c string) string {
	c = strings.ReplaceAll(c, "\t", "")
	c = strings.ReplaceAll(c, "\u0000", "")
	c = strings.ReplaceAll(c, "\n", "")
	c = strings.ReplaceAll(c, "\r", "")
	c = strings.Trim(c, " ")
	c = strings.Trim(c, "　")
	return c
}

func IsValidImageUrl(url string) bool {
	isJpg := strings.HasSuffix(url, ".jpg")
	isJpeg := strings.HasSuffix(url, "jpeg")
	isPng := strings.HasSuffix(url, "png")
	isGif := strings.HasSuffix(url, "gif")
	isWebp := strings.HasSuffix(url, "webp")
	return isJpg || isJpeg || isPng || isGif || isWebp
}

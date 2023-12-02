package lib

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Saver struct {
	ID        int
	imageUrls []string
}

func NewSaver(id int) *Saver {
	return &Saver{ID: id, imageUrls: make([]string, 0)}
}

func (s *Saver) AddImageUrl(imageUrl ...string) {
	s.imageUrls = append(s.imageUrls, imageUrl...)
}

func (s *Saver) GetImagePath(postId, imageUrl string) string {
	imageName := postId + "__" + imageUrl[strings.LastIndex(imageUrl, "/")+1:]
	return filepath.Join(s.GetImageDir(), imageName)
}

func (s *Saver) SavePost(postId, year, month, day, title, contnet string) error {
	fileName := fmt.Sprintf("%s-%s-%s %s.md", year, month, day, fmt.Sprintf("%s (%s)", title, postId))
	return os.WriteFile(filepath.Join(s.GetPostDir(), fileName), []byte(contnet), 0755)
}

func (s *Saver) SaveImage(postId, url string) error {
	log.Println("fetch", url)
fetch:
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Fail to get image,", err)
		return err
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		log.Println("Request is malformed")
		return nil
	}
	if resp.StatusCode >= 500 {
		log.Println("Server error, retry")
		goto fetch
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Fail to read body,", err)
		return err
	}
	resp.Body.Close()
	if _, _, err = image.Decode(bytes.NewReader(body)); err != nil {
		log.Println("Image is broken, retry")
		goto fetch
	}
	filename := filepath.Join(s.GetImagePath(postId, url))
	file, err := os.Create(filename)
	if err != nil {
		log.Println("Fail to create file,", filename)
		return err
	}
	defer file.Close()
	if _, err := file.Write(body); err != nil {
		log.Println("Fail to write image,", err)
		return err
	}
	log.Println("save image to", filename)
	return nil
}

func (s *Saver) GetPostDir() string {
	return filepath.Join("skz_blog", strconv.Itoa(s.ID), "posts")
}

func (s *Saver) GetImageDir() string {
	return filepath.Join("skz_blog", strconv.Itoa(s.ID), "images")
}

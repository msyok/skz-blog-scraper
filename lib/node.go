package lib

import "golang.org/x/net/html"

func GetLink(node *html.Node) (name, link string) {
	for i := 0; i < len(node.Attr); i++ {
		attr := node.Attr[i]
		if attr.Key == "href" {
			link = attr.Val
			break
		}
	}
	if node.FirstChild != nil {
		name = node.FirstChild.Data
	}
	return
}

func GetImage(node *html.Node) (name, url string) {
	for i := 0; i < len(node.Attr); i++ {
		attr := node.Attr[i]
		if attr.Key == "src" {
			url = attr.Val
		}
		if attr.Key == "alt" {
			name = attr.Val
		}
	}
	return
}

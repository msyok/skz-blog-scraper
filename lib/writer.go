package lib

import (
	"bytes"
	"fmt"
)

type writer struct {
	buf *bytes.Buffer
}

func NewWriter() *writer {
	return &writer{buf: new(bytes.Buffer)}
}

func (w *writer) WriteText(text string) {
	w.buf.WriteString(text)
}

func (w *writer) WriterLineBreak() {
	w.buf.WriteString("<br>")
}

func (w *writer) WriteLink(name, link string) {
	w.buf.WriteString(fmt.Sprintf("[%s](%s)", name, link))
}

func (w *writer) WriteImage(text, imagePath string) {
	w.buf.WriteString(fmt.Sprintf("![%s](%s)", text, imagePath))
}

func (w *writer) String() string {
	return w.buf.String()
}

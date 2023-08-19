package utils

import (
	"syscall"

	"golang.org/x/term"
)

type Window struct {
	Width  int
	Height int
}

func NewWindow() *Window {
	w := &Window{}
	w.UpdateSize()
	return w
}

func (w *Window) UpdateSize() {
	fd := syscall.Stdout
	width, height, _ := term.GetSize(int(fd))
	w.Width = width
	w.Height = height
}

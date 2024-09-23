package model

import (
	"fmt"
	"strings"
)

type Buffer struct {
	ID         int
	Content    []rune
	FilePath   string
	IsModified bool
}

func NewBuffer(id int, content string, filePath string) *Buffer {
	return &Buffer{
		ID:         id,
		Content:    []rune(content),
		FilePath:   filePath,
		IsModified: false,
	}
}

func (b *Buffer) GetLine(lineNum int) (string, error) {
	lines := strings.Split(string(b.Content), "\n")
	if lineNum < 0 || lineNum >= len(lines) {
		return "", fmt.Errorf("line number out of range")
	}
	return lines[lineNum], nil
}

func (b *Buffer) Save() {
	b.IsModified = false
}

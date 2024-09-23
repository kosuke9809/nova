package model

type Window struct {
	ID       int
	BufferID int
	Cursor   Cursor
	Viewport Viewport
}

type Cursor struct {
	Line   int
	Column int
}

type Viewport struct {
	TopLine    int
	LeftColumn int
	Height     int
	Width      int
}

func NewWindow(id int, bufferID int) *Window {
	return &Window{
		ID:       id,
		BufferID: bufferID,
		Cursor:   Cursor{Line: 0, Column: 0},
		Viewport: Viewport{TopLine: 0, LeftColumn: 0, Height: 24, Width: 80},
	}
}

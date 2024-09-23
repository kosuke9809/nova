package model

type Editor struct {
	Buffers     []*Buffer
	Windows     []*Window
	Tabs        []*Tab
	CurrentMode Mode
}

type Mode int

const (
	NormalMode Mode = iota
)

func NewEditor() *Editor {
	return &Editor{
		Buffers:     make([]*Buffer, 0),
		Windows:     make([]*Window, 0),
		Tabs:        make([]*Tab, 0),
		CurrentMode: NormalMode,
	}
}

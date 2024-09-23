package model

type Editor struct {
	Buffers     []*Buffer
	Windows     []*Window
	Tabs        []*Tab
	CurrentMode Mode
	Settings    EditorSettings
}

type EditorSettings struct {
	TabSize            int
	InsertSpacesForTab bool
	LineNumbers        bool
	SyntaxHighlighting bool
	Theme              string
}

func DefaultEditorSettings() EditorSettings {
	return EditorSettings{
		TabSize:            4,
		InsertSpacesForTab: true,
		LineNumbers:        true,
		SyntaxHighlighting: true,
		Theme:              "default",
	}
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
		Settings:    DefaultEditorSettings(),
	}
}

package memory

import (
	"encoding/json"
	"nova/domain/model"
	"nova/domain/repository"
	"os"
	"sync"
)

type editorRepository struct {
	editor     *model.Editor
	tabRepo    repository.ITabRepository
	widowRepo  repository.IWindowRepository
	bufferRepo repository.IBufferRepository
	filePath   string
	mu         sync.RWMutex
}

func NewEditorRepository(filePath string, tabRepo repository.ITabRepository, windowRepo repository.IWindowRepository, bufferRepo repository.IBufferRepository) repository.IEditorRepository {
	return &editorRepository{
		editor:     model.NewEditor(),
		tabRepo:    tabRepo,
		widowRepo:  windowRepo,
		bufferRepo: bufferRepo,
		filePath:   filePath,
	}
}

type editorState struct {
	TabIDs     []int
	WindowIDs  []int
	BufferIDs  []int
	CurretMode model.Mode
	Settings   model.EditorSettings
}

func (er *editorRepository) Get() (*model.Editor, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	return er.editor, nil
}

func (er *editorRepository) Update(editor *model.Editor) error {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.editor = editor
	return nil
}

func (er *editorRepository) Save(editor *model.Editor) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	state := editorState{
		TabIDs:     make([]int, len(editor.Tabs)),
		WindowIDs:  make([]int, len(editor.Windows)),
		BufferIDs:  make([]int, len(editor.Buffers)),
		CurretMode: editor.CurrentMode,
		Settings:   editor.Settings,
	}

	for i, tab := range editor.Tabs {
		state.TabIDs[i] = tab.ID
	}

	for i, window := range editor.Windows {
		state.WindowIDs[i] = window.ID
	}

	for i, buffer := range editor.Buffers {
		state.BufferIDs[i] = buffer.ID
	}

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(er.filePath, data, 0644)
}

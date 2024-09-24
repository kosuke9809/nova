package memory

import (
	"encoding/json"
	"fmt"
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
	repo := &editorRepository{
		tabRepo:    tabRepo,
		widowRepo:  windowRepo,
		bufferRepo: bufferRepo,
		filePath:   filePath,
	}
	editor, err := repo.Get()
	if err != nil {
		fmt.Printf("Failed to get editor: %v\n", err)
		editor = model.NewEditor()
	}
	repo.editor = editor
	return repo
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

	data, err := os.ReadFile(er.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			newEditor := model.NewEditor()
			er.editor = newEditor
			return newEditor, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		newEditor := model.NewEditor()
		er.editor = newEditor
		return newEditor, nil
	}

	var state editorState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	er.editor = &model.Editor{
		CurrentMode: state.CurretMode,
		Settings:    state.Settings,
		Tabs:        make([]*model.Tab, 0),
		Windows:     make([]*model.Window, 0),
		Buffers:     make([]*model.Buffer, 0),
	}

	for _, tabID := range state.TabIDs {
		tab, err := er.tabRepo.FindByID(tabID)
		if err != nil {
			return nil, err
		}
		er.editor.Tabs = append(er.editor.Tabs, tab)
	}

	for _, windowID := range state.WindowIDs {
		window, err := er.widowRepo.FindByID(windowID)
		if err != nil {
			return nil, err
		}
		er.editor.Windows = append(er.editor.Windows, window)
	}

	for _, bufferID := range state.BufferIDs {
		buffer, err := er.bufferRepo.FindByID(bufferID)
		if err != nil {
			return nil, err
		}
		er.editor.Buffers = append(er.editor.Buffers, buffer)
	}

	return er.editor, nil
}

func (er *editorRepository) Save(editor *model.Editor) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	if er.editor == nil {
		er.editor = editor
	}

	return er.saveToFile(editor)
}

func (er *editorRepository) Update(editor *model.Editor) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	// ファイルが存在するか確認
	if _, err := os.Stat(er.filePath); os.IsNotExist(err) {
		return fmt.Errorf("editor does not exist, cannot update")
	}

	// ファイルが空でないことを確認
	info, err := os.Stat(er.filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("editor does not exist, cannot update")
	}

	return er.saveToFile(editor)
}

func (er *editorRepository) saveToFile(editor *model.Editor) error {
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

package memory_test

import (
	"encoding/json"
	"nova/domain/model"
	"nova/infrastructure/persistence/memory"
	mocks "nova/mocks"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEditorRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTabRepo := mocks.NewMockITabRepository(ctrl)
	mockWindowRepo := mocks.NewMockIWindowRepository(ctrl)
	mockBufferRepo := mocks.NewMockIBufferRepository(ctrl)

	tempFile, err := os.CreateTemp("", "editor_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	t.Run("NewEditorRepository", func(t *testing.T) {
		repo := memory.NewEditorRepository(tempFile.Name(), mockTabRepo, mockWindowRepo, mockBufferRepo)
		assert.NotNil(t, repo, "NewEditorRepository should return a non-nil repository")
	})

	t.Run("Get", func(t *testing.T) {
		mockTabRepo.EXPECT().FindByID(gomock.Any()).Return(&model.Tab{}, nil).AnyTimes()
		mockWindowRepo.EXPECT().FindByID(gomock.Any()).Return(&model.Window{}, nil).AnyTimes()
		mockBufferRepo.EXPECT().FindByID(gomock.Any()).Return(&model.Buffer{}, nil).AnyTimes()

		repo := memory.NewEditorRepository(tempFile.Name(), mockTabRepo, mockWindowRepo, mockBufferRepo)
		editor, err := repo.Get()
		if err != nil {
			t.Fatalf("Failed to get editor: %v", err)
		}
		assert.NotNil(t, editor, "Get should return a non-nil editor")
		if editor != nil {
			assert.Equal(t, model.NormalMode, editor.CurrentMode, "Default CurrentMode should be NormalMode")
			assert.Equal(t, model.DefaultEditorSettings(), editor.Settings, "Default Settings should match DefaultEditorSettings")
		}
	})

	t.Run("Save", func(t *testing.T) {
		repo := memory.NewEditorRepository(tempFile.Name(), mockTabRepo, mockWindowRepo, mockBufferRepo)
		editor := model.NewEditor()
		editor.Settings.Theme = "dark"
		editor.Settings.TabSize = 2

		tab := &model.Tab{ID: 1}
		window := &model.Window{ID: 1}
		buffer := &model.Buffer{ID: 1}

		editor.Tabs = append(editor.Tabs, tab)
		editor.Windows = append(editor.Windows, window)
		editor.Buffers = append(editor.Buffers, buffer)

		err := repo.Save(editor)
		assert.NoError(t, err)

		// ファイルから保存されたデータを読み込んで検証
		data, err := os.ReadFile(tempFile.Name())
		assert.NoError(t, err)

		var state struct {
			TabIDs     []int
			WindowIDs  []int
			BufferIDs  []int
			CurretMode model.Mode
			Settings   model.EditorSettings
		}

		err = json.Unmarshal(data, &state)
		assert.NoError(t, err)

		assert.Equal(t, []int{1}, state.TabIDs)
		assert.Equal(t, []int{1}, state.WindowIDs)
		assert.Equal(t, []int{1}, state.BufferIDs)
		assert.Equal(t, model.NormalMode, state.CurretMode)
		assert.Equal(t, "dark", state.Settings.Theme)
		assert.Equal(t, 2, state.Settings.TabSize)
	})

	t.Run("Update", func(t *testing.T) {
		repo := memory.NewEditorRepository(tempFile.Name(), mockTabRepo, mockWindowRepo, mockBufferRepo)

		// 初期エディタを保存
		initialEditor := model.NewEditor()
		err := repo.Save(initialEditor)
		assert.NoError(t, err)

		// エディタを更新
		updatedEditor := model.NewEditor()
		updatedEditor.Settings.Theme = "light"
		updatedEditor.Settings.TabSize = 4

		err = repo.Update(updatedEditor)
		assert.NoError(t, err)

		// 更新されたデータを検証
		loadedEditor, err := repo.Get()
		assert.NoError(t, err)
		assert.NotNil(t, loadedEditor)
		assert.Equal(t, "light", loadedEditor.Settings.Theme)
		assert.Equal(t, 4, loadedEditor.Settings.TabSize)
	})

	t.Run("Update non-existent editor", func(t *testing.T) {
		// 新しいリポジトリを作成（エディタが存在しない状態）
		newTempFile, _ := os.CreateTemp("", "new_editor_test")
		defer os.Remove(newTempFile.Name())

		newRepo := memory.NewEditorRepository(newTempFile.Name(), mockTabRepo, mockWindowRepo, mockBufferRepo)

		// 存在しないエディタの更新を試みる
		err := newRepo.Update(model.NewEditor())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "editor does not exist, cannot update")
	})

	t.Run("SaveAndGet", func(t *testing.T) {
		// モックの振る舞いを定義
		mockTabRepo.EXPECT().FindByID(gomock.Any()).Return(&model.Tab{}, nil).AnyTimes()
		mockWindowRepo.EXPECT().FindByID(gomock.Any()).Return(&model.Window{}, nil).AnyTimes()
		mockBufferRepo.EXPECT().FindByID(gomock.Any()).Return(&model.Buffer{}, nil).AnyTimes()

		repo := memory.NewEditorRepository(tempFile.Name(), mockTabRepo, mockWindowRepo, mockBufferRepo)
		editor := model.NewEditor()
		editor.Settings.Theme = "dark"
		editor.Settings.TabSize = 2

		err := repo.Save(editor)
		assert.NoError(t, err)

		loadedEditor, err := repo.Get()
		assert.NoError(t, err)
		assert.NotNil(t, loadedEditor, "Get should return a non-nil editor after save")
		assert.Equal(t, "dark", loadedEditor.Settings.Theme)
		assert.Equal(t, 2, loadedEditor.Settings.TabSize)
	})
}

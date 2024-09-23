package repository

import "nova/domain/model"

type IEditorRepository interface {
	// Load() (*model.Editor, error)
	Get() (*model.Editor, error)
	Save(editor *model.Editor) error
	Update(editor *model.Editor) error
}

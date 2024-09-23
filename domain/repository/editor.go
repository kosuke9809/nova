package repository

import "nova/domain/model"

type EditorRepository interface {
	Save(editor *model.Editor) error
	Load() (*model.Editor, error)
}

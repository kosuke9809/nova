package repository

import "nova/domain/model"

type IWindowRepository interface {
	Save(window *model.Window) error
	Update(window *model.Window) error
	Delete(id int) error
	FindByID(id int) (*model.Window, error)
	List() ([]*model.Window, error)
	FindByBufferID(bufferID int) ([]*model.Window, error)
}

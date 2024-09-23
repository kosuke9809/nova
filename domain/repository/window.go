package repository

import "nova/domain/model"

type WindowRepository interface {
	Save(window *model.Window) error
	Update(window *model.Window) error
	Delete(id int) error
	FindByID(id int) (*model.Window, error)
	ListByTabID(tabID int) ([]*model.Window, error)
}

package repository

import "nova/domain/model"

type TabRepository interface {
	Save(tab *model.Tab) error
	Update(tab *model.Tab) error
	Delete(id int) error
	FindByID(id int) (*model.Tab, error)
	List() ([]*model.Tab, error)
}

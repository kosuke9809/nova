package repository

import "nova/domain/model"

type BufferRepository interface {
	Save(buffer *model.Buffer) error
	Update(buffer *model.Buffer) error
	Delete(id int) error
	FindByID(id int) (*model.Buffer, error)
	List() ([]*model.Buffer, error)
}

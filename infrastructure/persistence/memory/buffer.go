package memory

import (
	"fmt"
	"nova/domain/model"
	"nova/domain/repository"
	"sync"
)

type bufferRepository struct {
	buffers map[int]*model.Buffer
	nextID  int
	mu      sync.RWMutex
}

func NewBufferRepository() repository.IBufferRepository {
	return &bufferRepository{
		buffers: make(map[int]*model.Buffer),
		nextID:  1,
	}
}

func (br *bufferRepository) FindByID(id int) (*model.Buffer, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	buffer, exists := br.buffers[id]
	if !exists {
		return nil, fmt.Errorf("buffer with id %d not found", id)
	}
	return buffer, nil
}

func (br *bufferRepository) Save(buffer *model.Buffer) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if buffer.ID == 0 {
		buffer.ID = br.nextID
		br.nextID++
	}
	br.buffers[buffer.ID] = buffer
	return nil
}

func (br *bufferRepository) Update(buffer *model.Buffer) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if _, exists := br.buffers[buffer.ID]; !exists {
		return fmt.Errorf("buffer with id %d not found", buffer.ID)
	}
	br.buffers[buffer.ID] = buffer
	return nil
}

func (br *bufferRepository) Delete(id int) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if _, exists := br.buffers[id]; !exists {
		return fmt.Errorf("buffer with id %d not found", id)
	}
	delete(br.buffers, id)
	return nil
}

func (br *bufferRepository) List() ([]*model.Buffer, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	buffers := make([]*model.Buffer, 0, len(br.buffers))
	for _, buffer := range br.buffers {
		buffers = append(buffers, buffer)
	}
	return buffers, nil
}

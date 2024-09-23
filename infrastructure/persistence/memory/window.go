package memory

import (
	"fmt"
	"nova/domain/model"
	"nova/domain/repository"
	"sync"
)

type windowRepository struct {
	windows map[int]*model.Window
	nextID  int
	mu      sync.RWMutex
}

func NewWindowRepository() repository.IWindowRepository {
	return &windowRepository{
		windows: make(map[int]*model.Window),
		nextID:  1,
	}
}

func (wr *windowRepository) FindByID(id int) (*model.Window, error) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	window, exists := wr.windows[id]
	if !exists {
		return nil, fmt.Errorf("window with id %d not found", id)
	}
	return window, nil
}

func (wr *windowRepository) Save(window *model.Window) error {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	if window.ID == 0 {
		window.ID = wr.nextID
		wr.nextID++
	}
	wr.windows[window.ID] = window
	return nil
}

func (wr *windowRepository) Update(window *model.Window) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if _, exists := wr.windows[window.ID]; !exists {
		return fmt.Errorf("window with id %d not found", window.ID)
	}
	wr.windows[window.ID] = window
	return nil
}

func (wr *windowRepository) Delete(id int) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if _, exists := wr.windows[id]; !exists {
		return fmt.Errorf("window with id %d not found", id)
	}
	delete(wr.windows, id)
	return nil
}

func (wr *windowRepository) List() ([]*model.Window, error) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	windows := make([]*model.Window, 0, len(wr.windows))
	for _, window := range wr.windows {
		windows = append(windows, window)
	}
	return windows, nil
}

func (wr *windowRepository) FindByBufferID(bufferID int) ([]*model.Window, error) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	var result []*model.Window
	for _, window := range wr.windows {
		if window.BufferID == bufferID {
			result = append(result, window)
		}
	}
	return result, nil
}

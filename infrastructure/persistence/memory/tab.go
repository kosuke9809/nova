package memory

import (
	"fmt"
	"nova/domain/model"
	"nova/domain/repository"
	"sync"
)

type tabRepository struct {
	tabs   map[int]*model.Tab
	nextID int
	mu     sync.RWMutex
}

func NewTabRepository() repository.ITabRepository {
	return &tabRepository{
		tabs:   make(map[int]*model.Tab),
		nextID: 1,
	}
}

func (tr *tabRepository) FindByID(id int) (*model.Tab, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tab, exists := tr.tabs[id]
	if !exists {
		return nil, fmt.Errorf("tab with id %d not found", id)
	}
	return tab, nil
}

func (tr *tabRepository) Save(tab *model.Tab) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tab.ID == 0 {
		tab.ID = tr.nextID
		tr.nextID++
	}
	tr.tabs[tab.ID] = tab
	return nil
}

func (tr *tabRepository) Update(tab *model.Tab) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if _, exists := tr.tabs[tab.ID]; !exists {
		return fmt.Errorf("tab with id %d not found", tab.ID)
	}
	tr.tabs[tab.ID] = tab
	return nil
}

func (tr *tabRepository) Delete(id int) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if _, exists := tr.tabs[id]; !exists {
		return fmt.Errorf("tab with id %d not found", id)
	}
	delete(tr.tabs, id)
	return nil
}

func (tr *tabRepository) List() ([]*model.Tab, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tabs := make([]*model.Tab, 0, len(tr.tabs))
	for _, tab := range tr.tabs {
		tabs = append(tabs, tab)
	}
	return tabs, nil
}

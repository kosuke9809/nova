package memory_test

import (
	"nova/domain/model"
	"nova/infrastructure/persistence/memory"
	"sync"
	"testing"
)

func TestTabRepository(t *testing.T) {
	t.Run("NewTabRepository", func(t *testing.T) {
		repo := memory.NewTabRepository()
		if repo == nil {
			t.Error("NewTabRepository should return a non-nil repository")
		}
	})

	t.Run("SaveAndFindByID", func(t *testing.T) {
		repo := memory.NewTabRepository()
		tab := model.NewTab(0) // ID will be assigned by the repository

		err := repo.Save(tab)
		if err != nil {
			t.Errorf("Failed to save tab: %v", err)
		}

		if tab.ID == 0 {
			t.Error("Save should assign an ID to the tab")
		}

		foundTab, err := repo.FindByID(tab.ID)
		if err != nil {
			t.Errorf("Failed to find tab by ID: %v", err)
		}

		if foundTab.ID != tab.ID {
			t.Errorf("Expected tab ID %d, got %d", tab.ID, foundTab.ID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		repo := memory.NewTabRepository()
		tab := model.NewTab(0)
		window := &model.Window{ID: 1}
		tab.AddWindow(window)

		err := repo.Save(tab)
		if err != nil {
			t.Errorf("Failed to save tab: %v", err)
		}

		newWindow := &model.Window{ID: 2}
		tab.AddWindow(newWindow)
		tab.SetActiveWindow(2)

		err = repo.Update(tab)
		if err != nil {
			t.Errorf("Failed to update tab: %v", err)
		}

		updatedTab, err := repo.FindByID(tab.ID)
		if err != nil {
			t.Errorf("Failed to find updated tab: %v", err)
		}

		if len(updatedTab.Windows) != 2 {
			t.Errorf("Expected 2 windows, got %d", len(updatedTab.Windows))
		}

		if updatedTab.ActiveWindow.ID != 2 {
			t.Errorf("Expected active window ID 2, got %d", updatedTab.ActiveWindow.ID)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		repo := memory.NewTabRepository()
		tab := model.NewTab(0)

		err := repo.Save(tab)
		if err != nil {
			t.Errorf("Failed to save tab: %v", err)
		}

		err = repo.Delete(tab.ID)
		if err != nil {
			t.Errorf("Failed to delete tab: %v", err)
		}

		_, err = repo.FindByID(tab.ID)
		if err == nil {
			t.Error("Expected error when finding deleted tab, got nil")
		}
	})

	t.Run("List", func(t *testing.T) {
		repo := memory.NewTabRepository()
		numTabs := 3

		for i := 0; i < numTabs; i++ {
			tab := model.NewTab(0)
			err := repo.Save(tab)
			if err != nil {
				t.Errorf("Failed to save tab: %v", err)
			}
		}

		listedTabs, err := repo.List()
		if err != nil {
			t.Errorf("Failed to list tabs: %v", err)
		}

		if len(listedTabs) != numTabs {
			t.Errorf("Expected %d tabs, got %d", numTabs, len(listedTabs))
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		repo := memory.NewTabRepository()
		numGoroutines := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				tab := model.NewTab(0)
				window := &model.Window{ID: i}
				tab.AddWindow(window)
				err := repo.Save(tab)
				if err != nil {
					t.Errorf("Failed to save tab in goroutine: %v", err)
				}
			}(i)
		}

		wg.Wait()

		tabs, err := repo.List()
		if err != nil {
			t.Errorf("Failed to list tabs after concurrent saves: %v", err)
		}

		if len(tabs) != numGoroutines {
			t.Errorf("Expected %d tabs after concurrent saves, got %d", numGoroutines, len(tabs))
		}
	})

	t.Run("Window Operations", func(t *testing.T) {
		repo := memory.NewTabRepository()
		tab := model.NewTab(0)

		// Add windows
		tab.AddWindow(&model.Window{ID: 1})
		tab.AddWindow(&model.Window{ID: 2})
		tab.AddWindow(&model.Window{ID: 3})

		err := repo.Save(tab)
		if err != nil {
			t.Errorf("Failed to save tab: %v", err)
		}

		// Check active window
		if tab.ActiveWindow.ID != 1 {
			t.Errorf("Expected active window ID 1, got %d", tab.ActiveWindow.ID)
		}

		// Set active window
		tab.SetActiveWindow(2)
		err = repo.Update(tab)
		if err != nil {
			t.Errorf("Failed to update tab: %v", err)
		}

		updatedTab, err := repo.FindByID(tab.ID)
		if err != nil {
			t.Errorf("Failed to find updated tab: %v", err)
		}

		if updatedTab.ActiveWindow.ID != 2 {
			t.Errorf("Expected active window ID 2, got %d", updatedTab.ActiveWindow.ID)
		}

		// Remove window
		tab.RemoveWindow(2)
		err = repo.Update(tab)
		if err != nil {
			t.Errorf("Failed to update tab after removing window: %v", err)
		}

		finalTab, err := repo.FindByID(tab.ID)
		if err != nil {
			t.Errorf("Failed to find final tab: %v", err)
		}

		if len(finalTab.Windows) != 2 {
			t.Errorf("Expected 2 windows after removal, got %d", len(finalTab.Windows))
		}

		if finalTab.ActiveWindow.ID != 1 {
			t.Errorf("Expected active window to revert to ID 1, got %d", finalTab.ActiveWindow.ID)
		}
	})
}

package memory_test

import (
	"nova/domain/model"
	"nova/infrastructure/persistence/memory"
	"sync"
	"testing"
)

func TestWindowRepository(t *testing.T) {
	t.Run("NewWindowRepository", func(t *testing.T) {
		repo := memory.NewWindowRepository()
		if repo == nil {
			t.Error("NewWindowRepository should return a non-nil repository")
		}
	})

	t.Run("SaveAndFindByID", func(t *testing.T) {
		repo := memory.NewWindowRepository()
		window := &model.Window{BufferID: 1}

		err := repo.Save(window)
		if err != nil {
			t.Errorf("Failed to save window: %v", err)
		}

		if window.ID == 0 {
			t.Error("Save should assign an ID to the window")
		}

		foundWindow, err := repo.FindByID(window.ID)
		if err != nil {
			t.Errorf("Failed to find window by ID: %v", err)
		}

		if foundWindow.ID != window.ID || foundWindow.BufferID != window.BufferID {
			t.Errorf("Found window does not match saved window")
		}
	})

	t.Run("Update", func(t *testing.T) {
		repo := memory.NewWindowRepository()
		window := &model.Window{BufferID: 1}

		err := repo.Save(window)
		if err != nil {
			t.Errorf("Failed to save window: %v", err)
		}

		window.BufferID = 2
		err = repo.Update(window)
		if err != nil {
			t.Errorf("Failed to update window: %v", err)
		}

		updatedWindow, err := repo.FindByID(window.ID)
		if err != nil {
			t.Errorf("Failed to find updated window: %v", err)
		}

		if updatedWindow.BufferID != 2 {
			t.Errorf("Expected updated BufferID 2, got %d", updatedWindow.BufferID)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		repo := memory.NewWindowRepository()
		window := &model.Window{BufferID: 1}

		err := repo.Save(window)
		if err != nil {
			t.Errorf("Failed to save window: %v", err)
		}

		err = repo.Delete(window.ID)
		if err != nil {
			t.Errorf("Failed to delete window: %v", err)
		}

		_, err = repo.FindByID(window.ID)
		if err == nil {
			t.Error("Expected error when finding deleted window, got nil")
		}
	})

	t.Run("List", func(t *testing.T) {
		repo := memory.NewWindowRepository()
		numWindows := 3

		for i := 0; i < numWindows; i++ {
			window := &model.Window{BufferID: i + 1}
			err := repo.Save(window)
			if err != nil {
				t.Errorf("Failed to save window: %v", err)
			}
		}

		windows, err := repo.List()
		if err != nil {
			t.Errorf("Failed to list windows: %v", err)
		}

		if len(windows) != numWindows {
			t.Errorf("Expected %d windows, got %d", numWindows, len(windows))
		}
	})

	t.Run("FindByBufferID", func(t *testing.T) {
		repo := memory.NewWindowRepository()

		// Create windows with different BufferIDs
		window1 := &model.Window{BufferID: 1}
		window2 := &model.Window{BufferID: 2}
		window3 := &model.Window{BufferID: 1}

		repo.Save(window1)
		repo.Save(window2)
		repo.Save(window3)

		windows, err := repo.FindByBufferID(1)
		if err != nil {
			t.Errorf("Failed to find windows by BufferID: %v", err)
		}

		if len(windows) != 2 {
			t.Errorf("Expected 2 windows with BufferID 1, got %d", len(windows))
		}

		for _, w := range windows {
			if w.BufferID != 1 {
				t.Errorf("Found window with unexpected BufferID: %d", w.BufferID)
			}
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		repo := memory.NewWindowRepository()
		numGoroutines := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				window := &model.Window{BufferID: i}
				err := repo.Save(window)
				if err != nil {
					t.Errorf("Failed to save window in goroutine: %v", err)
				}
			}(i)
		}

		wg.Wait()

		windows, err := repo.List()
		if err != nil {
			t.Errorf("Failed to list windows after concurrent saves: %v", err)
		}

		if len(windows) != numGoroutines {
			t.Errorf("Expected %d windows after concurrent saves, got %d", numGoroutines, len(windows))
		}
	})
}

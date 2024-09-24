package memory_test

import (
	"fmt"
	"nova/domain/model"
	"nova/infrastructure/persistence/memory"
	"sync"
	"testing"
)

func TestBufferRepository(t *testing.T) {
	t.Run("NewBufferRepository", func(t *testing.T) {
		repo := memory.NewBufferRepository()
		if repo == nil {
			t.Error("NewBufferRepository should return a non-nil repository")
		}
	})

	t.Run("Save and FindByID", func(t *testing.T) {
		repo := memory.NewBufferRepository()
		buffer := &model.Buffer{Content: []rune("Test content")}

		err := repo.Save(buffer)
		if err != nil {
			t.Errorf("Failed to save buffer: %v", err)
		}

		if buffer.ID == 0 {
			t.Error("Save should assign an ID to the buffer")
		}

		foundBuffer, err := repo.FindByID(buffer.ID)
		if err != nil {
			t.Errorf("Failed to find buffer by ID: %v", err)
		}

		if string(foundBuffer.Content) != string(buffer.Content) {
			t.Errorf("Found buffer content does not match. Expected %s, got %s", string(buffer.Content), string(foundBuffer.Content))
		}
	})

	t.Run("Update", func(t *testing.T) {
		repo := memory.NewBufferRepository()
		buffer := &model.Buffer{Content: []rune("Original content")}

		err := repo.Save(buffer)
		if err != nil {
			t.Errorf("Failed to save buffer: %v", err)
		}

		buffer.Content = []rune("Updated content")
		err = repo.Update(buffer)
		if err != nil {
			t.Errorf("Failed to update buffer: %v", err)
		}

		updatedBuffer, err := repo.FindByID(buffer.ID)
		if err != nil {
			t.Errorf("Failed to find updated buffer: %v", err)
		}

		if string(updatedBuffer.Content) != "Updated content" {
			t.Errorf("Buffer content was not updated. Expected 'Updated content', got '%s'", string(updatedBuffer.Content))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		repo := memory.NewBufferRepository()
		buffer := &model.Buffer{Content: []rune("To be deleted")}

		err := repo.Save(buffer)
		if err != nil {
			t.Errorf("Failed to save buffer: %v", err)
		}

		err = repo.Delete(buffer.ID)
		if err != nil {
			t.Errorf("Failed to delete buffer: %v", err)
		}

		_, err = repo.FindByID(buffer.ID)
		if err == nil {
			t.Error("Expected error when finding deleted buffer, got nil")
		}
	})

	t.Run("List", func(t *testing.T) {
		repo := memory.NewBufferRepository()
		buffer1 := &model.Buffer{Content: []rune("Buffer1")}
		buffer2 := &model.Buffer{Content: []rune("Buffer 2")}

		repo.Save(buffer1)
		repo.Save(buffer2)

		buffers, err := repo.List()
		if err != nil {
			t.Errorf("Failed to list buffers: %v", err)
		}

		if len(buffers) != 2 {
			t.Errorf("Expected 2 buffers, got %d", len(buffers))
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		repo := memory.NewBufferRepository()
		numGoroutines := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				buffer := &model.Buffer{Content: []rune(fmt.Sprintf("Buffer %d", i))}
				err := repo.Save(buffer)
				if err != nil {
					t.Errorf("Failed to save buffer in goroutine: %v", err)
				}
			}(i)
		}

		wg.Wait()

		buffers, err := repo.List()
		if err != nil {
			t.Errorf("Failed to list buffers after concurrent saves: %v", err)
		}

		if len(buffers) != numGoroutines {
			t.Errorf("Expected %d buffers after concurrent saves, got %d", numGoroutines, len(buffers))
		}
	})
}

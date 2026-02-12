package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store manages task persistence
type Store struct {
	dir string
	mu  sync.RWMutex
}

// NewStore creates a new task store
func NewStore(configDir string) (*Store, error) {
	tasksDir := filepath.Join(configDir, "tasks")
	
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return nil, fmt.Errorf("create tasks dir: %w", err)
	}
	
	return &Store{dir: tasksDir}, nil
}

// Save persists a task
func (s *Store) Save(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data, err := task.ToJSON()
	if err != nil {
		return err
	}
	
	path := s.taskPath(task.ID)
	
	// Atomic write
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}
	
	return os.Rename(tmpPath, path)
}

// Get retrieves a task by ID
func (s *Store) Get(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	path := s.taskPath(id)
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	return FromJSON(data)
}

// List returns all tasks
func (s *Store) List() ([]*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	
	var tasks []*Task
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		
		path := filepath.Join(s.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		
		task, err := FromJSON(data)
		if err != nil {
			continue
		}
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// Delete removes a task
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	path := s.taskPath(id)
	return os.Remove(path)
}

// ListActive returns non-terminal tasks
func (s *Store) ListActive() ([]*Task, error) {
	all, err := s.List()
	if err != nil {
		return nil, err
	}
	
	var active []*Task
	for _, t := range all {
		if !t.IsTerminal() {
			active = append(active, t)
		}
	}
	
	return active, nil
}

// taskPath returns the file path for a task
func (s *Store) taskPath(id string) string {
	return filepath.Join(s.dir, id+".json")
}

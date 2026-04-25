package memory

import (
	"context"
	"fmt"
	"time"

	"carrotagent/carrot-agent/pkg/storage"
)

const (
	MemoryTypeSnapshot = "snapshot"
	MemoryTypeSession  = "session"
	MemoryTypeLongTerm = "longterm"
)

type MemoryManager struct {
	store    *storage.Store
	memories map[string]*storage.Memory
}

func NewMemoryManager(store *storage.Store) *MemoryManager {
	return &MemoryManager{
		store:    store,
		memories: make(map[string]*storage.Memory),
	}
}

func (m *MemoryManager) Load(ctx context.Context) error {
	memoryTypes := []string{MemoryTypeSnapshot, MemoryTypeSession, MemoryTypeLongTerm}

	for _, memType := range memoryTypes {
		memories, err := m.store.ListMemories(memType, 100)
		if err != nil {
			return fmt.Errorf("failed to load memories of type %s: %w", memType, err)
		}

		for _, mem := range memories {
			m.memories[mem.ID] = mem
		}
	}

	return nil
}

func (m *MemoryManager) Add(ctx context.Context, memType, content, metadata string) error {
	id := fmt.Sprintf("%s_%d", memType, time.Now().UnixNano())

	mem := &storage.Memory{
		ID:        id,
		Type:      memType,
		Content:   content,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	if err := m.store.SaveMemory(mem); err != nil {
		return fmt.Errorf("failed to save memory: %w", err)
	}

	m.memories[id] = mem
	return nil
}

func (m *MemoryManager) Get(id string) (*storage.Memory, error) {
	mem, ok := m.memories[id]
	if !ok {
		return m.store.GetMemory(id)
	}
	return mem, nil
}

func (m *MemoryManager) Update(ctx context.Context, id, content string) error {
	mem, err := m.Get(id)
	if err != nil {
		return err
	}
	if mem == nil {
		return fmt.Errorf("memory not found: %s", id)
	}

	mem.Content = content

	if err := m.store.SaveMemory(mem); err != nil {
		return fmt.Errorf("failed to update memory: %w", err)
	}

	m.memories[id] = mem
	return nil
}

func (m *MemoryManager) Delete(ctx context.Context, id string) error {
	if err := m.store.DeleteMemory(id); err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	delete(m.memories, id)
	return nil
}

func (m *MemoryManager) List(memType string, limit int) ([]*storage.Memory, error) {
	if limit <= 0 {
		limit = 50
	}

	return m.store.ListMemories(memType, limit)
}

func (m *MemoryManager) Search(keyword string, limit int) ([]*storage.Memory, error) {
	if limit <= 0 {
		limit = 50
	}

	return m.store.SearchMemories(keyword, limit)
}

func (m *MemoryManager) GetSnapshotContent() string {
	var snapshots []*storage.Memory

	for _, mem := range m.memories {
		if mem.Type == MemoryTypeSnapshot {
			snapshots = append(snapshots, mem)
		}
	}

	if len(snapshots) == 0 {
		return ""
	}

	content := ""
	for _, snap := range snapshots {
		content += snap.Content + "\n\n"
	}

	return content
}

func (m *MemoryManager) GetLongTermMemories(limit int) ([]*storage.Memory, error) {
	if limit <= 0 {
		limit = 50
	}

	return m.store.ListMemories(MemoryTypeLongTerm, limit)
}

func (m *MemoryManager) SaveSnapshot(ctx context.Context, content string) error {
	return m.Add(ctx, MemoryTypeSnapshot, content, "{}")
}

func (m *MemoryManager) SaveLongTermMemory(ctx context.Context, content, metadata string) error {
	return m.Add(ctx, MemoryTypeLongTerm, content, metadata)
}

func (m *MemoryManager) GetMemoryStats() map[string]int {
	stats := make(map[string]int)

	for _, mem := range m.memories {
		stats[mem.Type]++
	}

	return stats
}
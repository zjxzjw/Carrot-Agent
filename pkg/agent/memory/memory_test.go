package memory

import (
	"context"
	"testing"

	"carrotagent/carrot-agent/pkg/storage"
)

func TestNewMemoryManager(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	if mgr == nil {
		t.Fatal("Expected MemoryManager, got nil")
	}

	if mgr.store != store {
		t.Error("Expected store to be set")
	}

	if mgr.memories == nil {
		t.Error("Expected memories map to be initialized")
	}
}

func TestMemoryManagerLoad(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)

	mem := &storage.Memory{
		ID:        "snapshot_1",
		Type:      MemoryTypeSnapshot,
		Content:   "test content",
		Metadata:  "{}",
	}
	if err := store.SaveMemory(mem); err != nil {
		t.Fatalf("Failed to save memory: %v", err)
	}

	err = mgr.Load(context.Background())
	if err != nil {
		t.Fatalf("Failed to load memories: %v", err)
	}

	if len(mgr.memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(mgr.memories))
	}
}

func TestMemoryManagerAdd(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.Add(ctx, MemoryTypeSnapshot, "test content", "{}")
	if err != nil {
		t.Fatalf("Failed to add memory: %v", err)
	}

	if len(mgr.memories) != 1 {
		t.Errorf("Expected 1 memory in cache, got %d", len(mgr.memories))
	}

	memories, err := mgr.List(MemoryTypeSnapshot, 10)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory in storage, got %d", len(memories))
	}
}

func TestMemoryManagerGet(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.Add(ctx, MemoryTypeSnapshot, "test content", "{}")
	if err != nil {
		t.Fatalf("Failed to add memory: %v", err)
	}

	memories, _ := mgr.List(MemoryTypeSnapshot, 10)
	if len(memories) == 0 {
		t.Fatal("Expected at least 1 memory")
	}

	mem, err := mgr.Get(memories[0].ID)
	if err != nil {
		t.Fatalf("Failed to get memory: %v", err)
	}

	if mem == nil {
		t.Fatal("Expected memory, got nil")
	}

	if mem.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", mem.Content)
	}
}

func TestMemoryManagerGetNonExistent(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)

	mem, err := mgr.Get("non_existent")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mem != nil {
		t.Error("Expected nil for non-existent memory")
	}
}

func TestMemoryManagerUpdate(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.Add(ctx, MemoryTypeSnapshot, "original content", "{}")
	if err != nil {
		t.Fatalf("Failed to add memory: %v", err)
	}

	memories, _ := mgr.List(MemoryTypeSnapshot, 10)
	if len(memories) == 0 {
		t.Fatal("Expected at least 1 memory")
	}

	memoryID := memories[0].ID
	err = mgr.Update(ctx, memoryID, "updated content")
	if err != nil {
		t.Fatalf("Failed to update memory: %v", err)
	}

	memory, _ := mgr.Get(memoryID)
	if memory.Content != "updated content" {
		t.Errorf("Expected updated content, got '%s'", memory.Content)
	}
}

func TestMemoryManagerUpdateNonExistent(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.Update(ctx, "non_existent", "new content")
	if err == nil {
		t.Error("Expected error for updating non-existent memory")
	}
}

func TestMemoryManagerDelete(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.Add(ctx, MemoryTypeSnapshot, "test content", "{}")
	if err != nil {
		t.Fatalf("Failed to add memory: %v", err)
	}

	memories, _ := mgr.List(MemoryTypeSnapshot, 10)
	if len(memories) == 0 {
		t.Fatal("Expected at least 1 memory")
	}

	memoryID := memories[0].ID
	err = mgr.Delete(ctx, memoryID)
	if err != nil {
		t.Fatalf("Failed to delete memory: %v", err)
	}

	if len(mgr.memories) != 0 {
		t.Errorf("Expected 0 memories in cache, got %d", len(mgr.memories))
	}
}

func TestMemoryManagerList(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	mgr.Add(ctx, MemoryTypeSnapshot, "snapshot 1", "{}")
	mgr.Add(ctx, MemoryTypeSession, "session 1", "{}")
	mgr.Add(ctx, MemoryTypeLongTerm, "longterm 1", "{}")

	memories, err := mgr.List(MemoryTypeSnapshot, 10)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(memories))
	}
}

func TestMemoryManagerListWithLimit(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		mgr.Add(ctx, MemoryTypeSnapshot, "content", "{}")
	}

	memories, err := mgr.List(MemoryTypeSnapshot, 3)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(memories) != 3 {
		t.Errorf("Expected 3 memories (limit), got %d", len(memories))
	}
}

func TestMemoryManagerListWithZeroLimit(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	mgr.Add(ctx, MemoryTypeSnapshot, "content", "{}")

	memories, err := mgr.List(MemoryTypeSnapshot, 0)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory with zero limit (default), got %d", len(memories))
	}
}

func TestMemoryManagerSearch(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	mgr.Add(ctx, MemoryTypeSnapshot, "golang tutorial", "{}")
	mgr.Add(ctx, MemoryTypeSnapshot, "python tutorial", "{}")
	mgr.Add(ctx, MemoryTypeLongTerm, "golang best practices", "{}")

	results, err := mgr.Search("golang", 10)
	if err != nil {
		t.Fatalf("Failed to search memories: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestMemoryManagerGetSnapshotContent(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	mgr.Add(ctx, MemoryTypeSnapshot, "first snapshot", "{}")
	mgr.Add(ctx, MemoryTypeSnapshot, "second snapshot", "{}")
	mgr.Add(ctx, MemoryTypeSession, "session content", "{}")

	content := mgr.GetSnapshotContent()
	if content == "" {
		t.Fatal("Expected non-empty snapshot content")
	}
}

func TestMemoryManagerGetSnapshotContentEmpty(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)

	content := mgr.GetSnapshotContent()
	if content != "" {
		t.Errorf("Expected empty content, got '%s'", content)
	}
}

func TestMemoryManagerGetLongTermMemories(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	mgr.Add(ctx, MemoryTypeLongTerm, "longterm 1", "{}")
	mgr.Add(ctx, MemoryTypeLongTerm, "longterm 2", "{}")
	mgr.Add(ctx, MemoryTypeSnapshot, "snapshot 1", "{}")

	memories, err := mgr.GetLongTermMemories(10)
	if err != nil {
		t.Fatalf("Failed to get long term memories: %v", err)
	}

	if len(memories) != 2 {
		t.Errorf("Expected 2 longterm memories, got %d", len(memories))
	}
}

func TestMemoryManagerSaveSnapshot(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.SaveSnapshot(ctx, "snapshot content")
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	memories, _ := mgr.List(MemoryTypeSnapshot, 10)
	if len(memories) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(memories))
	}
}

func TestMemoryManagerSaveLongTermMemory(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	err = mgr.SaveLongTermMemory(ctx, "long term content", `{"key": "value"}`)
	if err != nil {
		t.Fatalf("Failed to save long term memory: %v", err)
	}

	memories, _ := mgr.List(MemoryTypeLongTerm, 10)
	if len(memories) != 1 {
		t.Errorf("Expected 1 longterm memory, got %d", len(memories))
	}
}

func TestMemoryManagerGetMemoryStats(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	mgr := NewMemoryManager(store)
	ctx := context.Background()

	mgr.Add(ctx, MemoryTypeSnapshot, "snapshot 1", "{}")
	mgr.Add(ctx, MemoryTypeSnapshot, "snapshot 2", "{}")
	mgr.Add(ctx, MemoryTypeSession, "session 1", "{}")
	mgr.Add(ctx, MemoryTypeLongTerm, "longterm 1", "{}")

	stats := mgr.GetMemoryStats()

	if stats[MemoryTypeSnapshot] != 2 {
		t.Errorf("Expected 2 snapshots, got %d", stats[MemoryTypeSnapshot])
	}

	if stats[MemoryTypeSession] != 1 {
		t.Errorf("Expected 1 session, got %d", stats[MemoryTypeSession])
	}

	if stats[MemoryTypeLongTerm] != 1 {
		t.Errorf("Expected 1 longterm, got %d", stats[MemoryTypeLongTerm])
	}
}
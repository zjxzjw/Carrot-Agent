package storage

import (
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	if store.db == nil {
		t.Fatal("Expected database connection to be initialized")
	}
}

func TestSaveAndGetSession(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	session := &Session{
		ID:        "test_session_1",
		UserID:    "user1",
		Messages:  `{"role": "user", "content": "hello"}`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = store.SaveSession(session)
	if err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	retrieved, err := store.GetSession("test_session_1")
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected session, got nil")
	}

	if retrieved.ID != session.ID {
		t.Errorf("Expected ID %s, got %s", session.ID, retrieved.ID)
	}

	if retrieved.UserID != session.UserID {
		t.Errorf("Expected UserID %s, got %s", session.UserID, retrieved.UserID)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	retrieved, err := store.GetSession("non_existent")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if retrieved != nil {
		t.Error("Expected nil for non-existent session")
	}
}

func TestListSessions(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	for i := 0; i < 3; i++ {
		session := &Session{
			ID:        "session_" + string(rune('a'+i)),
			UserID:    "user1",
			Messages:  "[]",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := store.SaveSession(session); err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}
	}

	sessions, err := store.ListSessions("user1", 10)
	if err != nil {
		t.Fatalf("Failed to list sessions: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sessions))
	}
}

func TestDeleteSession(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	session := &Session{
		ID:        "to_delete",
		UserID:    "user1",
		Messages:  "[]",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := store.SaveSession(session); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	if err := store.DeleteSession("to_delete"); err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	retrieved, _ := store.GetSession("to_delete")
	if retrieved != nil {
		t.Error("Expected session to be deleted")
	}
}

func TestSaveAndGetMemory(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	memory := &Memory{
		ID:        "memory_1",
		Type:      "snapshot",
		Content:   "Test memory content",
		Metadata:  "{}",
		CreatedAt: time.Now(),
	}

	err = store.SaveMemory(memory)
	if err != nil {
		t.Fatalf("Failed to save memory: %v", err)
	}

	retrieved, err := store.GetMemory("memory_1")
	if err != nil {
		t.Fatalf("Failed to get memory: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected memory, got nil")
	}

	if retrieved.Content != memory.Content {
		t.Errorf("Expected content %s, got %s", memory.Content, retrieved.Content)
	}
}

func TestListMemories(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	memTypes := []string{"snapshot", "session", "longterm"}
	for i, memType := range memTypes {
		memory := &Memory{
			ID:        memType + "_1",
			Type:      memType,
			Content:   "Test content",
			Metadata:  "{}",
			CreatedAt: time.Now(),
		}
		if err := store.SaveMemory(memory); err != nil {
			t.Fatalf("Failed to save memory: %v", err)
		}
		_ = i
	}

	snapshots, err := store.ListMemories("snapshot", 10)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot memory, got %d", len(snapshots))
	}
}

func TestSearchMemories(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	memories := []*Memory{
		{ID: "mem1", Type: "snapshot", Content: "golang tutorial", Metadata: "{}", CreatedAt: time.Now()},
		{ID: "mem2", Type: "snapshot", Content: "python tutorial", Metadata: "{}", CreatedAt: time.Now()},
		{ID: "mem3", Type: "longterm", Content: "golang best practices", Metadata: "{}", CreatedAt: time.Now()},
	}

	for _, m := range memories {
		if err := store.SaveMemory(m); err != nil {
			t.Fatalf("Failed to save memory: %v", err)
		}
	}

	results, err := store.SearchMemories("golang", 10)
	if err != nil {
		t.Fatalf("Failed to search memories: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestDeleteMemory(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	memory := &Memory{
		ID:        "to_delete",
		Type:      "snapshot",
		Content:   "Test",
		Metadata:  "{}",
		CreatedAt: time.Now(),
	}

	if err := store.SaveMemory(memory); err != nil {
		t.Fatalf("Failed to save memory: %v", err)
	}

	if err := store.DeleteMemory("to_delete"); err != nil {
		t.Fatalf("Failed to delete memory: %v", err)
	}

	retrieved, _ := store.GetMemory("to_delete")
	if retrieved != nil {
		t.Error("Expected memory to be deleted")
	}
}

func TestSaveAndGetSkill(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	skill := &Skill{
		ID:          "skill_1",
		Name:        "test-skill",
		Description: "A test skill",
		Version:     "1.0.0",
		Platforms:   `["macos","linux"]`,
		Content:     "# Test Skill\n\nThis is a test skill",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = store.SaveSkill(skill)
	if err != nil {
		t.Fatalf("Failed to save skill: %v", err)
	}

	retrieved, err := store.GetSkill("skill_1")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected skill, got nil")
	}

	if retrieved.Name != skill.Name {
		t.Errorf("Expected name %s, got %s", skill.Name, retrieved.Name)
	}
}

func TestGetSkillByName(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	skill := &Skill{
		ID:          "skill_by_name",
		Name:        "unique-skill-name",
		Description: "Test",
		Version:     "1.0.0",
		Platforms:   "[]",
		Content:     "Content",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := store.SaveSkill(skill); err != nil {
		t.Fatalf("Failed to save skill: %v", err)
	}

	retrieved, err := store.GetSkillByName("unique-skill-name")
	if err != nil {
		t.Fatalf("Failed to get skill by name: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected skill, got nil")
	}

	if retrieved.ID != skill.ID {
		t.Errorf("Expected ID %s, got %s", skill.ID, retrieved.ID)
	}
}

func TestListSkills(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	for i := 0; i < 5; i++ {
		skill := &Skill{
			ID:          "skill_" + string(rune('a'+i)),
			Name:        "skill_" + string(rune('a'+i)),
			Description: "Test skill",
			Version:     "1.0.0",
			Platforms:   "[]",
			Content:     "Content",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := store.SaveSkill(skill); err != nil {
			t.Fatalf("Failed to save skill: %v", err)
		}
	}

	skills, err := store.ListSkills(10)
	if err != nil {
		t.Fatalf("Failed to list skills: %v", err)
	}

	if len(skills) != 5 {
		t.Errorf("Expected 5 skills, got %d", len(skills))
	}
}

func TestSearchSkills(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	skills := []*Skill{
		{ID: "sk1", Name: "golang-expert", Description: "Go programming expert", Version: "1.0.0", Platforms: "[]", Content: "Content", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "sk2", Name: "python-dev", Description: "Python developer", Version: "1.0.0", Platforms: "[]", Content: "Content", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "sk3", Name: "golang-tutor", Description: "Teach Go programming", Version: "1.0.0", Platforms: "[]", Content: "Content", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, s := range skills {
		if err := store.SaveSkill(s); err != nil {
			t.Fatalf("Failed to save skill: %v", err)
		}
	}

	results, err := store.SearchSkills("golang", 10)
	if err != nil {
		t.Fatalf("Failed to search skills: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestDeleteSkill(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	skill := &Skill{
		ID:          "to_delete",
		Name:        "delete-me",
		Description: "Test",
		Version:     "1.0.0",
		Platforms:   "[]",
		Content:     "Content",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := store.SaveSkill(skill); err != nil {
		t.Fatalf("Failed to save skill: %v", err)
	}

	if err := store.DeleteSkill("to_delete"); err != nil {
		t.Fatalf("Failed to delete skill: %v", err)
	}

	retrieved, _ := store.GetSkill("to_delete")
	if retrieved != nil {
		t.Error("Expected skill to be deleted")
	}
}
package skill

import (
	"context"
	"testing"

	"carrotagent/carrot-agent/pkg/storage"
)

func TestNewSkillManager(t *testing.T) {
	// 创建内存存储
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 创建技能管理器
	skillManager := NewSkillManager(store)

	// 验证技能管理器创建成功
	if skillManager == nil {
		t.Fatal("Failed to create skill manager")
	}
}

func TestSkillManagerCreate(t *testing.T) {
	// 创建内存存储
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 创建技能管理器
	skillManager := NewSkillManager(store)

	// 初始化技能管理器
	err = skillManager.Load(context.Background())
	if err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	// 创建技能
	skillName := "test-skill"
	skillDescription := "Test skill"
	skillContent := GenerateSkillFile(skillName, skillDescription, "# Test Skill\n\nThis is a test skill")

	err = skillManager.Create(context.Background(), skillName, skillDescription, skillContent)
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	// 验证技能创建成功
	skills := skillManager.List(100)
	if len(skills) != 1 {
		t.Errorf("Expected 1 skill, got %d", len(skills))
	}

	if skills[0].Name != skillName {
		t.Errorf("Expected skill name to be '%s', got '%s'", skillName, skills[0].Name)
	}
}

func TestSkillManagerList(t *testing.T) {
	// 创建内存存储
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 创建技能管理器
	skillManager := NewSkillManager(store)

	// 初始化技能管理器
	err = skillManager.Load(context.Background())
	if err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	// 列出技能
	skills := skillManager.List(100)

	// 验证技能列表为空
	if len(skills) != 0 {
		t.Errorf("Expected 0 skills, got %d", len(skills))
	}
}

func TestSkillManagerGetSkillsIndex(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	err = skillManager.Create(context.Background(), "skill1", "First skill", GenerateSkillFile("skill1", "First skill", "Content 1"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	err = skillManager.Create(context.Background(), "skill2", "Second skill", GenerateSkillFile("skill2", "Second skill", "Content 2"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	index := skillManager.GetSkillsIndex()
	if index == "" {
		t.Error("Expected non-empty skills index")
	}
}

func TestSkillManagerGetSkillCount(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	count := skillManager.GetSkillCount()
	if count != 0 {
		t.Errorf("Expected 0 skills initially, got %d", count)
	}

	err = skillManager.Create(context.Background(), "test-skill", "A test skill", GenerateSkillFile("test-skill", "A test skill", "Content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	count = skillManager.GetSkillCount()
	if count != 1 {
		t.Errorf("Expected 1 skill after creation, got %d", count)
	}
}

func TestSkillManagerSearch(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	err = skillManager.Create(context.Background(), "search-test", "Testing search functionality", GenerateSkillFile("search-test", "Testing search functionality", "Content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	results, err := skillManager.Search("search", 10)
	if err != nil {
		t.Fatalf("Failed to search skills: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(results))
	}
}

func TestSkillManagerDelete(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	err = skillManager.Create(context.Background(), "delete-me", "Will be deleted", GenerateSkillFile("delete-me", "Will be deleted", "Content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	skills := skillManager.List(100)
	if len(skills) != 1 {
		t.Fatalf("Expected 1 skill before deletion, got %d", len(skills))
	}

	skillID := skills[0].ID
	err = skillManager.Delete(context.Background(), skillID)
	if err != nil {
		t.Fatalf("Failed to delete skill: %v", err)
	}

	skills = skillManager.List(100)
	if len(skills) != 0 {
		t.Errorf("Expected 0 skills after deletion, got %d", len(skills))
	}
}

func TestSkillManagerUpdate(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	err = skillManager.Create(context.Background(), "update-test", "Original description", GenerateSkillFile("update-test", "Original description", "Original content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	skills := skillManager.List(100)
	if len(skills) != 1 {
		t.Fatalf("Expected 1 skill, got %d", len(skills))
	}

	skillID := skills[0].ID
	newContent := GenerateSkillFile("update-test", "Updated description", "Updated content")
	err = skillManager.Update(context.Background(), skillID, newContent)
	if err != nil {
		t.Fatalf("Failed to update skill: %v", err)
	}

	updatedSkill, err := skillManager.Get(skillID)
	if err != nil {
		t.Fatalf("Failed to get updated skill: %v", err)
	}

	if updatedSkill.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", updatedSkill.Description)
	}
}

func TestParseSkillContent(t *testing.T) {
	content := `---
name: test-skill
description: A test skill
version: 1.0.0
platforms: [macos, linux]
---

# Skill Content
This is the actual skill content.`

	result, err := ParseSkillContent(content)
	if err != nil {
		t.Fatalf("Failed to parse skill content: %v", err)
	}

	if result["name"] != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%v'", result["name"])
	}

	if result["description"] != "A test skill" {
		t.Errorf("Expected description 'A test skill', got '%v'", result["description"])
	}

	if result["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", result["version"])
	}
}

func TestParseSkillContentWithContent(t *testing.T) {
	content := `---
name: with-content
description: Skill with content
---

# Main Content
This is the main content after frontmatter.`

	result, err := ParseSkillContent(content)
	if err != nil {
		t.Fatalf("Failed to parse skill content: %v", err)
	}

	if result["_content"] == "" {
		t.Error("Expected _content to be non-empty")
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"valid-name", "valid-name"},
		{"ValidName123", "ValidName123"},
		{"name with spaces", "name_with_spaces"},
		{"name@#$%", "name____"},
		{"name-with-dash", "name-with-dash"},
		{"name_with_underscore", "name_with_underscore"},
	}

	for _, test := range tests {
		result := sanitizeName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateSkillFile(t *testing.T) {
	result := GenerateSkillFile("test-skill", "A test skill", "# Test\n\nSkill content here.")

	if result == "" {
		t.Error("Expected non-empty skill file content")
	}

	if !contains(result, "name: test-skill") {
		t.Error("Expected skill file to contain 'name: test-skill'")
	}

	if !contains(result, "description: A test skill") {
		t.Error("Expected skill file to contain 'description: A test skill'")
	}

	if !contains(result, "# Test") {
		t.Error("Expected skill file to contain '# Test'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSkillManagerGetByName(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	skillName := "unique-test-skill"
	err = skillManager.Create(context.Background(), skillName, "Testing GetByName", GenerateSkillFile(skillName, "Testing GetByName", "Content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	skill, err := skillManager.GetByName(skillName)
	if err != nil {
		t.Fatalf("Failed to get skill by name: %v", err)
	}

	if skill == nil {
		t.Fatal("Expected to find skill by name")
	}

	if skill.Name != skillName {
		t.Errorf("Expected skill name '%s', got '%s'", skillName, skill.Name)
	}
}

func TestSkillManagerGetByNameNotFound(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	skill, err := skillManager.GetByName("non-existent-skill")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if skill != nil {
		t.Error("Expected nil for non-existent skill")
	}
}

func TestSkillManagerPatch(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	err = skillManager.Create(context.Background(), "patch-test", "Original", GenerateSkillFile("patch-test", "Original", "Original content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	skills := skillManager.List(100)
	if len(skills) != 1 {
		t.Fatalf("Expected 1 skill, got %d", len(skills))
	}

	skillID := skills[0].ID
	err = skillManager.Patch(context.Background(), skillID, "Original", "Patched")
	if err != nil {
		t.Fatalf("Failed to patch skill: %v", err)
	}

	updatedSkill, err := skillManager.Get(skillID)
	if err != nil {
		t.Fatalf("Failed to get patched skill: %v", err)
	}

	if !contains(updatedSkill.Content, "Patched") {
		t.Error("Expected patched content to contain 'Patched'")
	}
}

func TestSkillManagerPatchTextNotFound(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	skillManager := NewSkillManager(store)
	if err := skillManager.Load(context.Background()); err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	err = skillManager.Create(context.Background(), "patch-test-2", "Original", GenerateSkillFile("patch-test-2", "Original", "Original content"))
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	skills := skillManager.List(100)
	skillID := skills[0].ID

	err = skillManager.Patch(context.Background(), skillID, "NonExistent", "Replacement")
	if err == nil {
		t.Error("Expected error when patching with non-existent text")
	}
}

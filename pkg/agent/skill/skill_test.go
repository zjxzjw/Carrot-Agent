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

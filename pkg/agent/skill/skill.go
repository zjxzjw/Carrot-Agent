package skill

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"carrotagent/carrot-agent/pkg/storage"
)

type SkillManager struct {
	store  *storage.Store
	skills map[string]*storage.Skill
}

type SkillMetadata struct {
	Hermes          HermesMetadata `yaml:"hermes"`
	FallbackForTool string         `yaml:"fallback_for_toolsets"`
	RequiresTool    string         `yaml:"requires_toolsets"`
}

type HermesMetadata struct {
	Tags     []string `yaml:"tags"`
	Category string   `yaml:"category"`
}

func NewSkillManager(store *storage.Store) *SkillManager {
	return &SkillManager{
		store:  store,
		skills: make(map[string]*storage.Skill),
	}
}

func (m *SkillManager) Load(ctx context.Context) error {
	skills, err := m.store.ListSkills(1000)
	if err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}

	for _, skill := range skills {
		m.skills[skill.ID] = skill
	}

	return nil
}

func (m *SkillManager) Create(ctx context.Context, name, description, content string) error {
	id := fmt.Sprintf("skill_%s_%d", sanitizeName(name), time.Now().UnixNano())

	skillFile, err := ParseSkillContent(content)
	if err != nil {
		return fmt.Errorf("failed to parse skill content: %w", err)
	}

	version := "1.0.0"
	if v, ok := skillFile["version"].(string); ok {
		version = v
	}

	platforms := "[]"
	if p, ok := skillFile["platforms"].([]string); ok {
		platforms = "[\"" + strings.Join(p, "\",\"") + "\"]"
	} else if p, ok := skillFile["platforms"].(string); ok {
		platforms = p
	}

	skill := &storage.Skill{
		ID:          id,
		Name:        name,
		Description: description,
		Version:     version,
		Platforms:   platforms,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := m.store.SaveSkill(skill); err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	m.skills[id] = skill
	return nil
}

func (m *SkillManager) Update(ctx context.Context, id string, content string) error {
	skill, ok := m.skills[id]
	if !ok {
		var err error
		skill, err = m.store.GetSkill(id)
		if err != nil {
			return err
		}
		if skill == nil {
			return fmt.Errorf("skill not found: %s", id)
		}
	}

	skillFile, err := ParseSkillContent(content)
	if err != nil {
		return fmt.Errorf("failed to parse skill content: %w", err)
	}

	if name, ok := skillFile["name"].(string); ok {
		skill.Name = name
	}
	if desc, ok := skillFile["description"].(string); ok {
		skill.Description = desc
	}
	if v, ok := skillFile["version"].(string); ok {
		skill.Version = v
	}

	skill.Content = content
	skill.UpdatedAt = time.Now()

	if err := m.store.SaveSkill(skill); err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	m.skills[id] = skill
	return nil
}

func (m *SkillManager) Patch(ctx context.Context, id, oldText, newText string) error {
	skill, ok := m.skills[id]
	if !ok {
		var err error
		skill, err = m.store.GetSkill(id)
		if err != nil {
			return err
		}
		if skill == nil {
			return fmt.Errorf("skill not found: %s", id)
		}
	}

	if !strings.Contains(skill.Content, oldText) {
		return fmt.Errorf("old text not found in skill content")
	}

	skill.Content = strings.Replace(skill.Content, oldText, newText, 1)
	skill.UpdatedAt = time.Now()

	if err := m.store.SaveSkill(skill); err != nil {
		return fmt.Errorf("failed to patch skill: %w", err)
	}

	m.skills[id] = skill
	return nil
}

func (m *SkillManager) Delete(ctx context.Context, id string) error {
	if err := m.store.DeleteSkill(id); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	delete(m.skills, id)
	return nil
}

func (m *SkillManager) Get(id string) (*storage.Skill, error) {
	skill, ok := m.skills[id]
	if !ok {
		return m.store.GetSkill(id)
	}
	return skill, nil
}

func (m *SkillManager) GetByName(name string) (*storage.Skill, error) {
	for _, skill := range m.skills {
		if skill.Name == name {
			return skill, nil
		}
	}

	return m.store.GetSkillByName(name)
}

func (m *SkillManager) List(limit int) []*storage.Skill {
	if limit <= 0 {
		limit = 100
	}

	skills, err := m.store.ListSkills(limit)
	if err != nil {
		return nil
	}

	return skills
}

func (m *SkillManager) Search(keyword string, limit int) ([]*storage.Skill, error) {
	if limit <= 0 {
		limit = 50
	}

	return m.store.SearchSkills(keyword, limit)
}

func (m *SkillManager) GetSkillsIndex() string {
	var index []string
	index = append(index, "# Skills Index\n")

	for _, skill := range m.skills {
		index = append(index, fmt.Sprintf("## %s\n%s\n", skill.Name, skill.Description))
	}

	return strings.Join(index, "\n")
}

func (m *SkillManager) GetSkillCount() int {
	return len(m.skills)
}

func ParseSkillContent(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	lines := strings.Split(content, "\n")

	inFrontmatter := false
	frontmatterLines := []string{}
	contentLines := []string{}

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				contentLines = lines[i+1:]
				break
			}
		}

		if inFrontmatter {
			frontmatterLines = append(frontmatterLines, line)
		} else {
			contentLines = append(contentLines, line)
		}
	}

	for _, line := range frontmatterLines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			value = strings.Trim(value, "\"")

			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				value = strings.Trim(value, "[]")
				items := strings.Split(value, ",")
				var cleanItems []string
				for _, item := range items {
					cleanItems = append(cleanItems, strings.Trim(strings.TrimSpace(item), "\""))
				}
				result[key] = cleanItems
			} else {
				result[key] = value
			}
		}
	}

	result["_content"] = strings.Join(contentLines, "\n")

	return result, nil
}

func GenerateSkillFile(name, description, content string) string {
	return fmt.Sprintf(`---
name: %s
description: %s
version: 1.0.0
platforms: [macos, linux, docker]
metadata:
  hermes:
    tags: []
    category: custom
---

%s
`, name, description, content)
}

func sanitizeName(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return reg.ReplaceAllString(name, "_")
}
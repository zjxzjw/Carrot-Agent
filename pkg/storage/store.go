package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Messages  string    `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Memory struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
}

type Skill struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Platforms   string    `json:"platforms"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// 优化SQLite配置
	db.SetMaxOpenConns(1) // SQLite只支持一个写连接
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 设置PRAGMA优化
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=-64000", // 64MB cache
		"PRAGMA temp_store=MEMORY",
		"PRAGMA mmap_size=268435456", // 256MB
	}
	
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	store := &Store{db: db}
	if err := store.init(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return store, nil
}

func (s *Store) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		messages TEXT NOT NULL DEFAULT '[]',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS memories (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		metadata TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS skills (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		version TEXT NOT NULL DEFAULT '1.0.0',
		platforms TEXT NOT NULL DEFAULT '[]',
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_memories_type ON memories(type);
	CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);
	`

	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Ping 检查数据库连接是否正常
func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) SaveSession(session *Session) error {
	query := `
	INSERT INTO sessions (id, user_id, messages, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		messages = excluded.messages,
		updated_at = excluded.updated_at
	`

	_, err := s.db.Exec(query, session.ID, session.UserID, session.Messages, session.CreatedAt, session.UpdatedAt)
	return err
}

func (s *Store) GetSession(id string) (*Session, error) {
	query := `SELECT id, user_id, messages, created_at, updated_at FROM sessions WHERE id = ?`

	var session Session
	err := s.db.QueryRow(query, id).Scan(
		&session.ID, &session.UserID, &session.Messages, &session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) ListSessions(userID string, limit int) ([]*Session, error) {
	query := `SELECT id, user_id, messages, created_at, updated_at FROM sessions WHERE user_id = ? ORDER BY updated_at DESC LIMIT ?`

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.ID, &session.UserID, &session.Messages, &session.CreatedAt, &session.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (s *Store) DeleteSession(id string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

func (s *Store) SaveMemory(memory *Memory) error {
	query := `
	INSERT INTO memories (id, type, content, metadata, created_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		content = excluded.content,
		metadata = excluded.metadata
	`

	_, err := s.db.Exec(query, memory.ID, memory.Type, memory.Content, memory.Metadata, memory.CreatedAt)
	return err
}

func (s *Store) GetMemory(id string) (*Memory, error) {
	query := `SELECT id, type, content, metadata, created_at FROM memories WHERE id = ?`

	var memory Memory
	err := s.db.QueryRow(query, id).Scan(&memory.ID, &memory.Type, &memory.Content, &memory.Metadata, &memory.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &memory, nil
}

func (s *Store) ListMemories(memType string, limit int) ([]*Memory, error) {
	query := `SELECT id, type, content, metadata, created_at FROM memories WHERE type = ? ORDER BY created_at DESC LIMIT ?`

	rows, err := s.db.Query(query, memType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*Memory
	for rows.Next() {
		var memory Memory
		if err := rows.Scan(&memory.ID, &memory.Type, &memory.Content, &memory.Metadata, &memory.CreatedAt); err != nil {
			return nil, err
		}
		memories = append(memories, &memory)
	}

	return memories, nil
}

func (s *Store) DeleteMemory(id string) error {
	_, err := s.db.Exec(`DELETE FROM memories WHERE id = ?`, id)
	return err
}

func (s *Store) SaveSkill(skill *Skill) error {
	query := `
	INSERT INTO skills (id, name, description, version, platforms, content, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		description = excluded.description,
		version = excluded.version,
		platforms = excluded.platforms,
		content = excluded.content,
		updated_at = excluded.updated_at
	`

	_, err := s.db.Exec(query, skill.ID, skill.Name, skill.Description, skill.Version, skill.Platforms, skill.Content, skill.CreatedAt, skill.UpdatedAt)
	return err
}

func (s *Store) GetSkill(id string) (*Skill, error) {
	query := `SELECT id, name, description, version, platforms, content, created_at, updated_at FROM skills WHERE id = ?`

	var skill Skill
	err := s.db.QueryRow(query, id).Scan(
		&skill.ID, &skill.Name, &skill.Description, &skill.Version, &skill.Platforms, &skill.Content, &skill.CreatedAt, &skill.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &skill, nil
}

func (s *Store) GetSkillByName(name string) (*Skill, error) {
	query := `SELECT id, name, description, version, platforms, content, created_at, updated_at FROM skills WHERE name = ?`

	var skill Skill
	err := s.db.QueryRow(query, name).Scan(
		&skill.ID, &skill.Name, &skill.Description, &skill.Version, &skill.Platforms, &skill.Content, &skill.CreatedAt, &skill.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &skill, nil
}

func (s *Store) ListSkills(limit int) ([]*Skill, error) {
	query := `SELECT id, name, description, version, platforms, content, created_at, updated_at FROM skills ORDER BY updated_at DESC LIMIT ?`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(
			&skill.ID, &skill.Name, &skill.Description, &skill.Version, &skill.Platforms, &skill.Content, &skill.CreatedAt, &skill.UpdatedAt,
		); err != nil {
			return nil, err
		}
		skills = append(skills, &skill)
	}

	return skills, nil
}

func (s *Store) DeleteSkill(id string) error {
	_, err := s.db.Exec(`DELETE FROM skills WHERE id = ?`, id)
	return err
}

func (s *Store) SearchSkills(keyword string, limit int) ([]*Skill, error) {
	query := `SELECT id, name, description, version, platforms, content, created_at, updated_at FROM skills WHERE name LIKE ? OR description LIKE ? ORDER BY updated_at DESC LIMIT ?`

	pattern := "%" + keyword + "%"
	rows, err := s.db.Query(query, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(
			&skill.ID, &skill.Name, &skill.Description, &skill.Version, &skill.Platforms, &skill.Content, &skill.CreatedAt, &skill.UpdatedAt,
		); err != nil {
			return nil, err
		}
		skills = append(skills, &skill)
	}

	return skills, nil
}

func (s *Store) SearchMemories(keyword string, limit int) ([]*Memory, error) {
	query := `SELECT id, type, content, metadata, created_at FROM memories WHERE content LIKE ? ORDER BY created_at DESC LIMIT ?`

	pattern := "%" + keyword + "%"
	rows, err := s.db.Query(query, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*Memory
	for rows.Next() {
		var memory Memory
		if err := rows.Scan(&memory.ID, &memory.Type, &memory.Content, &memory.Metadata, &memory.CreatedAt); err != nil {
			return nil, err
		}
		memories = append(memories, &memory)
	}

	return memories, nil
}
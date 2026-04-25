package tool

import (
	"context"
	"testing"
)

func TestNewToolRegistry(t *testing.T) {
	registry := NewToolRegistry()

	if registry == nil {
		t.Fatal("Expected registry, got nil")
	}

	if registry.tools == nil {
		t.Error("Expected tools map to be initialized")
	}

	if registry.defs == nil {
		t.Error("Expected defs map to be initialized")
	}
}

func TestToolRegistryRegister(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "test result", nil
	}

	registry.Register("test_tool", "A test tool", map[string]interface{}{
		"param1": map[string]interface{}{"type": "string", "description": "Test param", "required": true},
	}, handler)

	if len(registry.tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(registry.tools))
	}

	if len(registry.defs) != 1 {
		t.Errorf("Expected 1 def, got %d", len(registry.defs))
	}

	tool, ok := registry.Get("test_tool")
	if !ok {
		t.Error("Expected to get registered tool")
	}

	if tool == nil {
		t.Error("Expected tool handler, got nil")
	}
}

func TestToolRegistryGet(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "result", nil
	}

	registry.Register("my_tool", "My tool", nil, handler)

	tool, ok := registry.Get("my_tool")
	if !ok {
		t.Error("Expected to get tool")
	}

	if tool == nil {
		t.Error("Expected tool handler")
	}
}

func TestToolRegistryGetNonExistent(t *testing.T) {
	registry := NewToolRegistry()

	_, ok := registry.Get("non_existent")
	if ok {
		t.Error("Expected ok to be false for non-existent tool")
	}
}

func TestToolRegistryGetDef(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return nil, nil
	}

	registry.Register("my_tool", "My tool description", map[string]interface{}{
		"arg1": map[string]interface{}{"type": "number", "description": "An argument"},
	}, handler)

	def, ok := registry.GetDef("my_tool")
	if !ok {
		t.Error("Expected to get tool def")
	}

	if def == nil {
		t.Error("Expected def, got nil")
	}

	if def.Name != "my_tool" {
		t.Errorf("Expected Name 'my_tool', got '%s'", def.Name)
	}

	if def.Description != "My tool description" {
		t.Errorf("Expected Description 'My tool description', got '%s'", def.Description)
	}
}

func TestToolRegistryGetDefNonExistent(t *testing.T) {
	registry := NewToolRegistry()

	_, ok := registry.GetDef("non_existent")
	if ok {
		t.Error("Expected ok to be false for non-existent def")
	}
}

func TestToolRegistryListDefs(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return nil, nil
	}

	registry.Register("tool1", "Tool 1", nil, handler)
	registry.Register("tool2", "Tool 2", nil, handler)
	registry.Register("tool3", "Tool 3", nil, handler)

	defs := registry.ListDefs()
	if len(defs) != 3 {
		t.Errorf("Expected 3 defs, got %d", len(defs))
	}
}

func TestToolRegistryExecute(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		if name, ok := args["name"].(string); ok {
			return "Hello, " + name, nil
		}
		return "Hello", nil
	}

	registry.Register("greet", "Greet someone", map[string]interface{}{
		"name": map[string]interface{}{"type": "string", "description": "Name to greet"},
	}, handler)

	result := registry.Execute(context.Background(), "greet", map[string]interface{}{"name": "World"})

	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.Error)
	}

	if result.Output != "Hello, World" {
		t.Errorf("Expected 'Hello, World', got '%v'", result.Output)
	}
}

func TestToolRegistryExecuteNonExistent(t *testing.T) {
	registry := NewToolRegistry()

	result := registry.Execute(context.Background(), "non_existent", nil)

	if result.Success {
		t.Error("Expected failure for non-existent tool")
	}

	if result.Error != "tool not found: non_existent" {
		t.Errorf("Expected 'tool not found: non_existent', got '%s'", result.Error)
	}
}

func TestToolRegistryExecuteWithError(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return nil, context.DeadlineExceeded
	}

	registry.Register("failing_tool", "A failing tool", nil, handler)

	result := registry.Execute(context.Background(), "failing_tool", nil)

	if result.Success {
		t.Error("Expected failure for tool that returns error")
	}

	if result.Error != context.DeadlineExceeded.Error() {
		t.Errorf("Expected '%v', got '%s'", context.DeadlineExceeded, result.Error)
	}
}

func TestToolRegistryGetToolsForPrompt(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return nil, nil
	}

	registry.Register("tool1", "First tool", map[string]interface{}{
		"param1": map[string]interface{}{"type": "string", "description": "Required param", "required": true},
		"param2": map[string]interface{}{"type": "number", "description": "Optional param", "required": false},
	}, handler)

	output := registry.GetToolsForPrompt()

	if output == "" {
		t.Error("Expected non-empty output")
	}
}

func TestToolRegistryGetToolsForPromptEmpty(t *testing.T) {
	registry := NewToolRegistry()

	output := registry.GetToolsForPrompt()

	if output != "No tools available." {
		t.Errorf("Expected 'No tools available.', got '%s'", output)
	}
}

func TestConvertToModelTools(t *testing.T) {
	registry := NewToolRegistry()

	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return nil, nil
	}

	registry.Register("test_tool", "A test tool", map[string]interface{}{
		"param1": map[string]interface{}{"type": "string", "description": "A string param", "required": true},
		"param2": map[string]interface{}{"type": "number", "description": "A number param", "required": false},
	}, handler)

	tools := ConvertToModelTools(registry)

	if len(tools) != 1 {
		t.Fatalf("Expected 1 tool, got %d", len(tools))
	}

	tool := tools[0]
	if tool["type"] != "function" {
		t.Errorf("Expected type 'function', got '%v'", tool["type"])
	}

	fn, ok := tool["function"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected function in tool")
	}

	if fn["name"] != "test_tool" {
		t.Errorf("Expected name 'test_tool', got '%v'", fn["name"])
	}

	if fn["description"] != "A test tool" {
		t.Errorf("Expected description 'A test tool', got '%v'", fn["description"])
	}
}

func TestConvertToModelToolsEmpty(t *testing.T) {
	registry := NewToolRegistry()

	tools := ConvertToModelTools(registry)

	if len(tools) != 0 {
		t.Errorf("Expected 0 tools, got %d", len(tools))
	}
}

func TestParseToolCall(t *testing.T) {
	toolCall := map[string]interface{}{
		"id":         "call_123",
		"name":       "test_function",
		"arguments":   `{"param1": "value1"}`,
	}

	parsed, err := ParseToolCall(toolCall)
	if err != nil {
		t.Fatalf("Failed to parse tool call: %v", err)
	}

	if parsed.Name != "test_function" {
		t.Errorf("Expected name 'test_function', got '%s'", parsed.Name)
	}

	if parsed.Arguments["param1"] != "value1" {
		t.Errorf("Expected param1 'value1', got '%v'", parsed.Arguments["param1"])
	}
}

func TestParseToolCallWithMapArguments(t *testing.T) {
	toolCall := map[string]interface{}{
		"id": "call_456",
		"name": "another_function",
		"arguments": map[string]interface{}{
			"arg1": "val1",
			"arg2": float64(42),
		},
	}

	parsed, err := ParseToolCall(toolCall)
	if err != nil {
		t.Fatalf("Failed to parse tool call: %v", err)
	}

	if parsed.Name != "another_function" {
		t.Errorf("Expected name 'another_function', got '%s'", parsed.Name)
	}

	if parsed.Arguments["arg1"] != "val1" {
		t.Errorf("Expected arg1 'val1', got '%v'", parsed.Arguments["arg1"])
	}

	if parsed.Arguments["arg2"] != float64(42) {
		t.Errorf("Expected arg2 42, got '%v'", parsed.Arguments["arg2"])
	}
}

func TestParseToolCallWithStringArguments(t *testing.T) {
	toolCall := map[string]interface{}{
		"id":         "call_789",
		"function":   map[string]interface{}{"name": "func_with_string_args"},
		"arguments":  `{"key": "value"}`,
	}

	parsed, err := ParseToolCall(toolCall)
	if err != nil {
		t.Fatalf("Failed to parse tool call: %v", err)
	}

	if parsed.Arguments["key"] != "value" {
		t.Errorf("Expected key 'value', got '%v'", parsed.Arguments["key"])
	}
}

func TestParseToolCallInvalidFormat(t *testing.T) {
	toolCall := map[string]interface{}{
		"id": "call_invalid",
	}

	_, err := ParseToolCall(toolCall)
	if err == nil {
		t.Error("Expected error for invalid format")
	}

	expected := "invalid tool call format"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestParseToolCallMissingFunctionName(t *testing.T) {
	toolCall := map[string]interface{}{
		"id": "call_no_name",
		"function": map[string]interface{}{},
	}

	_, err := ParseToolCall(toolCall)
	if err == nil {
		t.Error("Expected error for missing function name")
	}

	expected := "missing function name"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestParseToolCallInvalidJSON(t *testing.T) {
	toolCall := map[string]interface{}{
		"id":         "call_bad_json",
		"name":       "bad_json_func",
		"arguments":  `{invalid json`,
	}

	_, err := ParseToolCall(toolCall)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestToolCallStruct(t *testing.T) {
	call := ToolCall{
		Name: "test_call",
		Arguments: map[string]interface{}{
			"arg1": "value1",
		},
	}

	if call.Name != "test_call" {
		t.Errorf("Expected Name 'test_call', got '%s'", call.Name)
	}

	if call.Arguments["arg1"] != "value1" {
		t.Errorf("Expected arg1 'value1', got '%v'", call.Arguments["arg1"])
	}
}

func TestToolResultStruct(t *testing.T) {
	result := ToolResult{
		Success:  true,
		Output:   "test output",
		Error:    "",
		Metadata: map[string]interface{}{"key": "value"},
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.Output != "test output" {
		t.Errorf("Expected Output 'test output', got '%v'", result.Output)
	}

	if result.Metadata["key"] != "value" {
		t.Errorf("Expected Metadata key 'value', got '%v'", result.Metadata["key"])
	}
}

func TestToolDefStruct(t *testing.T) {
	def := ToolDef{
		Name:        "my_def",
		Description: "My def description",
		Parameters: map[string]interface{}{
			"param1": map[string]interface{}{"type": "string"},
		},
	}

	if def.Name != "my_def" {
		t.Errorf("Expected Name 'my_def', got '%s'", def.Name)
	}

	if def.Description != "My def description" {
		t.Errorf("Expected Description 'My def description', got '%s'", def.Description)
	}

	if def.Parameters["param1"] == nil {
		t.Error("Expected param1 to be set")
	}
}
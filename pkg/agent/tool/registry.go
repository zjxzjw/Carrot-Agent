package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Tool func(ctx context.Context, args map[string]interface{}) (interface{}, error)

type ToolDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ToolResult struct {
	Success bool                   `json:"success"`
	Output  interface{}            `json:"output,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type ToolRegistry struct {
	tools map[string]Tool
	defs  map[string]*ToolDef
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
		defs:  make(map[string]*ToolDef),
	}
}

func (r *ToolRegistry) Register(name string, description string, params map[string]interface{}, handler Tool) {
	r.tools[name] = handler
	r.defs[name] = &ToolDef{
		Name:        name,
		Description: description,
		Parameters:  params,
	}
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *ToolRegistry) GetDef(name string) (*ToolDef, bool) {
	d, ok := r.defs[name]
	return d, ok
}

func (r *ToolRegistry) ListDefs() []*ToolDef {
	defs := make([]*ToolDef, 0, len(r.defs))
	for _, d := range r.defs {
		defs = append(defs, d)
	}
	return defs
}

func (r *ToolRegistry) Execute(ctx context.Context, name string, args map[string]interface{}) *ToolResult {
	tool, ok := r.tools[name]
	if !ok {
		return &ToolResult{
			Success: false,
			Error:   fmt.Sprintf("tool not found: %s", name),
		}
	}

	output, err := tool(ctx, args)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &ToolResult{
		Success: true,
		Output:  output,
	}
}

func (r *ToolRegistry) GetToolsForPrompt() string {
	defs := r.ListDefs()
	if len(defs) == 0 {
		return "No tools available."
	}

	var lines []string
	lines = append(lines, "## Available Tools\n")

	for _, def := range defs {
		lines = append(lines, fmt.Sprintf("### %s\n%s\n", def.Name, def.Description))

		if len(def.Parameters) > 0 {
			lines = append(lines, "Parameters:")
			for paramName, paramDef := range def.Parameters {
				paramMap, ok := paramDef.(map[string]interface{})
				if !ok {
					continue
				}
				paramType, _ := paramMap["type"].(string)
				paramDesc, _ := paramMap["description"].(string)
				required, _ := paramMap["required"].(bool)

				reqStr := ""
				if required {
					reqStr = " (required)"
				}

				lines = append(lines, fmt.Sprintf("- %s: %s (%s)%s", paramName, paramDesc, paramType, reqStr))
			}
		}
		lines = append(lines, "")
	}

	return strings.Join(lines, "")
}

func ConvertToModelTools(registry *ToolRegistry) []map[string]interface{} {
	defs := registry.ListDefs()
	tools := make([]map[string]interface{}, 0, len(defs))

	for _, def := range defs {
		tool := map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        def.Name,
				"description": def.Description,
			},
		}

		// Always add parameters field, even for tools with no parameters
		properties := make(map[string]interface{})
		required := []string{}

		if len(def.Parameters) > 0 {
			for paramName, paramDef := range def.Parameters {
				paramMap, ok := paramDef.(map[string]interface{})
				if !ok {
					continue
				}

				paramType, _ := paramMap["type"].(string)
				paramDesc, _ := paramMap["description"].(string)

				properties[paramName] = map[string]interface{}{
					"type":        paramType,
					"description": paramDesc,
				}

				if requiredFlag, ok := paramMap["required"].(bool); ok && requiredFlag {
					required = append(required, paramName)
				}
			}
		}

		parameters := map[string]interface{}{
			"type":       "object",
			"properties": properties,
			"required":   required,
		}
		tool["function"].(map[string]interface{})["parameters"] = parameters

		tools = append(tools, tool)
	}

	return tools
}

type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func ParseToolCall(toolCall map[string]interface{}) (*ToolCall, error) {
	name, ok := toolCall["name"].(string)
	if !ok {
		funcData, ok := toolCall["function"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid tool call format")
		}
		name, ok = funcData["name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing function name")
		}
	}

	var args map[string]interface{}
	if argsStr, ok := toolCall["arguments"].(string); ok {
		if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
	} else if argsMap, ok := toolCall["arguments"].(map[string]interface{}); ok {
		args = argsMap
	} else {
		args = make(map[string]interface{})
	}

	return &ToolCall{
		Name:      name,
		Arguments: args,
	}, nil
}
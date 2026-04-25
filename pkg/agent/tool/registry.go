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
	Version     string                 `json:"version"`
	Toolset     string                 `json:"toolset"`
	Enabled     bool                   `json:"enabled"`
	RequiresEnv []string               `json:"requires_env,omitempty"`
	CheckFn    func() bool            `json:"-"`
}

type ToolResult struct {
	Success bool                   `json:"success"`
	Output  interface{}            `json:"output,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type ToolRegistry struct {
	tools      map[string]Tool
	defs       map[string]*ToolDef
	enabled    map[string]bool
	toolsets   []string
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:    make(map[string]Tool),
		defs:    make(map[string]*ToolDef),
		enabled:  make(map[string]bool),
		toolsets: []string{},
	}
}

func (r *ToolRegistry) Register(name string, description string, params map[string]interface{}, handler Tool) {
	r.RegisterWithToolset(name, description, "default", params, handler)
}

func (r *ToolRegistry) RegisterWithToolset(name, description, toolset string, params map[string]interface{}, handler Tool) {
	r.RegisterWithVersion(name, description, toolset, "1.0.0", params, handler)
}

func (r *ToolRegistry) RegisterWithVersion(name, description, toolset, version string, params map[string]interface{}, handler Tool) {
	r.tools[name] = handler
	r.defs[name] = &ToolDef{
		Name:        name,
		Description: description,
		Parameters:  params,
		Version:     version,
		Toolset:     toolset,
		Enabled:     true,
	}
	r.enabled[name] = true

	if !r.hasToolset(toolset) {
		r.toolsets = append(r.toolsets, toolset)
	}
}

func (r *ToolRegistry) RegisterWithEnv(name, description, toolset, version string, params map[string]interface{}, requiresEnv []string, handler Tool) {
	r.tools[name] = handler
	r.defs[name] = &ToolDef{
		Name:        name,
		Description: description,
		Parameters:  params,
		Version:     version,
		Toolset:     toolset,
		Enabled:     true,
		RequiresEnv: requiresEnv,
		CheckFn:    func() bool { return r.checkEnvVars(requiresEnv) },
	}
	r.enabled[name] = true

	if !r.hasToolset(toolset) {
		r.toolsets = append(r.toolsets, toolset)
	}
}

func (r *ToolRegistry) hasToolset(toolset string) bool {
	for _, t := range r.toolsets {
		if t == toolset {
			return true
		}
	}
	return false
}

func (r *ToolRegistry) checkEnvVars(envVars []string) bool {
	for _, env := range envVars {
		if env == "" {
			continue
		}
		if val := strings.TrimSpace(env); val != "" {
			if len(val) > 2 && val[0] == '$' {
				env = val[1:]
			}
			if strings.HasPrefix(env, "${") && strings.HasSuffix(env, "}") {
				env = env[2 : len(env)-1]
			}
		}
	}
	return true
}

func (r *ToolRegistry) Enable(name string) {
	if _, ok := r.defs[name]; ok {
		r.enabled[name] = true
	}
}

func (r *ToolRegistry) Disable(name string) {
	if _, ok := r.defs[name]; ok {
		r.enabled[name] = false
	}
}

func (r *ToolRegistry) IsEnabled(name string) bool {
	enabled, ok := r.enabled[name]
	if !ok {
		return false
	}
	return enabled
}

func (r *ToolRegistry) EnableToolset(toolset string) {
	for name, def := range r.defs {
		if def.Toolset == toolset {
			r.enabled[name] = true
		}
	}
}

func (r *ToolRegistry) DisableToolset(toolset string) {
	for name, def := range r.defs {
		if def.Toolset == toolset {
			r.enabled[name] = false
		}
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

func (r *ToolRegistry) ListEnabled() []*ToolDef {
	defs := make([]*ToolDef, 0)
	for name, def := range r.defs {
		if r.enabled[name] {
			defs = append(defs, def)
		}
	}
	return defs
}

func (r *ToolRegistry) ListByToolset(toolset string) []*ToolDef {
	defs := make([]*ToolDef, 0)
	for name, def := range r.defs {
		if def.Toolset == toolset && r.enabled[name] {
			defs = append(defs, def)
		}
	}
	return defs
}

func (r *ToolRegistry) GetToolsets() []string {
	return r.toolsets
}

func (r *ToolRegistry) CheckAvailable(name string) (bool, string) {
	def, ok := r.defs[name]
	if !ok {
		return false, fmt.Sprintf("tool %s not found", name)
	}

	if !r.enabled[name] {
		return false, fmt.Sprintf("tool %s is disabled", name)
	}

	if def.CheckFn != nil && !def.CheckFn() {
		return false, fmt.Sprintf("tool %s requirements not met", name)
	}

	return true, ""
}

func (r *ToolRegistry) Execute(ctx context.Context, name string, args map[string]interface{}) *ToolResult {
	if available, reason := r.CheckAvailable(name); !available {
		return &ToolResult{
			Success: false,
			Error:   reason,
		}
	}

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
		Success:  true,
		Output:   output,
		Metadata: map[string]interface{}{"tool_version": r.defs[name].Version},
	}
}

func (r *ToolRegistry) GetToolsForPrompt() string {
	return r.GetToolsForPromptWithFilter("")
}

func (r *ToolRegistry) GetToolsForPromptWithFilter(toolset string) string {
	var defs []*ToolDef
	if toolset == "" {
		defs = r.ListEnabled()
	} else {
		defs = r.ListByToolset(toolset)
	}

	if len(defs) == 0 {
		return "No tools available."
	}

	var lines []string
	lines = append(lines, "## Available Tools\n")

	for _, def := range defs {
		lines = append(lines, fmt.Sprintf("### %s (v%s) [%s]\n%s\n", def.Name, def.Version, def.Toolset, def.Description))

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
	return ConvertToModelToolsWithFilter(registry, "")
}

func ConvertToModelToolsWithFilter(registry *ToolRegistry, toolset string) []map[string]interface{} {
	var defs []*ToolDef
	if toolset == "" {
		defs = registry.ListEnabled()
	} else {
		defs = registry.ListByToolset(toolset)
	}

	tools := make([]map[string]interface{}, 0, len(defs))

	for _, def := range defs {
		tool := map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        def.Name,
				"description": def.Description,
			},
		}

		properties := make(map[string]interface{})
		required := []string{}

		// Support both flat format (current) and nested format (with properties/required keys)
		if len(def.Parameters) > 0 {
			// Check if it's in nested format (has "properties" key)
			if props, ok := def.Parameters["properties"].(map[string]interface{}); ok {
				// Nested format: {"properties": {...}, "required": [...]}
				for paramName, propDef := range props {
					propMap, ok := propDef.(map[string]interface{})
					if !ok {
						continue
					}

					paramType, _ := propMap["type"].(string)
					paramDesc, _ := propMap["description"].(string)

					properties[paramName] = map[string]interface{}{
						"type":        paramType,
						"description": paramDesc,
					}
				}

				if reqList, ok := def.Parameters["required"].([]interface{}); ok {
					for _, r := range reqList {
						if reqStr, ok := r.(string); ok {
							required = append(required, reqStr)
						}
					}
				}
			} else {
				// Flat format: each key is a parameter with its own type/description/required
				for paramName, paramDef := range def.Parameters {
					paramMap, ok := paramDef.(map[string]interface{})
					if !ok {
						continue
					}

					paramType, _ := paramMap["type"].(string)
					paramDesc, _ := paramMap["description"].(string)
					isRequired, _ := paramMap["required"].(bool)

					properties[paramName] = map[string]interface{}{
						"type":        paramType,
						"description": paramDesc,
					}

					if isRequired {
						required = append(required, paramName)
					}
				}
			}
		}

		parameters := map[string]interface{}{
			"type":       "object",
			"properties": properties,
		}
		if len(required) > 0 {
			parameters["required"] = required
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

type ToolInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Toolset     string   `json:"toolset"`
	Enabled     bool     `json:"enabled"`
	Parameters  int      `json:"parameters_count"`
}

func (r *ToolRegistry) GetToolInfo() []ToolInfo {
	info := make([]ToolInfo, 0, len(r.defs))
	for name, def := range r.defs {
		info = append(info, ToolInfo{
			Name:        name,
			Description: def.Description,
			Version:     def.Version,
			Toolset:     def.Toolset,
			Enabled:     r.enabled[name],
			Parameters:  len(def.Parameters),
		})
	}
	return info
}

func (r *ToolRegistry) GetToolsetInfo() map[string]interface{} {
	result := make(map[string]interface{})
	for _, ts := range r.toolsets {
		tools := r.ListByToolset(ts)
		toolNames := make([]string, 0, len(tools))
		for _, t := range tools {
			toolNames = append(toolNames, t.Name)
		}
		result[ts] = map[string]interface{}{
			"count": len(tools),
			"tools": toolNames,
		}
	}
	return result
}

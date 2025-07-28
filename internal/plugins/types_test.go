package plugins

import (
	"testing"
)

func TestPluginInfo(t *testing.T) {
	plugin := &PluginInfo{
		Name:        "test-plugin",
		Version:     "1.0.0",
		Description: "A test plugin",
		Author:      "Test Author",
		Enabled:     true,
	}

	if plugin.Name != "test-plugin" {
		t.Errorf("Expected Name to be 'test-plugin', got %q", plugin.Name)
	}

	if plugin.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got %q", plugin.Version)
	}

	if !plugin.Enabled {
		t.Errorf("Expected Enabled to be true, got %v", plugin.Enabled)
	}
}

func TestPluginCommand(t *testing.T) {
	command := &PluginCommand{
		Name:        "test-command",
		Description: "A test command",
		Usage:       "test-command [options]",
		Category:    "testing",
	}

	if command.Name != "test-command" {
		t.Errorf("Expected Name to be 'test-command', got %q", command.Name)
	}

	if command.Category != "testing" {
		t.Errorf("Expected Category to be 'testing', got %q", command.Category)
	}
}

func TestPluginHook(t *testing.T) {
	hook := &PluginHook{
		Name:        "test-hook",
		Event:       "pre-commit",
		Description: "A test hook",
		Priority:    100,
		Enabled:     true,
	}

	if hook.Name != "test-hook" {
		t.Errorf("Expected Name to be 'test-hook', got %q", hook.Name)
	}

	if hook.Event != "pre-commit" {
		t.Errorf("Expected Event to be 'pre-commit', got %q", hook.Event)
	}

	if hook.Priority != 100 {
		t.Errorf("Expected Priority to be 100, got %d", hook.Priority)
	}
}

func TestPluginConfig(t *testing.T) {
	config := &PluginConfig{
		EnabledPlugins: []string{"plugin1", "plugin2"},
		PluginPaths:    []string{"/path/to/plugins", "/another/path"},
		AutoLoad:       true,
	}

	if len(config.EnabledPlugins) != 2 {
		t.Errorf("Expected 2 enabled plugins, got %d", len(config.EnabledPlugins))
	}

	if config.EnabledPlugins[0] != "plugin1" {
		t.Errorf("Expected first plugin to be 'plugin1', got %q", config.EnabledPlugins[0])
	}

	if !config.AutoLoad {
		t.Errorf("Expected AutoLoad to be true, got %v", config.AutoLoad)
	}
}

func TestPluginMetadata(t *testing.T) {
	metadata := &PluginMetadata{
		APIVersion:    "1.0",
		MinCLIVersion: "2.0.0",
		MaxCLIVersion: "3.0.0",
		Dependencies:  []string{"dep1", "dep2"},
		Permissions:   []string{"read", "write"},
	}

	if metadata.APIVersion != "1.0" {
		t.Errorf("Expected APIVersion to be '1.0', got %q", metadata.APIVersion)
	}

	if len(metadata.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(metadata.Dependencies))
	}

	if len(metadata.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(metadata.Permissions))
	}
}
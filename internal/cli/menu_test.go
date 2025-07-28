package cli

import (
	"testing"
)

func TestMenuItem(t *testing.T) {
	// Test creating a menu item
	item := MenuItem{
		Label:       "Test Item",
		Description: "A test menu item",
		Handler: func() error {
			return nil
		},
		Condition: func() bool {
			return true
		},
	}

	if item.Label != "Test Item" {
		t.Errorf("Expected Label to be 'Test Item', got %q", item.Label)
	}

	if item.Description != "A test menu item" {
		t.Errorf("Expected Description to be 'A test menu item', got %q", item.Description)
	}

	// Test handler execution
	if err := item.Handler(); err != nil {
		t.Errorf("Handler() error = %v", err)
	}

	// Test condition evaluation
	if !item.Condition() {
		t.Errorf("Expected Condition() to return true")
	}
}

func TestMenuItemWithFailingHandler(t *testing.T) {
	item := MenuItem{
		Label:       "Failing Item",
		Description: "A menu item that fails",
		Handler: func() error {
			return &testError{"handler failed"}
		},
		Condition: func() bool {
			return true
		},
	}

	// Test handler failure
	if err := item.Handler(); err == nil {
		t.Errorf("Expected handler to return error")
	}
}

func TestMenuItemWithFalseCondition(t *testing.T) {
	item := MenuItem{
		Label:       "Conditional Item",
		Description: "A conditionally visible item",
		Handler: func() error {
			return nil
		},
		Condition: func() bool {
			return false
		},
	}

	// Test condition evaluation
	if item.Condition() {
		t.Errorf("Expected Condition() to return false")
	}
}

func TestBuildMenu(t *testing.T) {
	// Simulate building a menu with various items
	items := []MenuItem{
		{
			Label:       "Always Visible",
			Description: "This item is always visible",
			Handler:     func() error { return nil },
			Condition:   func() bool { return true },
		},
		{
			Label:       "Never Visible", 
			Description: "This item is never visible",
			Handler:     func() error { return nil },
			Condition:   func() bool { return false },
		},
		{
			Label:       "Sometimes Visible",
			Description: "This item is conditionally visible",
			Handler:     func() error { return nil },
			Condition:   func() bool { return true }, // Simulate condition is met
		},
	}

	// Filter items based on conditions
	var visibleItems []MenuItem
	for _, item := range items {
		if item.Condition() {
			visibleItems = append(visibleItems, item)
		}
	}

	if len(visibleItems) != 2 {
		t.Errorf("Expected 2 visible items, got %d", len(visibleItems))
	}

	// Verify the correct items are visible
	expectedLabels := []string{"Always Visible", "Sometimes Visible"}
	for i, item := range visibleItems {
		if item.Label != expectedLabels[i] {
			t.Errorf("Expected item %d to be %q, got %q", i, expectedLabels[i], item.Label)
		}
	}
}

func TestMenuNavigation(t *testing.T) {
	// Test menu navigation logic
	tests := []struct {
		name         string
		menuSize     int
		selection    int
		expectedIdx  int
		expectError  bool
	}{
		{
			name:        "valid selection - first item",
			menuSize:    5,
			selection:   1,
			expectedIdx: 0,
			expectError: false,
		},
		{
			name:        "valid selection - last item",
			menuSize:    5,
			selection:   5,
			expectedIdx: 4,
			expectError: false,
		},
		{
			name:        "invalid selection - zero",
			menuSize:    5,
			selection:   0,
			expectedIdx: -1,
			expectError: true,
		},
		{
			name:        "invalid selection - too high",
			menuSize:    5,
			selection:   6,
			expectedIdx: -1,
			expectError: true,
		},
		{
			name:        "invalid selection - negative",
			menuSize:    5,
			selection:   -1,
			expectedIdx: -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate menu selection validation
			var idx int
			var err error

			if tt.selection < 1 || tt.selection > tt.menuSize {
				idx = -1
				err = &testError{"invalid selection"}
			} else {
				idx = tt.selection - 1 // Convert to 0-based index
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if idx != tt.expectedIdx {
				t.Errorf("Expected index %d, got %d", tt.expectedIdx, idx)
			}
		})
	}
}

func TestMenuFormatting(t *testing.T) {
	// Test menu display formatting
	items := []MenuItem{
		{
			Label:       "🌿 Branch Operations",
			Description: "Create, delete, and switch branches",
			Handler:     func() error { return nil },
			Condition:   func() bool { return true },
		},
		{
			Label:       "📝 Commit Changes",
			Description: "Stage and commit your changes",
			Handler:     func() error { return nil },
			Condition:   func() bool { return true },
		},
	}

	// Test that menu items have expected properties
	for i, item := range items {
		if item.Label == "" {
			t.Errorf("Item %d has empty label", i)
		}

		if item.Description == "" {
			t.Errorf("Item %d has empty description", i)
		}

		if item.Handler == nil {
			t.Errorf("Item %d has nil handler", i)
		}

		if item.Condition == nil {
			t.Errorf("Item %d has nil condition", i)
		}
	}
}

func TestMenuStateManagement(t *testing.T) {
	// Test menu state tracking
	type MenuState struct {
		CurrentMenu string
		History     []string
		CanGoBack   bool
	}

	state := MenuState{
		CurrentMenu: "main",
		History:     []string{},
		CanGoBack:   false,
	}

	// Test initial state
	if state.CurrentMenu != "main" {
		t.Errorf("Expected current menu to be 'main', got %q", state.CurrentMenu)
	}

	if state.CanGoBack {
		t.Errorf("Should not be able to go back from main menu")
	}

	// Navigate to submenu
	state.History = append(state.History, state.CurrentMenu)
	state.CurrentMenu = "branch"
	state.CanGoBack = len(state.History) > 0

	if state.CurrentMenu != "branch" {
		t.Errorf("Expected current menu to be 'branch', got %q", state.CurrentMenu)
	}

	if !state.CanGoBack {
		t.Errorf("Should be able to go back after navigating to submenu")
	}

	// Go back to previous menu
	if state.CanGoBack && len(state.History) > 0 {
		state.CurrentMenu = state.History[len(state.History)-1]
		state.History = state.History[:len(state.History)-1]
		state.CanGoBack = len(state.History) > 0
	}

	if state.CurrentMenu != "main" {
		t.Errorf("Expected to be back at main menu, got %q", state.CurrentMenu)
	}

	if state.CanGoBack {
		t.Errorf("Should not be able to go back from main menu after returning")
	}
}

func TestMenuItemExecution(t *testing.T) {
	// Test different types of menu item executions
	var executionLog []string

	items := []MenuItem{
		{
			Label: "Success Item",
			Handler: func() error {
				executionLog = append(executionLog, "success")
				return nil
			},
			Condition: func() bool { return true },
		},
		{
			Label: "Error Item",
			Handler: func() error {
				executionLog = append(executionLog, "error")
				return &testError{"simulated error"}
			},
			Condition: func() bool { return true },
		},
		{
			Label: "Disabled Item",
			Handler: func() error {
				executionLog = append(executionLog, "disabled")
				return nil
			},
			Condition: func() bool { return false },
		},
	}

	// Execute available items
	for _, item := range items {
		if item.Condition() {
			err := item.Handler()
			if err != nil {
				t.Logf("Item %q returned error: %v", item.Label, err)
			}
		}
	}

	// Check execution log
	expectedLog := []string{"success", "error"}
	if len(executionLog) != len(expectedLog) {
		t.Errorf("Expected %d executions, got %d", len(expectedLog), len(executionLog))
	}

	for i, expected := range expectedLog {
		if i < len(executionLog) && executionLog[i] != expected {
			t.Errorf("Expected execution %d to be %q, got %q", i, expected, executionLog[i])
		}
	}
}
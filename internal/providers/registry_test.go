package providers

import (
	"context"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Errorf("NewRegistry() returned nil")
	}

	if registry.providers == nil {
		t.Errorf("Registry providers map should be initialized")
	}
}

func TestRegisterProvider(t *testing.T) {
	registry := NewRegistry()
	
	// Create a mock provider
	mockProvider := &MockProvider{
		name: "test-provider",
	}

	err := registry.RegisterProvider("test", mockProvider)
	if err != nil {
		t.Errorf("RegisterProvider() error = %v", err)
	}

	// Test that provider is registered
	providers := registry.ListProviders()
	found := false
	for _, name := range providers {
		if name == "test" {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Provider 'test' was not found in registry")
	}
}

func TestRegisterDuplicateProvider(t *testing.T) {
	registry := NewRegistry()
	
	mockProvider1 := &MockProvider{name: "provider1"}
	mockProvider2 := &MockProvider{name: "provider2"}

	// Register first provider
	err := registry.RegisterProvider("duplicate", mockProvider1)
	if err != nil {
		t.Errorf("RegisterProvider() error = %v", err)
	}

	// Try to register duplicate
	err = registry.RegisterProvider("duplicate", mockProvider2)
	if err == nil {
		t.Errorf("Expected error when registering duplicate provider")
	}
}

func TestGetProvider(t *testing.T) {
	registry := NewRegistry()
	
	mockProvider := &MockProvider{name: "test-provider"}
	registry.RegisterProvider("test", mockProvider)

	// Test getting existing provider
	provider, err := registry.GetProvider("test")
	if err != nil {
		t.Errorf("GetProvider() error = %v", err)
	}

	if provider == nil {
		t.Errorf("GetProvider() returned nil provider")
	}

	// Test getting non-existent provider
	provider, err = registry.GetProvider("nonexistent")
	if err == nil {
		t.Errorf("Expected error when getting non-existent provider")
	}

	if provider != nil {
		t.Errorf("Expected nil provider for non-existent provider")
	}
}

func TestUnregisterProvider(t *testing.T) {
	registry := NewRegistry()
	
	mockProvider := &MockProvider{name: "test-provider"}
	registry.RegisterProvider("test", mockProvider)

	// Verify provider exists
	_, err := registry.GetProvider("test")
	if err != nil {
		t.Errorf("Provider should exist before unregistering")
	}

	// Unregister provider
	err = registry.UnregisterProvider("test")
	if err != nil {
		t.Errorf("UnregisterProvider() error = %v", err)
	}

	// Verify provider is gone
	_, err = registry.GetProvider("test")
	if err == nil {
		t.Errorf("Provider should not exist after unregistering")
	}
}

func TestListProviders(t *testing.T) {
	registry := NewRegistry()

	// Initially should be empty
	providers := registry.ListProviders()
	if len(providers) != 0 {
		t.Errorf("Expected 0 providers initially, got %d", len(providers))
	}

	// Add some providers
	mockProvider1 := &MockProvider{name: "provider1"}
	mockProvider2 := &MockProvider{name: "provider2"}
	
	registry.RegisterProvider("github", mockProvider1)
	registry.RegisterProvider("gitlab", mockProvider2)

	providers = registry.ListProviders()
	if len(providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(providers))
	}

	// Check that both providers are listed
	found1, found2 := false, false
	for _, name := range providers {
		if name == "github" {
			found1 = true
		}
		if name == "gitlab" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("Not all registered providers were found in list")
	}
}

func TestGetDefaultProvider(t *testing.T) {
	registry := NewRegistry()
	
	// Should return empty string when no default is set
	defaultProvider := registry.GetDefaultProvider()
	if defaultProvider != "" {
		t.Errorf("Expected empty default provider, got %q", defaultProvider)
	}

	// Register a provider and set as default
	mockProvider := &MockProvider{name: "test-provider"}
	registry.RegisterProvider("github", mockProvider)
	
	err := registry.SetDefaultProvider("github")
	if err != nil {
		t.Errorf("SetDefaultProvider() error = %v", err)
	}

	defaultProvider = registry.GetDefaultProvider()
	if defaultProvider != "github" {
		t.Errorf("Expected default provider 'github', got %q", defaultProvider)
	}
}

func TestSetDefaultProvider(t *testing.T) {
	registry := NewRegistry()

	// Test setting non-existent provider as default
	err := registry.SetDefaultProvider("nonexistent")
	if err == nil {
		t.Errorf("Expected error when setting non-existent provider as default")
	}

	// Register provider and set as default
	mockProvider := &MockProvider{name: "test-provider"}
	registry.RegisterProvider("test", mockProvider)

	err = registry.SetDefaultProvider("test")
	if err != nil {
		t.Errorf("SetDefaultProvider() error = %v", err)
	}

	if registry.GetDefaultProvider() != "test" {
		t.Errorf("Default provider was not set correctly")
	}
}

// MockProvider is a test implementation of Provider
type MockProvider struct {
	name string
}

func (m *MockProvider) GetName() string {
	return m.name
}

func (m *MockProvider) GetType() string {
	return "mock"
}

func (m *MockProvider) Configure(config interface{}) error {
	return nil
}

func (m *MockProvider) IsConfigured() bool {
	return true
}

func (m *MockProvider) Authenticate(ctx context.Context, token string) error {
	return nil
}

func (m *MockProvider) IsAuthenticated() bool {
	return true
}

func (m *MockProvider) GetRepositoryInfo(url string) (*RepositoryInfo, error) {
	return &RepositoryInfo{
		Name:        "test-repo",
		FullName:    "test/test-repo",
		Description: "Test repository",
		URL:         url,
	}, nil
}

func (m *MockProvider) CreatePullRequest(repoURL, title, body, head, base string) (*PullRequest, error) {
	return &PullRequest{
		Number:      1,
		Title:       title,
		Description: body,
		State:       "open",
		URL:         "https://example.com/pr/1",
	}, nil
}

func (m *MockProvider) ListIssues(repoURL string) ([]*Issue, error) {
	return []*Issue{
		{
			Number:      1,
			Title:       "Test Issue",
			Description: "Test issue body",
			State:       "open",
			URL:         "https://example.com/issue/1",
		},
	}, nil
}
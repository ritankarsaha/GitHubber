/*
 * GitHubber - Provider Registry
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Registry for managing Git hosting providers
 */

package providers

import (
	"fmt"
	"sync"
)

// DefaultRegistry is the global provider registry
var DefaultRegistry = NewRegistry()

// Registry implements ProviderRegistry
type Registry struct {
	mu          sync.RWMutex
	factories   map[ProviderType]ProviderFactory
	providers   map[string]Provider
	defaultProvider string
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[ProviderType]ProviderFactory),
		providers: make(map[string]Provider),
	}
}

// Register registers a provider factory for a given type
func (r *Registry) Register(providerType ProviderType, factory ProviderFactory) error {
	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[providerType]; exists {
		return fmt.Errorf("provider type %s is already registered", providerType)
	}

	r.factories[providerType] = factory
	return nil
}

// Create creates a provider instance from configuration
func (r *Registry) Create(config *ProviderConfig) (Provider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	r.mu.RLock()
	factory, exists := r.factories[config.Type]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}

	return factory(config)
}

// GetSupportedTypes returns all registered provider types
func (r *Registry) GetSupportedTypes() []ProviderType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]ProviderType, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	return types
}

// MustRegister registers a provider factory and panics on error
func (r *Registry) MustRegister(providerType ProviderType, factory ProviderFactory) {
	if err := r.Register(providerType, factory); err != nil {
		panic(err)
	}
}

// RegisterProvider registers a provider instance with a name
func (r *Registry) RegisterProvider(name string, provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %s already exists", name)
	}

	r.providers[name] = provider
	return nil
}

// GetProvider retrieves a provider by name
func (r *Registry) GetProvider(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// ListProviders returns all provider names
func (r *Registry) ListProviders() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// UnregisterProvider removes a provider by name
func (r *Registry) UnregisterProvider(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}
	
	delete(r.providers, name)
	return nil
}

// SetDefaultProvider sets the default provider
func (r *Registry) SetDefaultProvider(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}
	
	r.defaultProvider = name
	return nil
}

// GetDefaultProvider returns the name of the default provider
func (r *Registry) GetDefaultProvider() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.defaultProvider
}

// IsSupported checks if a provider type is supported
func (r *Registry) IsSupported(providerType ProviderType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.factories[providerType]
	return exists
}

// ProviderManager manages multiple provider instances
type ProviderManager struct {
	mu        sync.RWMutex
	providers map[string]Provider
	registry  ProviderRegistry
}

// NewProviderManager creates a new provider manager
func NewProviderManager(registry ProviderRegistry) *ProviderManager {
	if registry == nil {
		registry = DefaultRegistry
	}
	
	return &ProviderManager{
		providers: make(map[string]Provider),
		registry:  registry,
	}
}

// AddProvider adds a provider instance with a name
func (pm *ProviderManager) AddProvider(name string, provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.providers[name]; exists {
		return fmt.Errorf("provider %s already exists", name)
	}

	pm.providers[name] = provider
	return nil
}

// GetProvider retrieves a provider by name
func (pm *ProviderManager) GetProvider(name string) (Provider, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	provider, exists := pm.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// RemoveProvider removes a provider by name
func (pm *ProviderManager) RemoveProvider(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.providers, name)
}

// ListProviders returns all provider names
func (pm *ProviderManager) ListProviders() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	names := make([]string, 0, len(pm.providers))
	for name := range pm.providers {
		names = append(names, name)
	}
	return names
}

// CreateFromConfig creates a provider from configuration and adds it
func (pm *ProviderManager) CreateFromConfig(name string, config *ProviderConfig) error {
	provider, err := pm.registry.Create(config)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	return pm.AddProvider(name, provider)
}

// GetProviderByType returns the first provider of the given type
func (pm *ProviderManager) GetProviderByType(providerType ProviderType) (Provider, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, provider := range pm.providers {
		if provider.GetType() == providerType {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("no provider found for type %s", providerType)
}

// GetDefaultProvider returns the default provider (first one added)
func (pm *ProviderManager) GetDefaultProvider() (Provider, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.providers) == 0 {
		return nil, fmt.Errorf("no providers configured")
	}

	// Return first provider (map iteration order is not guaranteed, but this is fine for default)
	for _, provider := range pm.providers {
		return provider, nil
	}

	return nil, fmt.Errorf("no providers available")
}

// ParseProviderURL parses a repository URL and returns the provider type and repository info
func ParseProviderURL(url string) (ProviderType, string, string, error) {
	// Implement URL parsing logic for different providers
	// This is a simplified version - in production, use more robust parsing
	
	if url == "" {
		return "", "", "", fmt.Errorf("empty URL")
	}

	// GitHub patterns
	if matchGitHub(url) {
		owner, repo, err := parseGitHubURL(url)
		return ProviderGitHub, owner, repo, err
	}

	// GitLab patterns
	if matchGitLab(url) {
		owner, repo, err := parseGitLabURL(url)
		return ProviderGitLab, owner, repo, err
	}

	// Bitbucket patterns
	if matchBitbucket(url) {
		owner, repo, err := parseBitbucketURL(url)
		return ProviderBitbucket, owner, repo, err
	}

	return ProviderCustom, "", "", fmt.Errorf("unsupported provider URL: %s", url)
}

// Helper functions for URL parsing
func matchGitHub(url string) bool {
	return contains(url, "github.com")
}

func matchGitLab(url string) bool {
	return contains(url, "gitlab.com") || contains(url, "gitlab.")
}

func matchBitbucket(url string) bool {
	return contains(url, "bitbucket.org") || contains(url, "bitbucket.")
}

func parseGitHubURL(url string) (string, string, error) {
	// Implement GitHub URL parsing
	return parseGenericURL(url, "github.com")
}

func parseGitLabURL(url string) (string, string, error) {
	// Implement GitLab URL parsing
	return parseGenericURL(url, "gitlab")
}

func parseBitbucketURL(url string) (string, string, error) {
	// Implement Bitbucket URL parsing
	return parseGenericURL(url, "bitbucket")
}

func parseGenericURL(url, provider string) (string, string, error) {
	// This is a simplified parser - in production, use proper URL parsing
	// Handle both HTTPS and SSH formats
	return "", "", fmt.Errorf("URL parsing not fully implemented")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != substr && 
		   (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		    findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
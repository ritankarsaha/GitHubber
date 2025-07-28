/*
 * GitHubber - CI/CD Manager Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Central manager for CI/CD provider operations
 */

package ci

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DefaultManager is the global CI/CD manager
var DefaultManager = NewManager()

// Manager implements CIManager
type Manager struct {
	mu           sync.RWMutex
	providers    map[string]CIProvider
	factories    map[CIPlatform]CIFactory
	repositories map[string][]string // repo URL -> provider names
	metrics      map[string]*CIMetrics
	status       map[string]*CIStatus
}

// NewManager creates a new CI/CD manager
func NewManager() *Manager {
	manager := &Manager{
		providers:    make(map[string]CIProvider),
		factories:    make(map[CIPlatform]CIFactory),
		repositories: make(map[string][]string),
		metrics:      make(map[string]*CIMetrics),
		status:       make(map[string]*CIStatus),
	}

	// Register default factories
	manager.registerDefaultFactories()
	
	return manager
}

// registerDefaultFactories registers built-in CI/CD provider factories
func (m *Manager) registerDefaultFactories() {
	// These would be implemented in separate files
	// m.RegisterFactory(PlatformGitHubActions, NewGitHubActionsFactory())
	// m.RegisterFactory(PlatformGitLabCI, NewGitLabCIFactory())
	// m.RegisterFactory(PlatformJenkins, NewJenkinsFactory())
}

// RegisterFactory registers a CI/CD provider factory
func (m *Manager) RegisterFactory(platform CIPlatform, factory CIFactory) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[platform] = factory
}

// CreateProvider creates a new CI/CD provider
func (m *Manager) CreateProvider(name string, config *CIConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	factory, exists := m.factories[config.Platform]
	if !exists {
		return fmt.Errorf("unsupported CI platform: %s", config.Platform)
	}

	provider, err := factory(config)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	m.providers[name] = provider
	
	// Initialize status
	m.status[name] = &CIStatus{
		Provider:    config.Platform,
		Connected:   provider.IsConnected(),
		LastCheck:   time.Now(),
		ErrorCount:  0,
		Capabilities: m.getProviderCapabilities(provider),
	}

	return nil
}

// RegisterProvider registers a CI/CD provider
func (m *Manager) RegisterProvider(name string, provider CIProvider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	m.providers[name] = provider
	
	// Initialize status
	m.status[name] = &CIStatus{
		Provider:    provider.GetPlatform(),
		Connected:   provider.IsConnected(),
		LastCheck:   time.Now(),
		ErrorCount:  0,
		Capabilities: m.getProviderCapabilities(provider),
	}

	return nil
}

// GetProvider retrieves a CI/CD provider
func (m *Manager) GetProvider(name string) (CIProvider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// ListProviders returns all registered providers
func (m *Manager) ListProviders() map[string]CIProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	providers := make(map[string]CIProvider, len(m.providers))
	for name, provider := range m.providers {
		providers[name] = provider
	}

	return providers
}

// RegisterRepository associates a repository with CI/CD providers
func (m *Manager) RegisterRepository(repoURL string, providers []string) error {
	if repoURL == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate providers exist
	for _, providerName := range providers {
		if _, exists := m.providers[providerName]; !exists {
			return fmt.Errorf("provider %s not found", providerName)
		}
	}

	m.repositories[repoURL] = providers
	return nil
}

// GetRepositoryProviders returns CI/CD providers for a repository
func (m *Manager) GetRepositoryProviders(repoURL string) ([]CIProvider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	providerNames, exists := m.repositories[repoURL]
	if !exists {
		return nil, fmt.Errorf("repository %s not registered", repoURL)
	}

	providers := make([]CIProvider, 0, len(providerNames))
	for _, name := range providerNames {
		if provider, exists := m.providers[name]; exists {
			providers = append(providers, provider)
		}
	}

	return providers, nil
}

// TriggerBuilds triggers builds across all providers for a repository
func (m *Manager) TriggerBuilds(ctx context.Context, repoURL string, request *TriggerPipelineRequest) ([]*Pipeline, error) {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return nil, err
	}

	var allPipelines []*Pipeline
	var errors []string

	for _, provider := range providers {
		pipeline, err := provider.TriggerPipeline(ctx, repoURL, request)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", provider.GetName(), err))
			m.incrementErrorCount(provider.GetName())
			continue
		}
		allPipelines = append(allPipelines, pipeline)
	}

	if len(errors) > 0 && len(allPipelines) == 0 {
		return nil, fmt.Errorf("all providers failed: %v", errors)
	}

	return allPipelines, nil
}

// GetPipelineStatus gets pipeline status from the appropriate provider
func (m *Manager) GetPipelineStatus(ctx context.Context, repoURL, pipelineID string) (*Pipeline, error) {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return nil, err
	}

	// Try each provider until we find the pipeline
	for _, provider := range providers {
		pipeline, err := provider.GetPipeline(ctx, repoURL, pipelineID)
		if err == nil {
			return pipeline, nil
		}
	}

	return nil, fmt.Errorf("pipeline %s not found in any provider", pipelineID)
}

// CancelPipelines cancels pipelines across all providers
func (m *Manager) CancelPipelines(ctx context.Context, repoURL, pipelineID string) error {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return err
	}

	var errors []string
	var success bool

	for _, provider := range providers {
		err := provider.CancelPipeline(ctx, repoURL, pipelineID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", provider.GetName(), err))
			m.incrementErrorCount(provider.GetName())
			continue
		}
		success = true
	}

	if !success && len(errors) > 0 {
		return fmt.Errorf("failed to cancel pipeline: %v", errors)
	}

	return nil
}

// GetPipelinesByStatus returns pipelines with specific status
func (m *Manager) GetPipelinesByStatus(ctx context.Context, repoURL string, status BuildStatus, limit int) ([]*Pipeline, error) {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return nil, err
	}

	options := &ListPipelineOptions{
		Status: status,
		Limit:  limit,
	}

	var allPipelines []*Pipeline
	
	for _, provider := range providers {
		pipelines, err := provider.ListPipelines(ctx, repoURL, options)
		if err != nil {
			m.incrementErrorCount(provider.GetName())
			continue
		}
		allPipelines = append(allPipelines, pipelines...)
	}

	return allPipelines, nil
}

// GetBuildLogs retrieves build logs for a specific build
func (m *Manager) GetBuildLogs(ctx context.Context, repoURL, buildID string) (*BuildLogs, error) {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return nil, err
	}

	// Try each provider until we find the build
	for _, provider := range providers {
		logs, err := provider.GetBuildLogs(ctx, repoURL, buildID)
		if err == nil {
			return logs, nil
		}
	}

	return nil, fmt.Errorf("build %s not found in any provider", buildID)
}

// GetArtifacts retrieves artifacts for a pipeline
func (m *Manager) GetArtifacts(ctx context.Context, repoURL, pipelineID string) ([]*Artifact, error) {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return nil, err
	}

	var allArtifacts []*Artifact
	
	for _, provider := range providers {
		artifacts, err := provider.ListArtifacts(ctx, repoURL, pipelineID)
		if err != nil {
			continue // Try next provider
		}
		allArtifacts = append(allArtifacts, artifacts...)
	}

	return allArtifacts, nil
}

// GetEnvironments returns environments across all providers
func (m *Manager) GetEnvironments(ctx context.Context, repoURL string) ([]*Environment, error) {
	providers, err := m.GetRepositoryProviders(repoURL)
	if err != nil {
		return nil, err
	}

	var allEnvironments []*Environment
	
	for _, provider := range providers {
		environments, err := provider.ListEnvironments(ctx, repoURL)
		if err != nil {
			continue
		}
		allEnvironments = append(allEnvironments, environments...)
	}

	return allEnvironments, nil
}

// Health and monitoring methods

// GetProviderStatus returns the status of a provider
func (m *Manager) GetProviderStatus(name string) (*CIStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, exists := m.status[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return status, nil
}

// GetProviderMetrics returns metrics for a provider
func (m *Manager) GetProviderMetrics(name string) (*CIMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.metrics[name]
	if !exists {
		// Initialize metrics if not found
		metrics = &CIMetrics{
			CollectedAt: time.Now(),
		}
		m.metrics[name] = metrics
	}

	return metrics, nil
}

// UpdateProviderMetrics updates metrics for a provider
func (m *Manager) UpdateProviderMetrics(name string, metrics *CIMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[name] = metrics
}

// CheckProviderHealth performs health checks on all providers
func (m *Manager) CheckProviderHealth(ctx context.Context) map[string]*CIStatus {
	m.mu.RLock()
	providers := make(map[string]CIProvider, len(m.providers))
	for name, provider := range m.providers {
		providers[name] = provider
	}
	m.mu.RUnlock()

	results := make(map[string]*CIStatus)
	
	for name, provider := range providers {
		start := time.Now()
		connected := provider.IsConnected()
		responseTime := time.Since(start)

		m.mu.Lock()
		status := m.status[name]
		if status == nil {
			status = &CIStatus{
				Provider: provider.GetPlatform(),
			}
			m.status[name] = status
		}
		
		status.Connected = connected
		status.LastCheck = time.Now()
		status.ResponseTime = responseTime
		
		if !connected {
			status.ErrorCount++
		}
		
		results[name] = status
		m.mu.Unlock()
	}

	return results
}

// CollectMetrics collects metrics from all providers
func (m *Manager) CollectMetrics(ctx context.Context) map[string]*CIMetrics {
	// This would typically collect metrics from providers
	// For now, return cached metrics
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]*CIMetrics, len(m.metrics))
	for name, metric := range m.metrics {
		metrics[name] = metric
	}

	return metrics
}

// Configuration management

// SaveConfiguration saves CI/CD configuration
func (m *Manager) SaveConfiguration(path string) error {
	// Implementation would save provider configurations to file
	return nil
}

// LoadConfiguration loads CI/CD configuration
func (m *Manager) LoadConfiguration(path string) error {
	// Implementation would load provider configurations from file
	return nil
}

// Event handling

// EmitEvent emits a CI/CD event
func (m *Manager) EmitEvent(event *CIEvent) error {
	// Implementation would handle event distribution
	// This could integrate with webhooks, notifications, etc.
	return nil
}

// Template management

// GetTemplates returns available pipeline templates
func (m *Manager) GetTemplates(ctx context.Context, platform CIPlatform) ([]*PipelineTemplate, error) {
	var templates []*PipelineTemplate
	
	for _, provider := range m.providers {
		if provider.GetPlatform() != platform {
			continue
		}
		
		// This would fetch templates from the provider
		// For now, return empty list
	}
	
	return templates, nil
}

// ValidateConfiguration validates pipeline configuration
func (m *Manager) ValidateConfiguration(ctx context.Context, platform CIPlatform, config []byte) (*ValidationResult, error) {
	for _, provider := range m.providers {
		if provider.GetPlatform() != platform {
			continue
		}
		
		return provider.ValidatePipelineConfig(ctx, config)
	}
	
	return nil, fmt.Errorf("no provider found for platform %s", platform)
}

// Helper methods

func (m *Manager) incrementErrorCount(providerName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if status, exists := m.status[providerName]; exists {
		status.ErrorCount++
	}
}

func (m *Manager) getProviderCapabilities(provider CIProvider) []string {
	capabilities := []string{"pipelines", "builds"}
	
	// This would inspect the provider to determine capabilities
	// For now, return basic capabilities
	return capabilities
}

// Shutdown gracefully shuts down all providers
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []string
	
	for name, provider := range m.providers {
		// If provider implements Shutdown method, call it
		if shutdowner, ok := provider.(interface{ Shutdown(context.Context) error }); ok {
			if err := shutdowner.Shutdown(ctx); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", name, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("shutdown errors: %v", errors)
	}

	return nil
}

// RemoveProvider removes a provider
func (m *Manager) RemoveProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}

	delete(m.providers, name)
	delete(m.status, name)
	delete(m.metrics, name)

	// Remove from repositories
	for repoURL, providers := range m.repositories {
		var filtered []string
		for _, providerName := range providers {
			if providerName != name {
				filtered = append(filtered, providerName)
			}
		}
		if len(filtered) == 0 {
			delete(m.repositories, repoURL)
		} else {
			m.repositories[repoURL] = filtered
		}
	}

	return nil
}

// GetProvidersByPlatform returns providers for a specific platform
func (m *Manager) GetProvidersByPlatform(platform CIPlatform) []CIProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var providers []CIProvider
	for _, provider := range m.providers {
		if provider.GetPlatform() == platform {
			providers = append(providers, provider)
		}
	}

	return providers
}
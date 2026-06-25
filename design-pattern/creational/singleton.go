package creational

import (
	"fmt"
	"sync"
)

// Config holds application configuration
type Config struct {
	AppName    string
	Version    string
	MaxRetries int
}

// ConfigManager is a singleton that manages configuration
type ConfigManager struct {
	config *Config
	mu     sync.RWMutex
}

var (
	instance *ConfigManager
	once     sync.Once
)

// GetConfigManager returns the singleton instance of ConfigManager
func GetConfigManager() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{
			config: &Config{
				AppName:    "MyApp",
				Version:    "1.0.0",
				MaxRetries: 3,
			},
		}
	})
	return instance
}

func (c *ConfigManager) GetConfig() *Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

func (c *ConfigManager) UpdateConfig(updates map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if appName, ok := updates["app_name"].(string); ok {
		c.config.AppName = appName
	}
	if version, ok := updates["version"].(string); ok {
		c.config.Version = version
	}
	if maxRetries, ok := updates["max_retries"].(int); ok {
		c.config.MaxRetries = maxRetries
	}
}

func (c *ConfigManager) GetAppName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.AppName
}

// SingletonExampleUsage demonstrates the Singleton pattern
func SingletonExampleUsage() {
	config1 := GetConfigManager()
	config2 := GetConfigManager()

	fmt.Printf("Same instance? %v\n", config1 == config2)
	fmt.Printf("AppName: %s\n", config1.GetAppName())

	config1.UpdateConfig(map[string]interface{}{
		"app_name": "MyApp-v2",
	})

	fmt.Printf("Updated AppName via config2: %s\n", config2.GetAppName())
}

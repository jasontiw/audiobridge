package config

import (
	"errors"
	"fmt"
)

// Config stub - full implementation in Phase 3
var (
	ErrNotImplemented = errors.New("config loading not implemented yet - coming in Phase 3")
)

// Config holds the application configuration
type Config struct {
	General   GeneralConfig
	Send      SendConfig
	Receive   ReceiveConfig
	Discovery DiscoveryConfig
}

// GeneralConfig holds general settings
type GeneralConfig struct {
	LogLevel string `toml:"log_level"`
	Port     int    `toml:"port"`
}

// SendConfig holds settings for sending audio
type SendConfig struct {
	Source     string `toml:"source"`
	Target     string `toml:"target"`
	Codec      string `toml:"codec"`
	Bitrate    int    `toml:"bitrate"`
	Channels   int    `toml:"channels"`
	SampleRate int    `toml:"sample_rate"`
}

// ReceiveConfig holds settings for receiving audio
type ReceiveConfig struct {
	OutputDevice   string `toml:"output_device"`
	JitterBufferMs int    `toml:"jitter_buffer_ms"`
}

// DiscoveryConfig holds settings for mDNS discovery
type DiscoveryConfig struct {
	Enabled      bool   `toml:"enabled"`
	AnnounceName string `toml:"announce_name"`
}

// Load loads configuration from a TOML file
func Load(path string) (*Config, error) {
	// TODO: Implement in Phase 3 with BurntSushi/toml
	return nil, ErrNotImplemented
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// TODO: Implement in Phase 3
	if c.General.Port < 1024 || c.General.Port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535, got %d", c.General.Port)
	}
	return nil
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		General: GeneralConfig{
			LogLevel: "info",
			Port:     9876,
		},
		Send: SendConfig{
			Source:     "microphone",
			Codec:      "opus",
			Bitrate:    64000,
			Channels:   1,
			SampleRate: 44100,
		},
		Receive: ReceiveConfig{
			JitterBufferMs: 50,
		},
		Discovery: DiscoveryConfig{
			Enabled: true,
		},
	}
}

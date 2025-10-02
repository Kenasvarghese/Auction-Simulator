package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration values for the auction simulator.
type Config struct {
	NumBidders       int `default:"100" split_words:"true"`
	NumAuctions      int `default:"10" split_words:"true"`
	AuctionTimeoutMs int `default:"100" split_words:"true"`
	AuctionVCPU      int `default:"1" split_words:"true"`
	AuctionMemory    int `default:"10" split_words:"true"`
	VCPU             int `required:"true" split_words:"true"`
	Memory           int `required:"true" split_words:"true"`
}

// LoadConfig loads configuration from environment variables into a Config struct.
func LoadConfig() *Config {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Println("error loading config:", err)
	}

	// validate only critical values
	if err := cfg.Validate(); err != nil {
		fmt.Println("invalid config:", err)
	}

	return &cfg
}

// Validate ensures required values are sane.
// Only the truly critical values are checked.
func (c *Config) Validate() error {
	if c.NumBidders <= 0 {
		return fmt.Errorf("number of bidders must be > 0")
	}
	if c.NumAuctions <= 0 {
		return fmt.Errorf("number of auctions must be > 0")
	}
	if c.AuctionTimeoutMs <= 0 {
		return fmt.Errorf("auction timeout must be > 0")
	}
	if c.VCPU <= 0 || c.VCPU < c.AuctionVCPU {
		return fmt.Errorf("VCPU must be > 0 and >= auction VCPU")
	}
	if c.Memory <= 0 || c.Memory < c.AuctionMemory {
		return fmt.Errorf("memory must be > 0 and >= auction Memory")
	}
	return nil
}

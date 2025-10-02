package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration values for the auction simulator.
type Config struct {
	NumBidders       int `required:"true" split_words:"true"`
	NumAttributes    int `required:"true" split_words:"true"`
	NumAuctions      int `required:"true" split_words:"true"`
	AuctionTimeoutMs int `required:"true" split_words:"true"`
	AuctionVCPU      int `required:"true" split_words:"true"`
	AuctionMemory    int `required:"true" split_words:"true"`
	VCPU             int `required:"true" split_words:"true"`
	Memory           int `required:"true" split_words:"true"`
}

// LoadConfig loads configuration from environment variables into a Config struct.
func LoadConfig() *Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		fmt.Println(err)
	}

	return &cfg
}

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
	if c.AuctionVCPU <= 0 {
		return fmt.Errorf("number of cpu per auction must be > 0")
	}
	if c.AuctionMemory <= 0 {
		return fmt.Errorf("memory per auction must be > 0")
	}
	if c.VCPU <= 0 {
		return fmt.Errorf("vCPU must be > 0")
	}
	if c.Memory <= 0 {
		return fmt.Errorf("memory must be > 0")
	}
	return nil
}

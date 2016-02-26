package db

import (
	//	"fmt"
//	"os"
	"errors"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megdc/subd"
)
const (
	SETTINGS = "settings"
)
type Config struct {
	Common    *subd.Config     `toml:"common"`
}

func (c Config) String() string {
	return ("\n" +
		c.Common.String() + "\n")
}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Config {
	c := &Config{}
	c.Common = subd.NewConfig()
	return c
}

// NewDemoConfig returns the config that runs when no config is specified.
func NewDemoConfig() (*Config, error) {
	c := NewConfig()
	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.Common.Home == "" {
		return errors.New("Home Dir must be specified")
	}
	return nil
}

func ParseConfig() (*Config, error) {

	config := NewConfig()
	path := config.Common.Home
	if path == "" {
		path = config.Common.Home + "/megdc.conf"
	}
	log.Warnf("Using configuration at: %s", path)
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	log.Debug(config)
	return config, nil
}

func StoreDB(data interface{}) error{
		t := TableInfo{
			Name: SETTINGS,
			Pks: []string{"name"},
			Ccms: []string{},
		}
 err := Write(t ,data)
 return err
}

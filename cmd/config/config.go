package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Cluster struct {
	Name        string
	Addrs       []string `yaml:"addrs"`
	Password    string   `yaml:"password"`
	Description string   `yaml:"desc"`
}

type Config struct {
	CurrentCluster  string `yaml:"current-cluster"`
	ClusterOverride string
	Clusters        []*Cluster `yaml:"clusters"`
}

func (c *Config) SetCurrentCluster(name string) error {
	var oldCluster string
	if c.ActiveCluster() != nil {
		oldCluster = c.ActiveCluster().Name
	}
	for _, cluster := range c.Clusters {
		if cluster.Name == name {
			c.CurrentCluster = name

			if err := c.Write(); err != nil {
				// "Revert" change to the cluster struct, either
				// everything is successful or nothing.
				c.CurrentCluster = oldCluster
				return err

			}
			return nil

		}
	}
	return fmt.Errorf("Could not find cluster with name %v", name)
}

func (c *Config) ActiveCluster() *Cluster {
	if c == nil {
		return nil
	}

	toSearch := c.ClusterOverride
	if c.ClusterOverride == "" {
		toSearch = c.CurrentCluster
	}

	if toSearch == "" {
		return nil
	}

	for _, cluster := range c.Clusters {
		if cluster.Name == toSearch {
			// Make copy of cluster struct, using a pointer leads to unintended
			// behavior where modifications on currentCluster are written back
			// into the config
			c := *cluster
			return &c
		}
	}
	return nil
}

func (c *Config) Write() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".godis")
	_ = os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config")

	file, err := os.OpenFile(configPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(&c)
}

func ReadConfig(cfgPath string) (c Config, err error) {
	file, err := os.OpenFile(getConfigPath(cfgPath), os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getConfigPath(cfgPath string) string {
	if !fileExists(cfgPath) {
		return getDefaultConfigPath()
	}
	return cfgPath
}

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".godis", "config")
}

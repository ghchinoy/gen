package model

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config is the configuration for the application.
type Config struct {
	ProjectID      string
	RegionID       string
	ConfigFile     string
	LogType        string
	OutputType     string
	ModelParameters map[string]interface{}
}

// ConfigBuilder is a builder for the Config struct.
type ConfigBuilder struct {
	projectID      string
	regionID       string
	configFile     string
	logType        string
	outputType     string
	modelParameters map[string]interface{}
}

// ProjectID sets the project ID.
func (b *ConfigBuilder) ProjectID(p string) *ConfigBuilder {
	b.projectID = p
	return b
}

// RegionID sets the region ID.
// Allowed values are: us-central1, us-east1, us-west1, etc.
func (b *ConfigBuilder) RegionID(r string) *ConfigBuilder {
	b.regionID = r
	return b
}

// ConfigFile sets the config file.
func (b *ConfigBuilder) ConfigFile(configFile string) *ConfigBuilder {
	b.configFile = configFile
	return b
}

// LogType sets the log type.
// Allowed values are: none, quiet, verbose.
func (b *ConfigBuilder) LogType(logType string) *ConfigBuilder {
	if logType != "none" && logType != "quiet" && logType != "verbose" {
		logType = "none"
	}
	b.logType = logType
	return b
}

// OutputType sets the output type.
// Allowed values are: text, json.
func (b *ConfigBuilder) OutputType(outputType string) *ConfigBuilder {
	if outputType != "text" && outputType != "json" {
		outputType = "text"
	}
	b.outputType = outputType
	return b
}

// Describe returns a string description of the ConfigBuilder.
func (b *ConfigBuilder) Describe() string {
	return fmt.Sprintf("%+v", b)
}

// Build builds the Config struct.
func (b *ConfigBuilder) Build() (Config, error) {

	cfg := Config{}

	if b.projectID == "" {
		return cfg, fmt.Errorf("need a valid GCP project ID")
	}
	cfg.ProjectID = b.projectID

	if b.regionID == "" {
		cfg.RegionID = "us-central1"
	} else {
		cfg.RegionID = b.regionID
	}

	cfg.ConfigFile = b.configFile
	cfg.LogType = b.logType
	cfg.OutputType = b.outputType

	if b.configFile != "" {
		data, err := os.ReadFile(b.configFile)
		if err != nil {
			return cfg, fmt.Errorf("error reading model config: %v", err)
		}

		err = json.Unmarshal(data, &b.modelParameters)
		if err != nil {
			return cfg, fmt.Errorf("error unmarshalling model config: %v", err)
		}
		cfg.ModelParameters = b.modelParameters
	}

	return cfg, nil
}

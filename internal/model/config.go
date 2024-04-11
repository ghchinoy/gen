package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	ProjectID  string
	RegionID   string
	ConfigFile string
	LogType    string
	OutputType string
}

type ConfigBuilder struct {
	projectID  string
	regionID   string
	configFile string
	logType    string
	outputType string
}

func (b *ConfigBuilder) ProjectID(p string) *ConfigBuilder {
	b.projectID = p
	return b
}

func (b *ConfigBuilder) RegionID(r string) *ConfigBuilder {
	b.regionID = r
	return b
}
func (b *ConfigBuilder) ConfigFile(configFile string) *ConfigBuilder {
	b.configFile = configFile
	return b
}

func (b *ConfigBuilder) LogType(logType string) *ConfigBuilder {
	b.logType = logType
	return b
}

func (b *ConfigBuilder) OutputType(outputType string) *ConfigBuilder {
	b.outputType = outputType
	return b
}

func (b *ConfigBuilder) Describe() string {
	return fmt.Sprintf("%+v", b)
}

func (b *ConfigBuilder) Build() (Config, error) {

	cfg := Config{}

	if b.projectID == "" {
		log.Fatalln("Need a valid GCP project ID")
	} else {
		cfg.ProjectID = b.projectID
	}

	if b.regionID == "" {
		cfg.RegionID = "us-central1"
	} else {
		cfg.RegionID = b.regionID
	}

	if b.configFile == "" {
		cfg.ConfigFile = b.configFile
	}

	if b.logType == "" {
		cfg.LogType = "none"
	} else {
		cfg.LogType = b.logType
	}

	if b.outputType == "" {
		cfg.OutputType = b.outputType
	}

	return cfg, nil
}

// TODO - Revisit ReadModelConfigFile() and decide whether it should be exported or not
// and whether it should instead set the fields of the struct as a way
// to initialize a model configuration from a file.  If thats the case, then this
// should be a method on the ConfigBuilder, that way when the Build() method is invoked
// it would validate the inputs and return a valid Config instance.

// readModelConfigFile reads the model configuration file (JSON text file)
func (cfg Config) ReadModelConfigFile() (map[string]interface{}, error) {

	var config map[string]interface{}
	data, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		return config, fmt.Errorf("error reading model config: %v", err)

	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling model config: %v", err)
	}
	return config, nil
}

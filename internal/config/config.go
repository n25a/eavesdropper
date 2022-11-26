package config

import (
	"os"

	"go.uber.org/zap"

	"gopkg.in/yaml.v3"

	"github.com/n25a/eavesdropper/internal/database"
	"github.com/n25a/eavesdropper/internal/log"
	"github.com/n25a/eavesdropper/internal/mq"
)

// C - global config
var C *Config

// Config - config struct
type Config struct {
	MQ         MQ       `yaml:"mq"`
	Database   Database `yaml:"database"`
	SchemaPath string   `yaml:"schema_path"`
}

// MQ - message queue config
type MQ struct {
	Type mq.MQType `yaml:"type"`
	Conf mq.Conf   `yaml:"conf"`
}

// Database - database config
type Database struct {
	Type          database.DatabaseType `yaml:"type"`
	MigrationPath string                `yaml:"migration_path"`
	Conf          database.Conf         `yaml:"conf"`
}

// LoadConfig - load config from file
func LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Logger.Warn("error in closing schema file", zap.Error(err))
		}
	}(file)

	d := yaml.NewDecoder(file)
	if err := d.Decode(&C); err != nil {
		return err
	}

	return nil
}

package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/n25a/eavesdropper/internal/database"
	"github.com/n25a/eavesdropper/internal/mq"
)

var C *Config

type Config struct {
	Log      Log      `yaml:"log"`
	MQ       MQ       `yaml:"mq"`
	Database Database `yaml:"database"`
}

type Log struct {
	Level string `yaml:"level"`
}

type MQ struct {
	Type mq.MQType `yaml:"type"`
	Conf mq.MQConf `yaml:"conf"`
}

type Database struct {
	Type database.DatabaseType `yaml:"type"`
	Conf database.DatabaseConf `yaml:"conf"`
}

func LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	d := yaml.NewDecoder(file)
	if err := d.Decode(&C); err != nil {
		return err
	}

	return nil
}

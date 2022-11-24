package app

import (
	"encoding/json"
	"os"

	"github.com/n25a/eavesdropper/internal/config"
	"github.com/n25a/eavesdropper/internal/database"
	"github.com/n25a/eavesdropper/internal/mq"
	"gopkg.in/yaml.v3"
)

var A *App

type App struct {
	DB      database.Database
	MQ      mq.MessageQueue
	Schemas map[string]Schema
}

type Schema struct {
	Fields []string
	Query  string
}

type schemaBinder struct {
	Subject string `yaml:"subject"`
	Data    string `yaml:"data"`
}

func InitApp() error {
	A = &App{
		DB:      database.NewDatabase(config.C.Database.Type),
		MQ:      mq.NewMessageQueue(mq.NATS),
		Schemas: map[string]Schema{},
	}

	// Load schema from file
	file, err := os.Open(config.C.SchemaPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var binder []schemaBinder
	d := yaml.NewDecoder(file)
	if err := d.Decode(&binder); err != nil {
		return err
	}

	// Parse schema
	for _, b := range binder {
		// Parse fields
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(b.Data), &data); err != nil {
			return err
		}
		var fields []string
		for k, _ := range data {
			fields = append(fields, k)
		}

		// TODO: Parse query (after add query builder)

		// Add schema
		A.Schemas[b.Subject] = Schema{
			Fields: fields,
		}
	}

	return nil
}

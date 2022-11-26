package app

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"

	"github.com/n25a/eavesdropper/internal/config"
	"github.com/n25a/eavesdropper/internal/database"
	"github.com/n25a/eavesdropper/internal/log"
	"github.com/n25a/eavesdropper/internal/mq"
	"gopkg.in/yaml.v3"
)

// A - application
var A *App

// App - main app struct
type App struct {
	DB      database.Database
	MQ      mq.MessageQueue
	Schemas map[string]Schema
}

// Schema - schema struct for each subject on message queue
type Schema struct {
	Fields []string
	Query  string
}

type schemaBinder struct {
	Subject string `yaml:"subject"`
	Table   string `yaml:"table"`
	Data    string `yaml:"data"`
}

// InitApp - initialize app
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
			log.Logger.Warn("error in closing schema file", zap.Error(err))
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

		query := A.DB.BuildInsertQuery(b.Table, fields)

		// Add schema
		A.Schemas[b.Subject] = Schema{
			Fields: fields,
			Query:  query,
		}
	}

	return nil
}

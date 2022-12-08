package app

import (
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
	Schemas map[string][]Schema
}

// Schema - schema struct for each subject on message queue
type Schema struct {
	Fields []string
	Query  string
}

type storage struct {
	Table           string            `yaml:"table"`
	fieldToDBColumn map[string]string `yaml:"field_to_db_column"`
}

type schemaBinder struct {
	Subject string    `yaml:"subject"`
	Storage []storage `yaml:"storage"`
}

// InitApp - initialize app
func InitApp() error {
	A = &App{
		DB:      database.NewDatabase(config.C.Database.Type),
		MQ:      mq.NewMessageQueue(mq.NATS),
		Schemas: map[string][]Schema{},
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
		var schemas []Schema
		for _, s := range b.Storage {
			var dbColumns []string
			var fields []string
			for field, dbColumn := range s.fieldToDBColumn {
				dbColumns = append(dbColumns, dbColumn)
				fields = append(fields, field)
			}

			query := A.DB.BuildInsertQuery(s.Table, dbColumns)
			schemas = append(schemas, Schema{
				Fields: fields,
				Query:  query,
			})
		}

		// Add schema
		A.Schemas[b.Subject] = schemas
	}

	return nil
}

package cmd

import (
	"go.uber.org/zap"

	"github.com/n25a/eavesdropper/internal/app"
	"github.com/n25a/eavesdropper/internal/config"
	"github.com/n25a/eavesdropper/internal/database"
	"github.com/n25a/eavesdropper/internal/log"
	"github.com/spf13/cobra"
)

var migrationCMD = &cobra.Command{
	Use:   "migration",
	Short: "Migrate Migration files to database",
	Run: func(cmd *cobra.Command, args []string) {
		Migration()
	},
}

func Migration() {
	if err := config.LoadConfig(ConfigPath); err != nil {
		log.Logger.Panic("error in loading config", zap.Error(err))
	}

	if err := app.InitApp(); err != nil {
		log.Logger.Panic("error in initializing app", zap.Error(err))
	}

	if err := app.A.DB.Connect(); err != nil {
		log.Logger.Panic("error in connecting to database", zap.Error(err))
	}

	defer func(DB database.Database) {
		err := DB.Close()
		if err != nil {
			log.Logger.Warn("error in closing connection from database", zap.Error(err))
		}
	}(app.A.DB)

	err := app.A.DB.Migrate()
	if err != nil {
		log.Logger.Panic("error in migrating", zap.Error(err))
	}
	log.Logger.Info("migration completed")
}

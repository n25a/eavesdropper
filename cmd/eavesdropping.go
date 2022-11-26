package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/n25a/eavesdropper/internal/app"
	"github.com/n25a/eavesdropper/internal/config"
	"github.com/n25a/eavesdropper/internal/database"
	"github.com/n25a/eavesdropper/internal/log"
	"github.com/n25a/eavesdropper/internal/mq"

	"github.com/spf13/cobra"
)

var configPath string

var eavesdroppingCMD = &cobra.Command{
	Use:     "eavesdropping",
	Aliases: []string{"ed"},
	Short:   "Consuming messages and store it in db",
	Run: func(cmd *cobra.Command, args []string) {
		Eavesdropping()
	},
}

func init() {
	eavesdroppingCMD.Flags().StringVarP(&configPath, "config", "c", "", "config file")
}

// Eavesdropping - consuming messages and store it in db
func Eavesdropping() {
	if err := config.LoadConfig(configPath); err != nil {
		log.Logger.Panic("error in loading config", zap.Error(err))
	}

	if err := app.InitApp(); err != nil {
		log.Logger.Panic("error in initializing app", zap.Error(err))
	}

	if err := app.A.MQ.Connect(); err != nil {
		log.Logger.Panic("error in connecting to message queue", zap.Error(err))
	}

	defer func(MQ mq.MessageQueue) {
		err := MQ.UnSubscribe()
		if err != nil {
			log.Logger.Warn("error in Unsubscribing", zap.Error(err))
		}
	}(app.A.MQ)

	if err := app.A.DB.Connect(); err != nil {
		log.Logger.Panic("error in connecting to database", zap.Error(err))
	}

	defer func(DB database.Database) {
		err := DB.Close()
		if err != nil {
			log.Logger.Warn("error in closing connection from database", zap.Error(err))
		}
	}(app.A.DB)

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for subject, _ := range app.A.Schemas {
		err := app.A.MQ.Subscribe(subject, app.A.DB.Insert)
		if err != nil {
			log.Logger.Panic("error in subscribing", zap.Error(err), zap.String("subject", subject))
		}
	}

	<-shutdown
}

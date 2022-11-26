package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/n25a/eavesdropper/internal/database"

	"github.com/n25a/eavesdropper/internal/mq"

	"github.com/n25a/eavesdropper/internal/app"

	"github.com/n25a/eavesdropper/internal/config"

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

func Eavesdropping() {
	if err := config.LoadConfig(configPath); err != nil {
		panic(err)
	}

	if err := app.InitApp(); err != nil {
		panic(err)
	}

	if err := app.A.MQ.Connect(); err != nil {
		panic(err)
	}

	defer func(MQ mq.MessageQueue) {
		err := MQ.UnSubscribe()
		if err != nil {
			log.Println(err)
		}
	}(app.A.MQ)

	if err := app.A.DB.Connect(); err != nil {
		panic(err)
	}

	defer func(DB database.Database) {
		err := DB.Close()
		if err != nil {

		}
	}(app.A.DB)

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for subject, _ := range app.A.Schemas {
		err := app.A.MQ.Subscribe(subject, app.A.DB.Insert)
		if err != nil {
			panic(err)
		}
	}

	<-shutdown
}

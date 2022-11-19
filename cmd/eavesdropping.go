package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var eavesdroppingCMD = &cobra.Command{
	Use:   "eavesdropping",
	Short: "Consuming messages and store it in db",
	Run: func(cmd *cobra.Command, args []string) {
		eavesdropping()
	},
}

func eavesdropping() {

	// consume

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	// shutdown service

	//l.Logger.Info("Shutting down the application")
	//l.Teardown()
}

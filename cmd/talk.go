package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var talkCMD = &cobra.Command{
	Use:   "talk",
	Short: "publish fake messages",
	Run: func(cmd *cobra.Command, args []string) {
		talk()
	},
}

func talk() {

	// send messages

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	// shutdown service

	//l.Logger.Info("Shutting down the application")
	//l.Teardown()
}

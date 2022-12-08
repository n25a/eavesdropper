package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "eavesdropper",
	Short: "CLI for subscribe messages and store them",
	Long:  `Eavesdropper is a simple CLI to receives messages from different message queues and store them to a database`,
}

var (
	ConfigPath string

	asciiArt = `
███████╗ █████╗ ██╗   ██╗███████╗███████╗██████╗ ██████╗  ██████╗ ██████╗ ██████╗ ███████╗██████╗ 
██╔════╝██╔══██╗██║   ██║██╔════╝██╔════╝██╔══██╗██╔══██╗██╔═══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗
█████╗  ███████║██║   ██║█████╗  ███████╗██║  ██║██████╔╝██║   ██║██████╔╝██████╔╝█████╗  ██████╔╝
██╔══╝  ██╔══██║╚██╗ ██╔╝██╔══╝  ╚════██║██║  ██║██╔══██╗██║   ██║██╔═══╝ ██╔═══╝ ██╔══╝  ██╔══██╗
███████╗██║  ██║ ╚████╔╝ ███████╗███████║██████╔╝██║  ██║╚██████╔╝██║     ██║     ███████╗██║  ██║
╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚══════╝╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝     ╚══════╝╚═╝  ╚═╝
`
)

func init() {
	fmt.Println(asciiArt)
	rootCMD.AddCommand(eavesdroppingCMD, migrationCMD)
	rootCMD.Flags().StringVarP(&ConfigPath, "config", "c", "", "config file")
}

// Execute executes the root command.
func Execute() error {
	return rootCMD.Execute()
}

package main

import (
	"github.com/milligan22963/circus/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "Dark Circus",
		Short: "An escape room application",
		Long:  `An application to interface with the local device and present a web page for controlling`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			print("starting up\n")
		},
	}

	rootCmd.AddCommand(cmd.ServerCmd)
	cmd.ServerCmd.Flags().String("config", "settings.yaml", "configuration file")

	rootCmd.AddCommand(cmd.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		println("failed to initialize: ", err.Error())
	}
}

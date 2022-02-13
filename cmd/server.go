// Package cmd is for all commands in the cmd tree
package cmd

import (
	"github.com/milligan22963/circus/config"
	"github.com/milligan22963/circus/pkg/server"
	"github.com/spf13/cobra"
)

// ServerCmd is the main server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server hosts the web site",
	Long:  `A base server interface running on a device`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			panic("Unable to find config flag")
		}

		appConfig := config.NewSiteConfiguration(configFile)

		serverInstance := server.ServerInstance{}

		serverInstance.Run(appConfig)
	},
}

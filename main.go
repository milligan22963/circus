package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	local "pkg/management"

	"github.com/milligan22963/cmra/cmd/subcmd"
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

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

	rootCmd.AddCommand(subcmd.ServerCmd)
	subcmd.ServerCmd.Flags().String("config", "settings.yaml", "configuration file")

	rootCmd.AddCommand(subcmd.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		println("failed to initialize: ", err.Error())
	}

	err := adapter.Enable()

	if err != nil {
		fmt.Printf("failed to enable adapter: %+v", err)
		return
	}

	targetName := "Dark Circus"
	ch := make(chan bluetooth.ScanResult, 1)

	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		fmt.Printf("%+v, payload: %+v", result, result.AdvertisementPayload)
		//		println("found device:", result.Address.String(), result.RSSI, result.LocalName())
		if result.LocalName() == targetName {
			adapter.StopScan()
			ch <- result
		}
	})

	if err != nil {
		fmt.Printf("failed to enable scan: %+v", err)
		return
	}

	result := <-ch

	var skull local.Skull
	err = skull.Connect(adapter, result.Address)
	if err != nil {
		println(err.Error())
		return
	}

	println("connected to ", result.Address.String())

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	err = skull.Disconnect()
	if err != nil {
		println(err)
	}
}

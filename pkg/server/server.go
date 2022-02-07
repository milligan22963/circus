// Package server is for all server related items
package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/milligan22963/circus/config"
	"github.com/milligan22963/circus/pkg/management"
	"github.com/milligan22963/circus/pkg/web"
	"tinygo.org/x/bluetooth"
)

// HTTPResponse is a structure defining what a response should look like
type HTTPResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
}

// ServerInstance is an instance of server
type ServerInstance struct {
	ServerPort int
	Skull      *management.Skull
}

func (server *ServerInstance) waitForExit() {
	signals := make(chan os.Signal, 1)
	doneFlag := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("\nending operation")
		doneFlag <- true
	}()

	<-doneFlag
}

var adapter = bluetooth.DefaultAdapter

func (server *ServerInstance) SetupSkull() (*management.Skull, error) {
	err := adapter.Enable()

	if err != nil {
		return nil, err
	}

	targetName := "Dark Circus"
	ch := make(chan bluetooth.ScanResult, 1)

	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		fmt.Printf("%+v, payload: %+v", result, result.AdvertisementPayload)
		if result.LocalName() == targetName {
			adapter.StopScan()
			ch <- result
		}
	})

	if err != nil {
		return nil, err
	}

	result := <-ch

	skull := &management.Skull{}
	err = skull.Connect(adapter, result.Address)
	if err != nil {
		return nil, err
	}
	return skull, nil
}

func (server *ServerInstance) Run(appConfig *config.AppConfiguration) {
	defer func() {
		server.Skull.Disconnect()
	}()

	var err error
	server.Skull, err = server.SetupSkull()

	if err != nil {
		return
	}

	webServer := web.WebServer{}

	// server up the world
	go webServer.SetupWebserver(appConfig)

	server.waitForExit()

	appConfig.AppActive <- struct{}{}
}

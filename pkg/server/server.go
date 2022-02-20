// Package server is for all server related items
package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/milligan22963/circus/config"
	"github.com/milligan22963/circus/pkg/artifacts"
	"github.com/milligan22963/circus/pkg/management"
	"github.com/milligan22963/circus/pkg/web"
	"github.com/stianeikeland/go-rpio"
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

func (server *ServerInstance) waitForSkull(appConfig *config.AppConfiguration) {
	select {
	case skull := <-appConfig.Skull:
		if skull != nil {
			server.Skull = skull
			appConfig.GetLogger().Information("the skull is ready")
		}
	case <-appConfig.AppActive:
		return
	}
}

func (server *ServerInstance) Run(appConfig *config.AppConfiguration) {
	defer func() {
		if server.Skull != nil {
			err := server.Skull.Disconnect()
			if err != nil {
				appConfig.GetLogger().Errorf("failed to disconnect from skull: %+v", err)
			}
		}
	}()

	webServer := web.WebServer{}

	// server up the world
	go webServer.SetupWebserver(appConfig)

	go server.waitForSkull(appConfig)

	// Setup gpio
	err := rpio.Open()
	if err != nil {
		appConfig.GetLogger().Fatalf("failed to open GPIO %v.", err)
	}

	// Load up artifacts
	tableArtifacts := artifacts.Artifacts{}

	go tableArtifacts.SetupArtifacts(appConfig)

	// Assuming we have opened everything
	defer rpio.Close()

	server.waitForExit()

	appConfig.AppActive <- struct{}{}
}

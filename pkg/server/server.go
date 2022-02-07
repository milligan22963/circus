// Package server is for all server related items
package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// HTTPResponse is a structure defining what a response should look like
type HTTPResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
}

// ServerInstance is an instance of server
type ServerInstance struct {
	ServerPort int
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

func (server *ServerInstance) Run(appConfig *config.AppConfiguration) {
	defer func() {
		if appConfig.Database != nil {
			appConfig.Database.Close()
		}
	}()

	webServer := web.WebServer{}

	// server up the world
	go webServer.SetupWebserver(appConfig)

	server.waitForExit()

	appConfig.AppActive <- struct{}{}
}

// Package artifacts are for all artifact related processing
package artifacts

import (
	"time"

	"github.com/milligan22963/afmlog"
	"github.com/milligan22963/circus/config"
	"github.com/milligan22963/circus/pkg/management"
	"github.com/stianeikeland/go-rpio"
	"tinygo.org/x/bluetooth"
)

const (
	artifactReadyWaitTime = 1
)

type Artifact struct {
	pinID     uint
	activePin rpio.Pin
}

type Artifacts struct {
	artifacts []Artifact
	logger    *afmlog.Log
	connected bool
}

func CreateArtifact(pinID uint) (*Artifact, error) {
	artifact := Artifact{pinID: pinID}

	artifact.activePin = rpio.Pin(pinID)

	artifact.activePin.PullUp()

	return &artifact, nil
}

func (artifact *Artifact) Monitor() rpio.State {
	return artifact.activePin.Read()
}

func (a *Artifacts) AddArtifact(artifact *Artifact) {
	a.artifacts = append(a.artifacts, *artifact)
}

func (a *Artifacts) connectToSkull(adapter *bluetooth.Adapter) *management.Skull {
	targetName := "Dark Circus"
	ch := make(chan bluetooth.ScanResult, 1)

	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() == targetName {
			a.logger.Information("Detected a skull")
			err := adapter.StopScan()
			if err != nil {
				a.logger.Errorf("failed to stop bluetooth scanning: %+v", err)
			}
			ch <- result
		}
	})

	if err != nil {
		return nil
	}

	result := <-ch

	skull := &management.Skull{ActiveDevice: nil}
	err = skull.Connect(adapter, result.Address)
	if err != nil {
		return nil
	}
	return skull
}

func (a *Artifacts) SetupArtifacts(appConfig *config.AppConfiguration) {
	a.logger = appConfig.GetLogger()

	// create each of the artifacts
	go a.Monitor(appConfig.Adapter, appConfig.AppActive, appConfig.Skull)

	<-appConfig.AppActive
}

func (a *Artifacts) Monitor(adapter *bluetooth.Adapter, appActive chan struct{}, skull chan *management.Skull) {
	ticker := time.NewTicker(artifactReadyWaitTime * time.Second)
	select {
	case <-appActive:
		return
	case <-ticker.C:
		if a.Ready() && !a.connected {
			// Connect to skull
			a.logger.Information("connecting to skull.")
			skullConnection := a.connectToSkull(adapter)
			if skullConnection != nil {
				skull <- skullConnection
				a.connected = true
			}
		}
	}
}

func (a *Artifacts) Ready() bool {
	// scan each artifact and return true if all in place
	for _, artifact := range a.artifacts {
		if artifact.Monitor() == rpio.Low {
			return false
		}
	}
	return true
}

func (a *Artifacts) Reset() {
	a.connected = false

	// do we want to cache the artifacts so if they move after being detected, we don't care?
}

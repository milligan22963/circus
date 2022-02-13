// Package artifacts are for all artifact related processing
package artifacts

import (
	"fmt"
	"time"

	"github.com/milligan22963/afmlog"
	"github.com/milligan22963/circus/config"
	"github.com/milligan22963/circus/pkg/management"
	"github.com/stianeikeland/go-rpio"
	"tinygo.org/x/bluetooth"
)

type Artifact struct {
	pinID     uint
	activePin rpio.Pin
}

type Artifacts struct {
	artifacts []Artifact
	logger    *afmlog.Log
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
		//fmt.Printf("%+v, payload: %+v", result, result.AdvertisementPayload)
		if result.LocalName() == targetName {
			fmt.Printf("Found it!")
			adapter.StopScan()
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
	go a.Monitor(appConfig.Adapter, appConfig.AppActive)

	<-appConfig.AppActive
}

func (a *Artifacts) Monitor(adapter *bluetooth.Adapter, appActive chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	select {
	case <-appActive:
		return
	case <-ticker.C:
		if a.Ready() {
			// Connect to skull
			a.connectToSkull(adapter)
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

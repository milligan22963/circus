// Package management is for all bluetooth management
package management

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

type Skull struct {
	ActiveDevice *bluetooth.Device
}

func (skull *Skull) Connect(adapter *bluetooth.Adapter, bdAddr bluetooth.Addresser) error {
	skull.ActiveDevice = nil
	device, err := adapter.Connect(bdAddr, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	skull.ActiveDevice = device
	return nil
}

func (skull *Skull) ProcessRFID(chipID string) {
	// see if it is one we recognize and if in the right order

}

func (skull *Skull) Disconnect() error {
	if skull.ActiveDevice != nil {
		return skull.ActiveDevice.Disconnect()
	}
	return fmt.Errorf("skull not connected")
}

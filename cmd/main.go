package main

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	err := adapter.Enable()

	if err != nil {
		fmt.Printf("failed to enable adapter: %+v", err)
		return
	}

	connectAddress := "123"
	ch := make(chan bluetooth.ScanResult, 1)

	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		println("found device:", result.Address.String(), result.RSSI, result.LocalName())
		if result.Address.String() == connectAddress {
			adapter.StopScan()
			ch <- result
		}
	})

	if err != nil {
		fmt.Printf("failed to eenable scan: %+v", err)
		return
	}

	var device *bluetooth.Device
	result := <-ch

	var skull management.Skull
	err = skull.Connect(adapter, result.Address)
	if err != nil {
		println(err.Error())
		return
	}

	println("connected to ", result.Address.String())

	err = device.Disconnect()
	if err != nil {
		println(err)
	}
}

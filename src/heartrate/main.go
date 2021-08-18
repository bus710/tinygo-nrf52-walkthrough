package main

import (
	"math/rand"
	"time"

	_ "github.com/bus710/tinygo-nrf52-walkthrough/src/heartrate/app"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// TODO: use atomics to access this value.
var heartRate uint8 = 75 // 75bpm

func main() {

	// r := C.print()
	// println(r)
	_ := app.print()

	println("starting")
	println("ServiceUUIDHeartRate: ", bluetooth.ServiceUUIDHeartRate.String())
	println("CharacteristicUUIDHeartRateMeasurement: ", bluetooth.CharacteristicUUIDHeartRateMeasurement.String())

	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	svc := make([]bluetooth.UUID, 0)
	svc = append(svc, bluetooth.ServiceUUIDHeartRate)
	svc = append(svc, bluetooth.ServiceUUIDBondManagement)
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Go HRS 2",
		ServiceUUIDs: svc,
	}))
	must("start adv", adv.Start())

	var heartRateMeasurement bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHeartRate,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &heartRateMeasurement,
				UUID:   bluetooth.CharacteristicUUIDHeartRateMeasurement,
				Value:  []byte{0, heartRate},
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))

	tickFunc := func() {
		for {
			println("tick", time.Now().Format("04:05.000"))
			time.Sleep(time.Second)

			// random variation in heartrate
			heartRate = randomInt(65, 85)

			// and push the next notification
			heartRateMeasurement.Write([]byte{0, heartRate})
		}
	}
	tickFunc()
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

// Returns an int >= min, < max
func randomInt(min, max int) uint8 {
	return uint8(min + rand.Intn(max-min))
}

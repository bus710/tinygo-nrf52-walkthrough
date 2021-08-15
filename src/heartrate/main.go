// 2021/08/15: I tried adding BMS protocol but it is in progress.

package main

import (
	"math/rand"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// TODO: use atomics to access this value.
var heartRate uint8 = 75 // 75bpm

func main() {
	println("starting")
	println("ServiceUUIDHeartRate: ", bluetooth.ServiceUUIDHeartRate.String())
	println("CharacteristicUUIDHeartRateMeasurement: ", bluetooth.CharacteristicUUIDHeartRateMeasurement.String())

	println("ServiceUUIDBondManagement: ", bluetooth.ServiceUUIDBondManagement.String())
	println("CharacteristicUUIDHeartRateMeasurement: ", bluetooth.CharacteristicUUIDBondManagementControlPoint.String())

	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	svc := make([]bluetooth.UUID, 0)
	svc = append(svc, bluetooth.ServiceUUIDHeartRate)
	svc = append(svc, bluetooth.ServiceUUIDBondManagement)
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Go HRS 3",
		ServiceUUIDs: svc,
	}))
	must("start adv", adv.Start())

	var bondManagement bluetooth.Characteristic
	value := make([]byte, 5)
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDBondManagement,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &bondManagement,
				UUID:   bluetooth.CharacteristicUUIDBondManagementControlPoint,
				Value:  value,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))

	checkFunc := func() {
		for {
			time.Sleep(time.Second)
			println(string(value))
		}
	}
	go checkFunc()

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
		nextBeat := time.Now()
		for {
			nextBeat = nextBeat.Add(time.Minute / time.Duration(heartRate))
			println("tick", time.Now().Format("04:05.000"))
			time.Sleep(nextBeat.Sub(time.Now()))

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

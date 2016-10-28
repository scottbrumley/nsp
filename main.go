package main

import (
	"github.com/scottbrumley/nsp"
	"fmt"
)

func main() {
	hosts := make(chan nsp.Sensors)
	results := make(chan string)

	myParms := nsp.GetParams()
	//fmt.Print(nsp.SshCommand(myParms))
	deviceList := nsp.GetSensorList("sensors.csv")

	fmt.Println("Spin up workers")
	// Spin up 10 workers which block until there is work
	for w := 1; w <= 10; w++ {
		go nsp.SshMultiple(hosts,results,myParms)
	}

	fmt.Println("Fill Channel")
	// Fill channels with devices to ssh
	for _,host := range deviceList {
		fmt.Println("Channel Add: " + host.Name)
		hosts <- host
	}

	// Close channel to let everyone know we are done
	close(hosts)

	for _,device := range deviceList {
		fmt.Println(device.Name)
		<-results
	}
}

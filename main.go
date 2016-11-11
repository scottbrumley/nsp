package main

import (
	"github.com/scottbrumley/nsp"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	hosts := make(chan nsp.Sensors)
	results := make(chan string)
	var totalBlackListFiles int64


	myParms := nsp.GetParams()
	deviceList := nsp.GetSensorList("sensors.csv")

	fmt.Println("Spin up workers")
	// Spin up 10 workers which block until there is work
	for w := 1; w <= 100; w++ {
		go nsp.SshMultiple(hosts,results,myParms)
	}

	fmt.Println("Fill Channel")
	// Fill channels with devices to ssh
	for _,host := range deviceList {
		fmt.Println("Channel Add: " + host.Name)
		hosts <- host
	}

	allResults := make(map[string]string)

	//Collect results
	for _,device := range deviceList {
		allResults[device.Name] = <-results
	}

	// Print Out Results in a nice format
	for key,value := range allResults {
		if (value != "") {
//			fmt.Print("Device: " + key + " Results: " + value)
			fmt.Println("Device: " + key)
			stanzaString := getStanza(value,"MALWARE STATISTICS FOR BLACKLIST ENGINE:")
			totalBlackListFiles = totalBlackListFiles + getValue(stanzaString,"Number of files sent:")
		}
	}

	// Close channel to let everyone know we are done
	close(hosts)
	fmt.Printf("Total Files Black Listed: %d \n", totalBlackListFiles)
}

func getStanza(resultStr string, stringMatch string)(string){
	retVal := strings.Split(resultStr,stringMatch)
	return retVal[1]
}

func getValue(resultStr string, stringMatch string)(int64){
	var retVal int64
	var match []string

	if ( (len(stringMatch) > 0) && (len(resultStr) > 0) ) {
		myString := stringMatch + "\\s*(\\d*)\\s"
		re := regexp.MustCompile(myString)
		match = re.FindStringSubmatch(resultStr)
	}

	if len(match) <= 0 {
		retVal = 0
	} else {
		if i, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			retVal = i
		}

	}
	return retVal

}
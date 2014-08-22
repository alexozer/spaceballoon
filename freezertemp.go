package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

const logFile = "/home/alex/log/temperature.csv"
const tempSrc = "/sys/class/thermal/thermal_zone0/temp"
const interval = 3

var startTime time.Time

func main() {
	logfile, err := os.Create(logFile)
	if err != nil {
		log.Fatal(err)
	}

	csvWriter := csv.NewWriter(logfile)
	defer logfile.Close()
	csvWriter.Write([]string{"Time (s)", "Temperature (C)"})
	csvWriter.Flush()

	ticker := time.NewTicker(time.Second * interval)
	startTime = time.Now()
	for time := range ticker.C {
		csvWriter.Write([]string{
			strconv.FormatFloat(time.Sub(startTime).Seconds(), 'f', 3, 64),
			strconv.FormatFloat(getTemperature(), 'f', 3, 64),
		})
		csvWriter.Flush()
	}
}

func getTemperature() float64 {
	file, err := os.Open(tempSrc)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tempBytes := make([]byte, 5)
	_, err = file.Read(tempBytes)
	if err != nil {
		log.Fatal(err)
	}

	val, err := strconv.ParseInt(string(tempBytes), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return float64(val) / 1000.0
}

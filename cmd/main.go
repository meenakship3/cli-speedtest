package main

import (
	"fmt"

	"github.com/showwin/speedtest-go/speedtest"
)

func checkSpeed() (float64, float64, error) {
	var speedtestClient = speedtest.New()

	serverList, err := speedtestClient.FetchServers()
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}

	var downloadSpeed, uploadSpeed float64

	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		fmt.Printf("Latency: %s, Download: %s, Upload: %s\n", s.Latency, s.DLSpeed, s.ULSpeed)
		downloadSpeed = s.DLSpeed.Mbps()
		uploadSpeed = s.ULSpeed.Mbps()
		s.Context.Reset()
	}
	return downloadSpeed, uploadSpeed, nil
}

func main() {
	downloadSpeed, uploadSpeed, err := checkSpeed()
	if err != nil {
		fmt.Print(err)
		return
	}
	if downloadSpeed < 50.00 { // TODO: remove hardcoding
		fmt.Println("Download speed is below 50 mbps")
	}
	fmt.Print(downloadSpeed, uploadSpeed, err)
}

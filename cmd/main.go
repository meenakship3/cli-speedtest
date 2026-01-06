package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/showwin/speedtest-go/speedtest"
)

func getFlags() (float64, float64, bool, float64, bool) {
	dlT := flag.Float64("download-threshold", 50.00, "sets alert threshold (mbps) for download speeds")
	ulT := flag.Float64("upload-threshold", 50.00, "sets alert threshold (mbps) for upload speeds")
	ulM := flag.Bool("upload-monitoring", false, "checks upload speed along with download speed")
	interval := flag.Float64("interval", 60.00, "sets check interval (in minutes)")
	once := flag.Bool("run-once", true, "sets test frequency")

	flag.Parse()

	return *dlT, *ulT, *ulM, *interval, *once
}

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

func getWebhookURL() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	webhookURL, isSet := os.LookupEnv("WEBHOOK_URL")
	if !isSet || webhookURL == "" {
		log.Fatal("Error getting Slack Webhook URL")
	}
	return webhookURL
}
func notify(messageContent string) string {
	webhookURL := getWebhookURL()

	postBody, _ := json.Marshal(map[string]string{
		"text": messageContent,
	})

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(webhookURL, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return sb
}

func getCurrentDateTime() string {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return ""
	}

	now := time.Now().In(loc)

	formattedTime := now.Format("02/01/2006 15:04:05 MST")
	return formattedTime
}

func main() {
	downloadThreshold, uploadThreshold, uploadMonitoring, interval, once := getFlags()
	duration := time.Duration(interval) * time.Minute

	for {
		downloadSpeed, uploadSpeed, err := checkSpeed()
		if err != nil {
			fmt.Print(err)
			return
		}
		if downloadSpeed < downloadThreshold {
			fmt.Printf("Download speed is below %.2f mbps", downloadThreshold)
			currentDateTime := getCurrentDateTime()
			messageContent := fmt.Sprintf("Uh oh! Download speed at %s was below %.2f mbps. Current download speed is %.2f mbps.", currentDateTime, downloadThreshold, downloadSpeed)
			notify(messageContent)
		}
		if uploadMonitoring && uploadSpeed < uploadThreshold {
			fmt.Printf("Upload speed is below %.2f mbps", uploadThreshold)
			currentDateTime := getCurrentDateTime()
			messageContent := fmt.Sprintf("Uh oh! Upload speed at %s was below %.2f mbps. Current upload speed is %.2f mbps.", currentDateTime, uploadThreshold, uploadSpeed)
			notify(messageContent)
		}
		if once {
			break
		}
		time.Sleep(duration)
	}

}

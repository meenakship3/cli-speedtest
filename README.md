# Internet Speed Monitor

A Go-based CLI tool for monitoring internet connection speeds with configurable thresholds and Slack notifications.

## Features

- **Configurable thresholds** via CLI flags for download and upload speeds.
- **Slack webhook integration** for instant notifications when speeds drop below thresholds.
- **Flexible monitoring intervals** - set custom check frequencies.
- **Timezone-aware logging** for accurate alert timestamps. Currently set to Asia/Kolkata.

## Installation

```bash
go get github.com/showwin/speedtest-go
go get github.com/joho/godotenv
go build -o speedmonitor
```

## Usage

Create a `.env` file with your Slack webhook URL:
```
WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

Run with custom thresholds:
```bash
./speedmonitor --download-threshold 50 --upload-threshold 10 --interval 60
```

### Available Flags

- `--download-threshold` - Alert threshold for download speed in Mbps (default: 50)
- `--upload-threshold` - Alert threshold for upload speed in Mbps (default: 50)
- `--upload-monitoring` - Enable upload speed monitoring (default: false)
- `--interval` - Check interval in minutes (default: 60)
- `--once` - Set test/run frequency (default: true)
- `--help` - View available flags

## Technical Implementation

- Built with Go's standard library (`net/http`, `encoding/json`, `flag`)
- Uses `showwin/speedtest-go` for reliable speed testing via Speedtest.net infrastructure
- Implements RESTful HTTP POST requests for Slack webhook integration
- Continuous monitoring loop with configurable intervals

## Use Case

Ideal for monitoring home/office internet reliability, tracking ISP performance, or alerting teams to connectivity issues in real-time.

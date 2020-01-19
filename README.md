# Route 53 record updater

[![Go Report Card](https://goreportcard.com/badge/github.com/fsaravia/route-53-updater)](https://goreportcard.com/report/github.com/fsaravia/route-53-updater)

## Usage:

```bash
RESOLVER_URL='<AN_URL>' \
API_KEY='<AN_API_KEY>' \
HOSTED_ZONE_ID='<YOUR_ZONE_ID>' \
RECORD_SET='<YOUR_RECORD_SET>' \
go run ./main.go
```

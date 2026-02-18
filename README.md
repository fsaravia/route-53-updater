# Route 53 record updater

[![Go Report Card](https://goreportcard.com/badge/github.com/fsaravia/route-53-updater)](https://goreportcard.com/report/github.com/fsaravia/route-53-updater)

## What?

This is a small Go tool that updates an AWS Route53 `A` record with a dynamic public IP address.

The process has two stages:

### 1. Resolve the current public IP address

The tool expects an external HTTP endpoint that accepts a `GET` request and returns JSON in either of these forms:

```json
{
  "ip": "127.0.0.1"
}
```

or:

```json
{
  "ip_address": "127.0.0.1"
}
```

An example implementation of such a service can be found [here](https://gist.github.com/fsaravia/13f4b94d5a370b1198f8474422c8b862).

### 2. Update Route53

Once the public IP is resolved, the tool sends an UPSERT request to Route53 for the configured `A` record.

## Configuration

Configuration is provided through environment variables. If any required variable is missing, the process exits with an error.

- `RESOLVER_URL`: URL of the IP resolver endpoint.
- `API_KEY`: API key sent as the `x-api-key` request header.
- `HOSTED_ZONE_ID`: Route53 hosted zone ID.
- `RECORD_SET`: FQDN of the record to update.

## Logging

The application logs extensively to STDOUT with UTC timestamps, including:

- startup and configuration loading (without printing secrets),
- resolver request/response metadata and IP parsing,
- Route53 client initialization and UPSERT submission,
- success/failure status and error context.

## Usage

```bash
RESOLVER_URL='<AN_URL>' \
API_KEY='<AN_API_KEY>' \
HOSTED_ZONE_ID='<YOUR_ZONE_ID>' \
RECORD_SET='<YOUR_RECORD_SET>' \
go run ./main.go
```

## Example

```bash
RESOLVER_URL='https://resolver.example.org/whatsmyip' \
API_KEY='123456' \
HOSTED_ZONE_ID='FOO' \
RECORD_SET='dynamic.example.org' \
go run ./main.go
```

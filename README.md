# Route 53 record updater

[![Go Report Card](https://goreportcard.com/badge/github.com/fsaravia/route-53-updater)](https://goreportcard.com/report/github.com/fsaravia/route-53-updater)

## What?

This is a small tool that updates a record in AWS Route53 with a dynamic IP address.

It works on two stages:

### Resolve local IP address

In order to resolve the local IP address, this tool requires an external service that accepts a GET request and returns a JSON response in the following form:

```json
{
	"ip": "127.0.0.1"
}
```

An example for the code of such a service can be found [here](https://gist.github.com/fsaravia/13f4b94d5a370b1198f8474422c8b862)

### Update IP address on Route53

Once the local IP address has been obtained, it calls Route53 and updates an `A` record with the value of said IP address.

## Configuration

This tool depends entirely on environment variables, if unset or empty, it will fail.

* `RESOLVER_URL`: The URL of the IP resolver to which a GET request will be sent.
* `API_KEY`: The API key of the URL resolver. This API key will be sent as a header under the value `x-api-key`, following the AWS Lambda conventions.
* `HOSTED_ZONE_ID`: The ID of the AWS Route53 hosted zone.
* `RECORD_SET`: The FQDN of the record set to update.

## Usage:

```bash
RESOLVER_URL='<AN_URL>' \
API_KEY='<AN_API_KEY>' \
HOSTED_ZONE_ID='<YOUR_ZONE_ID>' \
RECORD_SET='<YOUR_RECORD_SET>' \
go run ./main.go
```

## Example:

```bash
RESOLVER_URL='https://resolver.example.org/whatsmyip' \
API_KEY='123456' \
HOSTED_ZONE_ID='FOO' \
RECORD_SET='dynamic.example.org' \
go run ./main.go
```

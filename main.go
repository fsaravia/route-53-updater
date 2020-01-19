package main

import (
	"./ipresolver"
	"./route53"
	"fmt"
	"os"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic(fmt.Sprintf("Please set the %s environment variable", key))
	}

	return value
}

func resolverURL() string {
	return getEnv("RESOLVER_URL")
}

func resolverAPIKey() string {
	return getEnv("API_KEY")
}

func route53HostedZoneId() string {
	return getEnv("HOSTED_ZONE_ID")
}

func recordSet() string {
	return getEnv("RECORD_SET")
}

func main() {
	ip, err := ipresolver.ResolveIp(resolverURL(), resolverAPIKey())
	if err != nil {
		panic(err)
	}

	session := route53.CreateSession()

	output, err := route53.UpsertZone(session, route53HostedZoneId(), recordSet(), ip)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}

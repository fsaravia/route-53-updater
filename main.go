package main

import (
	"fmt"
	"os"
	"./ipResolver"
	"./route53"
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

func main() {
	ip, err := ipResolver.ResolveIp(resolverURL(), resolverAPIKey())
	if err != nil {
		panic(err)
	}

	fmt.Printf("The IP address of this device is: %s", ip)
}

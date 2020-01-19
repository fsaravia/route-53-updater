package main

import (
	"fmt"
	"./ipResolver"
)

func main() {
	ip, err := ipResolver.ResolveIp()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The IP address of this device is: %s", ip)
}

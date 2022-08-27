package cmd

import (
	"fmt"
	"log"
	"strings"

	gmx "github.com/arinn1204/gmx/pkg"
)

// Run is the main entry point for the cli
func Run(domain string, name string, operation string, args string, types string, containerTypes string) {
	client := gmx.CreateClient()
	err := client.Initialize()

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	id, err := client.RegisterBean("127.0.0.1", 9001)

	if err != nil {
		log.Fatal(err)
	}

	mbeanArgs := make([]gmx.MBeanArgs, 0)

	splitArgs := strings.Split(args, ",")
	splitTypes := strings.Split(types, ",")
	splitContainers := strings.Split(containerTypes, ",")

	for i := range splitArgs {
		containerType := ""
		if i < len(splitContainers) {
			containerType = splitContainers[i]
		}
		mbeanArgs = append(mbeanArgs, gmx.MBeanArgs{
			Value:             splitArgs[i],
			JavaType:          splitTypes[i],
			JavaContainerType: containerType,
		})
	}

	operator := client.GetOperator()

	result, err := operator.ExecuteAgainstID(
		*id,
		domain,
		name,
		operation,
		mbeanArgs...,
	)

	if err != nil {
		log.Fatal(err)
	} else if result != "" {
		fmt.Println(result)
	}
}

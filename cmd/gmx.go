package cmd

import (
	"fmt"
	"log"
	"strings"

	gmx "github.com/arinn1204/gmx/pkg"
)

func Run(domain string, name string, operation string, args string, types string) {
	client := &gmx.Client{}
	err := client.Initialize()

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	id, err := client.Connect("127.0.0.1", 9001)

	if err != nil {
		log.Fatal(err)
	}

	mbeanArgs := make([]gmx.MBeanArgs, 0)

	splitArgs := strings.Split(args, ",")
	splitTypes := strings.Split(types, ",")

	for i := range splitArgs {
		mbeanArgs = append(mbeanArgs, gmx.MBeanArgs{
			Value:    splitArgs[i],
			JavaType: splitTypes[i],
		})
	}

	result, err := client.ExecuteAgainstID(
		*id,
		domain,
		name,
		operation,
		mbeanArgs...,
	)

	if err != nil {
		log.Fatal(err)
	} else if result != nil && (result != "") {
		fmt.Println(result)
	}
}

package main

import (
	"flag"
	"fmt"
	"gmx/pkg/gmx"
	"log"
	"strings"
)

var (
	domain    *string
	name      *string
	operation *string
	args      *string
	types     *string
)

func init() {
	domain = flag.String("domain", "", "The domain for the mbean")
	name = flag.String("name", "", "The name of the mbean itself")
	operation = flag.String("operation", "", "The operation that is being executed")
	args = flag.String("arguments", "", "the optional comma separated arguments passed into the operation")
	types = flag.String("types", "", "the comma separated types that correspond with the arguments passed in")
}

func validateArgs() {
	if domain == nil || *domain == "" {
		log.Fatal("did not receive a domain to execute against")
	}

	if name == nil || *name == "" {
		log.Fatal("did not receive a name to execute against")
	}

	if operation == nil || *operation == "" {
		log.Fatal("did not receive a domain to operation")
	}

	if (args != nil && types == nil) ||
		(args == nil && types != nil) ||
		(strings.Count(*args, ",") != strings.Count(*types, ",")) {
		log.Fatal("must provide types for all the args provided")
	}
}

func main() {
	flag.Parse()
	validateArgs()

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

	splitArgs := strings.Split(*args, ",")
	splitTypes := strings.Split(*types, ",")

	for i := range splitArgs {
		mbeanArgs = append(mbeanArgs, gmx.MBeanArgs{
			Value:    splitArgs[i],
			JavaType: splitTypes[i],
		})
	}

	result, err := client.ExecuteAgainstID(
		*id,
		*domain,
		*name,
		*operation,
		mbeanArgs...,
	)

	if err != nil {
		log.Fatal(err)
	} else if result != nil && (result != "") {
		fmt.Println(result)
	}
}

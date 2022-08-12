package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/arinn1204/gmx/cmd"
)

var (
	domain    *string
	name      *string
	operation *string
	hostname  *string
	port      *int
	args      *string
	types     *string
)

func init() {
	hostname = flag.String("hostname", "", "[Required] The hostname or IP of the JMX RMI server")
	port = flag.Int("port", 0, "[Required] The port of the JMX RMI server")
	domain = flag.String("domain", "", "[Required] The domain for the mbean")
	name = flag.String("name", "", "[Required] The name of the mbean itself")
	operation = flag.String("operation", "", "[Required] The operation that is being executed")
	args = flag.String("arguments", "", "[Optional] The comma separated arguments passed into the operation")
	types = flag.String("types", "", "[Optional] The comma separated types that correspond with the arguments passed in")
}

func exit(predicate func() bool) {
	if predicate() {
		flag.Usage()
		os.Exit(1)
	}
}

func validateArgs() {
	exit(func() bool { return domain == nil || *domain == "" })
	exit(func() bool { return name == nil || *name == "" })
	exit(func() bool { return operation == nil || *operation == "" })
	exit(func() bool { return hostname == nil || *hostname == "" })
	exit(func() bool { return port == nil || *port == 0 })

	if (args != nil && types == nil) ||
		(args == nil && types != nil) ||
		(*args != "" && *types == "") ||
		(*args == "" && *types != "") ||
		(strings.Count(*args, ",") != strings.Count(*types, ",")) {
		flag.Usage()
		log.Fatal("\nMust provide types for all the args provided")
	}
}

func main() {
	flag.Parse()
	validateArgs()
	cmd.Run(*domain, *name, *operation, *args, *types)
}

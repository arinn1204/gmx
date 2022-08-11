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
	args      *string
	types     *string
)

func init() {
	domain = flag.String("domain", "", "[Required] The domain for the mbean")
	name = flag.String("name", "", "[Required] The name of the mbean itself")
	operation = flag.String("operation", "", "[Required] The operation that is being executed")
	args = flag.String("arguments", "", "[Optional] The comma separated arguments passed into the operation")
	types = flag.String("types", "", "[Optional] The comma separated types that correspond with the arguments passed in")
}

func validateArgs() {
	if domain == nil || *domain == "" {
		flag.Usage()
		os.Exit(0)
	}

	if name == nil || *name == "" {
		flag.Usage()
		os.Exit(0)
	}

	if operation == nil || *operation == "" {
		flag.Usage()
		os.Exit(0)
	}

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

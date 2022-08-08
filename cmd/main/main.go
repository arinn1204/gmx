package main

import (
	"fmt"
	"gmx/internal/java"
)

func main() {

	mbean := &java.MBean{}

	err := mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	defer mbean.Close()
	operation := java.MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "putString",
		Args: []java.MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
			{
				Value: "fan369",
				Type:  "java.lang.String",
			},
		},
	}
	mbean.Execute(operation)

	operation = java.MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "getString",
		Args: []java.MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
		},
	}
	res, err := mbean.Execute(operation)

	fmt.Println(res)
	fmt.Println(err)
}

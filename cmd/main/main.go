package main

import (
	"fmt"
	"gmx/internal/mbean"
	"log"
)

func main() {

	beanExecutor := &mbean.MBean{}

	if err := beanExecutor.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi"); err != nil {
		log.Panicf("failed to initialize the connection::%s", err.Error())
	}

	defer beanExecutor.Close()
	operation := mbean.MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "putString",
		Args: []mbean.MBeanOperationArgs{
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
	beanExecutor.Execute(operation)

	operation = mbean.MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "getString",
		Args: []mbean.MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
		},
	}
	res, err := beanExecutor.Execute(operation)

	fmt.Println(res)
	fmt.Println(err)
}

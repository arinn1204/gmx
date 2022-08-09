package main

import (
	"fmt"
	"gmx/internal/jvm"
	"gmx/internal/mbean"
	"log"
)

func main() {

	javaVm, err := jvm.CreateJvm()

	if err != nil {
		log.Panicf("failed to start jvm::%s", err.Error())
	}

	defer javaVm.ShutdownJvm()

	beanExecutor, err := jvm.CreateMBeanConnection(javaVm, "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")

	if err != nil {
		log.Panicf("failed to initialize the connection::%s", err.Error())
	}

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
